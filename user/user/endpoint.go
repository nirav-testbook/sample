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

type AddRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AddResponse struct {
	Id  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

func (r AddResponse) Error() error {return r.Err}

func MakeAddEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddRequest)
		id, err := s.Add(ctx, req.Name, req.Username, req.Password)
		return AddResponse{Id: id, Err: err}, nil
	}
}

func (e AddEndpoint) Add(ctx context.Context, name string, username string, password string) (id string, err error) {
	request := AddRequest{
		Name:     name,
		Username: username,
		Password: password,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	resp := response.(AddResponse)
	return resp.Id, resp.Err
}

type GetRequest struct {
	Username string `schema:"username"`
}

type GetResponse struct {
	User model.User `json:"user"`
	Err  error      `json:"error,omitempty"`
}

func (r GetResponse) Error() error {return r.Err}

func MakeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRequest)
		user, err := s.Get(ctx, req.Username)
		return GetResponse{User: user, Err: err}, nil
	}
}

func (e GetEndpoint) Get(ctx context.Context, username string) (user model.User, err error) {
	request := GetRequest{
		Username: username,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	resp := response.(GetResponse)
	return resp.User, resp.Err
}

type CheckPasswordRequest struct {
	Username string `schema:"username"`
	Password string `schema:"password"`
}

type CheckPasswordResponse struct {
	Id  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

func (r CheckPasswordResponse) Error() error {return r.Err}

func MakeCheckPasswordEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CheckPasswordRequest)
		id, err := s.CheckPassword(ctx, req.Username, req.Password)
		return CheckPasswordResponse{Id: id, Err: err}, nil
	}
}

func (e CheckPasswordEndpoint) CheckPassword(ctx context.Context, username string, password string) (id string, err error) {
	request := CheckPasswordRequest{
		Username: username,
		Password: password,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	resp := response.(CheckPasswordResponse)
	return resp.Id, resp.Err
}
