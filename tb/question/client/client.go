package client

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	kitconsul "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	kithttp "github.com/go-kit/kit/transport/http"

	"sample/common/auth/token"
	chttp "sample/common/http"
	"sample/tb/question"
)

func New(instance string, client *http.Client) (question.Service, error) {
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	opts := []kithttp.ClientOption{
		kithttp.SetClient(client),
		kithttp.ClientBefore(token.HTTPTokenFromContext),
	}

	addEndpoint := kithttp.NewClient(
		http.MethodPost,
		copyURL(u, "/question"),
		chttp.EncodeJSONRequest,
		question.DecodeAddResponse,
		opts...,
	).Endpoint()

	getEndpoint := kithttp.NewClient(
		http.MethodGet,
		copyURL(u, "/question"),
		chttp.EncodeSchemaRequest,
		question.DecodeGetResponse,
		opts...,
	).Endpoint()

	listEndpoint := kithttp.NewClient(
		http.MethodGet,
		copyURL(u, "/question/list"),
		chttp.EncodeSchemaRequest,
		question.DecodeListResponse,
		opts...,
	).Endpoint()

	return &question.Endpoint{
		AddEndpoint:  question.AddEndpoint(addEndpoint),
		GetEndpoint:  question.GetEndpoint(getEndpoint),
		ListEndpoint: question.ListEndpoint(listEndpoint),
	}, nil
}

func NewWithLB(instancer *kitconsul.Instancer, retryMax int, retryTimeout time.Duration, logger kitlog.Logger, client *http.Client) question.Service {
	var endpoints question.Endpoint
	{
		factory := factoryFor(question.MakeAddEndpoint, client)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.AddEndpoint = question.AddEndpoint(retry)
	}
	{
		factory := factoryFor(question.MakeGetEndpoint, client)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.GetEndpoint = question.GetEndpoint(retry)
	}
	{
		factory := factoryFor(question.MakeListEndpoint, client)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.ListEndpoint = question.ListEndpoint(retry)
	}
	return endpoints
}

func copyURL(u *url.URL, path string) *url.URL {
	c := *u
	c.Path = path
	return &c
}

func factoryFor(makeEndpoint func(question.Service) endpoint.Endpoint, client *http.Client) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		service, err := New(instance, client)
		if err != nil {
			return nil, nil, err
		}
		return makeEndpoint(service), nil, nil
	}
}
