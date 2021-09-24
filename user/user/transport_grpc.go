package user

import (
	"context"

	kitgrpc "github.com/go-kit/kit/transport/grpc"

	"sample/user/user/pb"
)

type grpcHandler struct {
	pb.UnimplementedUserServer
	add           kitgrpc.Handler
}

func NewGRPCHandler(s Service) pb.UserServer {
	var opts []kitgrpc.ServerOption
	add := kitgrpc.NewServer(
		MakeAddEndpoint(s),
		decodeGRPCAddRequest,
		encodeGRPCAddResponse,
		opts...,
	)

	/*
		get := kithttp.NewServer(
			MakeGetEndpoint(s),
			DecodeGetRequest,
			chttp.EncodeJSONResponse,
			opts...,
		)

		checkPassword := kithttp.NewServer(
			MakeCheckPasswordEndpoint(s),
			DecodeCheckPasswordRequest,
			chttp.EncodeJSONResponse,
			opts...,
		)
	*/

	return &grpcHandler{
		add: add,
	}
}

func (s *grpcHandler) Add(ctx context.Context, req *pb.AddReq) (*pb.AddResp, error) {
	_, rep, err := s.add.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.AddResp), nil
}

func decodeGRPCAddRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.AddReq)
	return addRequest{Name: req.Name, Username: req.Username, Password: req.Password}, nil
}

func encodeGRPCAddResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(addResponse)
	return &pb.AddResp{Id: resp.Id}, nil
}

func EncodeGRPCAddRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(addRequest)
	return &pb.AddReq{Name: req.Name, Username: req.Username, Password: req.Password}, nil
}

func DecodeGRPCAddResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	resp := grpcReply.(*pb.AddResp)
	return addResponse{Id: resp.Id}, nil
}
