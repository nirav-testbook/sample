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
		copyURL(u, "/lesson"),
		chttp.EncodeJSONRequest,
		lesson.DecodeAddResponse,
		opts...,
	).Endpoint()

	getEndpoint := kithttp.NewClient(
		http.MethodGet,
		copyURL(u, "/lesson"),
		chttp.EncodeSchemaRequest,
		lesson.DecodeGetResponse,
		opts...,
	).Endpoint()

	get1Endpoint := kithttp.NewClient(
		http.MethodGet,
		copyURL(u, "/lesson/1"),
		chttp.EncodeSchemaRequest,
		lesson.DecodeGet1Response,
		opts...,
	).Endpoint()

	return &lesson.Endpoint{
		AddEndpoint:  lesson.AddEndpoint(addEndpoint),
		GetEndpoint:  lesson.GetEndpoint(getEndpoint),
		Get1Endpoint: lesson.Get1Endpoint(get1Endpoint),
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
		factory := factoryFor(lesson.MakeGet1Endpoint, client)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.Get1Endpoint = lesson.Get1Endpoint(retry)
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
