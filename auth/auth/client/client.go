package authclient

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"sample/auth/auth"
	"sample/common/auth/token"
	chttp "sample/common/http"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	kitconsul "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	kithttp "github.com/go-kit/kit/transport/http"
)

func New(instance string, client *http.Client) (auth.Service, error) {
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

	signinEndpoint := kithttp.NewClient(
		http.MethodPost,
		copyURL(u, "/auth/signin"),
		chttp.EncodeJsonReq,
		auth.DecodeSigninResponse,
		opts...,
	).Endpoint()

	verifyTokenEndpoint := kithttp.NewClient(
		http.MethodPost,
		copyURL(u, "/auth/verify"),
		chttp.EncodeJsonReq,
		auth.DecodeVerifyTokenResponse,
		opts...,
	).Endpoint()

	return &auth.Endpoint{
		SigninEndpoint:      auth.SigninEndpoint(signinEndpoint),
		VerifyTokenEndpoint: auth.VerifyTokenEndpoint(verifyTokenEndpoint),
	}, nil
}

func NewWithLB(instancer *kitconsul.Instancer, retryMax int, retryTimeout time.Duration, logger kitlog.Logger, client *http.Client) auth.Service {
	var endpoints auth.Endpoint
	{
		factory := factoryFor(auth.MakeSigninEndpoint, client)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.SigninEndpoint = auth.SigninEndpoint(retry)
	}
	{
		factory := factoryFor(auth.MakeVerifyTokenEndpoint, client)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.VerifyTokenEndpoint = auth.VerifyTokenEndpoint(retry)
	}
	return endpoints
}

func copyURL(u *url.URL, path string) *url.URL {
	c := *u
	c.Path = path
	return &c
}

func factoryFor(makeEndpoint func(auth.Service) endpoint.Endpoint, client *http.Client) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		service, err := New(instance, client)
		if err != nil {
			return nil, nil, err
		}
		return makeEndpoint(service), nil, nil
	}
}
