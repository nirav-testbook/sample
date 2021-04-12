package user

import (
	"context"
	"sample/user/model"

	"github.com/go-kit/kit/endpoint"
)

type AddEndpoint endpoint.Endpoint
type GetEndpoint endpoint.Endpoint
type CheckPasswordEndpoint endpoint.Endpoint

type Endpoint struct {
	AddEndpoint
	GetEndpoint
	CheckPasswordEndpoint
}

type addRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type addResponse struct {
	Id string `json:"id"`
}

func MakeAddEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addRequest)
		id, err := s.Add(ctx, req.Name, req.Username, req.Password)
		return addResponse{Id: id}, err
	}
}

func (e AddEndpoint) Add(ctx context.Context, name string, username string, password string) (id string, err error) {
	request := addRequest{
		Name:     name,
		Username: username,
		Password: password,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	resp := response.(addResponse)
	return resp.Id, nil
}

type getRequest struct {
	Username string `schema:"username"`
}

type getResponse struct {
	User model.User `json:"user"`
}

func MakeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getRequest)
		user, err := s.Get(ctx, req.Username)
		return getResponse{User: user}, err
	}
}

func (e GetEndpoint) Get(ctx context.Context, username string) (user model.User, err error) {
	request := getRequest{
		Username: username,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	resp := response.(getResponse)
	return resp.User, nil
}

type checkPasswordRequest struct {
	Username string `schema:"username"`
	Password string `schema:"password"`
}

type checkPasswordResponse struct {
	Id string `json:"id"`
}

func MakeCheckPasswordEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(checkPasswordRequest)
		id, err := s.CheckPassword(ctx, req.Username, req.Password)
		return checkPasswordResponse{Id: id}, err
	}
}

func (e CheckPasswordEndpoint) CheckPassword(ctx context.Context, username string, password string) (id string, err error) {
	request := checkPasswordRequest{
		Username: username,
		Password: password,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	resp := response.(checkPasswordResponse)
	return resp.Id, nil
}
