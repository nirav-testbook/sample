package user

import (
	"context"

	kitgrpc "github.com/go-kit/kit/transport/grpc"

	cgrpc "sample/common/grpc"
	"sample/user/model"
	"sample/user/user/pb"
)

type grpcHandler struct {
	pb.UnimplementedUserServer
	add           kitgrpc.Handler
	get           kitgrpc.Handler
	checkPassword kitgrpc.Handler
}

func NewGRPCHandler(s Service) pb.UserServer {
	var opts []kitgrpc.ServerOption

	add := kitgrpc.NewServer(
		MakeAddEndpoint(s),
		decodeGRPCAddRequest,
		encodeGRPCAddResponse,
		opts...,
	)

	get := kitgrpc.NewServer(
		MakeGetEndpoint(s),
		decodeGRPCGetRequest,
		encodeGRPCGetResponse,
		opts...,
	)

	checkPassword := kitgrpc.NewServer(
		MakeCheckPasswordEndpoint(s),
		decodeGRPCCheckPasswordRequest,
		encodeGRPCCheckPasswordResponse,
		opts...,
	)

	return &grpcHandler{
		add:           add,
		get:           get,
		checkPassword: checkPassword,
	}
}

func (s *grpcHandler) Add(ctx context.Context, req *pb.AddReq) (*pb.AddResp, error) {
	_, rep, err := s.add.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.AddResp), nil
}

func (s *grpcHandler) Get(ctx context.Context, req *pb.GetReq) (*pb.GetResp, error) {
	_, rep, err := s.get.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetResp), nil
}

func (s *grpcHandler) CheckPassword(ctx context.Context, req *pb.CheckPasswordReq) (*pb.CheckPasswordResp, error) {
	_, rep, err := s.checkPassword.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CheckPasswordResp), nil
}

func decodeGRPCAddRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.AddReq)
	return AddRequest{Name: req.Name, Username: req.Username, Password: req.Password}, nil
}

func encodeGRPCAddResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(AddResponse)
	return &pb.AddResp{Id: resp.Id, Err: cgrpc.ErrorToStr(resp.Err)}, nil
}

func EncodeGRPCAddRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(AddRequest)
	return &pb.AddReq{Name: req.Name, Username: req.Username, Password: req.Password}, nil
}

func DecodeGRPCAddResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	resp := grpcReply.(*pb.AddResp)
	return AddResponse{Id: resp.Id, Err: cgrpc.ErrorFromStr(resp.Err)}, nil
}

func decodeGRPCGetRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetReq)
	return GetRequest{Username: req.Username}, nil
}

func encodeGRPCGetResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(GetResponse)
	return &pb.GetResp{User: &pb.UserModel{
		Id:       resp.User.Id,
		Name:     resp.User.Name,
		Username: resp.User.Username,
	}, Err: cgrpc.ErrorToStr(resp.Err)}, nil
}

func EncodeGRPCGetRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(GetRequest)
	return &pb.GetReq{Username: req.Username}, nil
}

func DecodeGRPCGetResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	resp := grpcReply.(*pb.GetResp)
	var u model.User
	if resp.User != nil {
		u.Id = resp.User.Id
		u.Name = resp.User.Name
		u.Username = resp.User.Username
	}
	return GetResponse{User: u, Err: cgrpc.ErrorFromStr(resp.Err)}, nil
}

func decodeGRPCCheckPasswordRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CheckPasswordReq)
	return CheckPasswordRequest{Username: req.Username, Password: req.Password}, nil
}

func encodeGRPCCheckPasswordResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(CheckPasswordResponse)
	return &pb.CheckPasswordResp{Id: resp.Id, Err: cgrpc.ErrorToStr(resp.Err)}, nil
}

func EncodeGRPCCheckPasswordRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(CheckPasswordRequest)
	return &pb.CheckPasswordReq{Username: req.Username, Password: req.Password}, nil
}

func DecodeGRPCCheckPasswordResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	resp := grpcReply.(*pb.CheckPasswordResp)
	return CheckPasswordResponse{Id: resp.Id, Err: cgrpc.ErrorFromStr(resp.Err)}, nil
}
