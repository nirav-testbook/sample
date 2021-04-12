package client

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"sample/common/auth/token"
	chttp "sample/common/http"
	"sample/tb/lesson"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	kithttp "github.com/go-kit/kit/transport/http"
)

func New(instance string, client *http.Client) (lesson.Service, error) {
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
		lesson.DecodeAddResponse,
		opts...,
	).Endpoint()

	getEndpoint := kithttp.NewClient(
		http.MethodGet,
		copyURL(u, "/question"),
		chttp.EncodeSchemaRequest,
		lesson.DecodeGetResponse,
		opts...,
	).Endpoint()

	listEndpoint := kithttp.NewClient(
		http.MethodGet,
		copyURL(u, "/question/list"),
		chttp.EncodeSchemaRequest,
		lesson.DecodeListResponse,
		opts...,
	).Endpoint()

	return &lesson.Endpoint{
		AddEndpoint:  lesson.AddEndpoint(addEndpoint),
		GetEndpoint:  lesson.GetEndpoint(getEndpoint),
		ListEndpoint: lesson.ListEndpoint(listEndpoint),
	}, nil
}

func NewWithLB(instancer *kitconsul.Instancer, retryMax int, retryTimeout time.Duration, logger kitlog.Logger, client *http.Client) lesson.Service {
	var endpoints lesson.Endpoint
	{
		factory := factoryFor(lesson.MakeAddEndpoint, client)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.AddEndpoint = lesson.AddEndpoint(retry)
	}
	{
		factory := factoryFor(lesson.MakeGetEndpoint, client)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.GetEndpoint = lesson.GetEndpoint(retry)
	}
	{
		factory := factoryFor(lesson.MakeListEndpoint, client)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.ListEndpoint = lesson.ListEndpoint(retry)
	}
	return endpoints
}

func copyURL(u *url.URL, path string) *url.URL {
	c := *u
	c.Path = path
	return &c
}

func factoryFor(makeEndpoint func(lesson.Service) endpoint.Endpoint, client *http.Client) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		service, err := New(instance, client)
		if err != nil {
			return nil, nil, err
		}
		return makeEndpoint(service), nil, nil
	}
}
