package client

import (
	"io"
	"time"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	kitconsul "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"

	"sample/user/user"
	"sample/user/user/pb"
)

func NewGRPCClient(conn *grpc.ClientConn) (user.Service, error) {
	var options []kitgrpc.ClientOption

	addEndpoint := kitgrpc.NewClient(
		conn,
		"User",
		"Add",
		user.EncodeGRPCAddRequest,
		user.DecodeGRPCAddResponse,
		pb.AddResp{},
		options...,
	).Endpoint()

	getEndpoint := kitgrpc.NewClient(
		conn,
		"User",
		"Get",
		user.EncodeGRPCGetRequest,
		user.DecodeGRPCGetResponse,
		pb.GetResp{},
		options...,
	).Endpoint()

	checkPasswordEndpoint := kitgrpc.NewClient(
		conn,
		"User",
		"CheckPassword",
		user.EncodeGRPCCheckPasswordRequest,
		user.DecodeGRPCCheckPasswordResponse,
		pb.CheckPasswordResp{},
		options...,
	).Endpoint()

	return &user.Endpoint{
		AddEndpoint:           user.AddEndpoint(addEndpoint),
		GetEndpoint:           user.GetEndpoint(getEndpoint),
		CheckPasswordEndpoint: user.CheckPasswordEndpoint(checkPasswordEndpoint),
	}, nil
}

func NewGRPCWithLB(instancer *kitconsul.Instancer, retryMax int, retryTimeout time.Duration, logger kitlog.Logger) user.Service {
	var endpoints user.Endpoint
	{
		factory := grpcFactoryFor(user.MakeAddEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.AddEndpoint = user.AddEndpoint(retry)
	}
	{
		factory := grpcFactoryFor(user.MakeGetEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.GetEndpoint = user.GetEndpoint(retry)
	}
	{
		factory := grpcFactoryFor(user.MakeCheckPasswordEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.CheckPasswordEndpoint = user.CheckPasswordEndpoint(retry)
	}
	return endpoints
}

func grpcFactoryFor(makeEndpoint func(user.Service) endpoint.Endpoint) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		service, err := NewGRPCClient(conn)
		if err != nil {
			return nil, nil, err
		}
		return makeEndpoint(service), nil, nil
	}
}
