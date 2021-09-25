package client

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"sample/common/auth/token"
	chttp "sample/common/http"
	"sample/user/user"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	kitconsul "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	kithttp "github.com/go-kit/kit/transport/http"
)

func New(instance string, client *http.Client) (user.Service, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	opts := []kithttp.ClientOption{
		kithttp.SetClient(client),
		kithttp.ClientBefore(token.HTTPTokenFromContext),
	}

	addEndPoint := kithttp.NewClient(
		http.MethodPost,
		copyURL(u, "/user"),
		chttp.EncodeJsonReq,
		user.DecodeAddResponse,
		opts...,
	).Endpoint()

	getEndpoint := kithttp.NewClient(
		http.MethodGet,
		copyURL(u, "/user"),
		chttp.EncodeQueryReq,
		user.DecodeGetResponse,
		opts...,
	).Endpoint()

	checkPasswordEndpoint := kithttp.NewClient(
		http.MethodGet,
		copyURL(u, "/user/check"),
		chttp.EncodeQueryReq,
		user.DecodeCheckPasswordResponse,
		opts...,
	).Endpoint()

	return &user.Endpoint{
		AddEndpoint:           user.AddEndpoint(addEndPoint),
		GetEndpoint:           user.GetEndpoint(getEndpoint),
		CheckPasswordEndpoint: user.CheckPasswordEndpoint(checkPasswordEndpoint),
	}, nil
}

func NewWithLB(instancer *kitconsul.Instancer, retryMax int, retryTimeout time.Duration, logger kitlog.Logger, client *http.Client) user.Service {
	var endpoints user.Endpoint
	{
		factory := factoryFor(user.MakeAddEndpoint, client)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.AddEndpoint = user.AddEndpoint(retry)
	}
	{
		factory := factoryFor(user.MakeGetEndpoint, client)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.GetEndpoint = user.GetEndpoint(retry)
	}
	{
		factory := factoryFor(user.MakeCheckPasswordEndpoint, client)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.CheckPasswordEndpoint = user.CheckPasswordEndpoint(retry)
	}
	return endpoints
}

func copyURL(u *url.URL, path string) *url.URL {
	c := *u
	c.Path = path
	return &c
}

func factoryFor(makeEndpoint func(user.Service) endpoint.Endpoint, client *http.Client) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		service, err := New(instance, client)
		if err != nil {
			return nil, nil, err
		}
		return makeEndpoint(service), nil, nil
	}
}
