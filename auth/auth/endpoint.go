package auth

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type SigninEndpoint endpoint.Endpoint
type VerifyTokenEndpoint endpoint.Endpoint

type Endpoint struct {
	SigninEndpoint
	VerifyTokenEndpoint
}

type SigninRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SigninResponse struct {
	Token string `json:"token"`
	Err   error  `json:"error,omitempty"`
}

func (r SigninResponse) Error() error {return r.Err}

func MakeSigninEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SigninRequest)
		token, err := s.Signin(ctx, req.Username, req.Password)
		return SigninResponse{Token: token, Err: err}, nil
	}
}

func (e SigninEndpoint) Signin(ctx context.Context, username string, password string) (token string, err error) {
	request := SigninRequest{
		Username: username,
		Password: password,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	resp := response.(SigninResponse)
	return resp.Token, resp.Err
}

type VerifyTokenRequest struct {
	Token string `schema:"token"`
}

type VerifyTokenResponse struct {
	UserID string `json:"user_id"`
	Err    error  `json:"error,omitempty"`
}

func (r VerifyTokenResponse) Error() error {return r.Err}

func MakeVerifyTokenEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(VerifyTokenRequest)
		userID, err := s.VerifyToken(ctx, req.Token)
		return VerifyTokenResponse{UserID: userID, Err: err}, nil
	}
}

func (e VerifyTokenEndpoint) VerifyToken(ctx context.Context, token string) (userID string, err error) {
	request := VerifyTokenRequest{
		Token: token,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	resp := response.(VerifyTokenResponse)
	return resp.UserID, resp.Err
}
