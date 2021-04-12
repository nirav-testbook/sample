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

type signinRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signinResponse struct {
	Token string `json:"token"`
}

func MakeSigninEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(signinRequest)
		token, err := s.Signin(ctx, req.Username, req.Password)
		return signinResponse{Token: token}, err
	}
}

func (e SigninEndpoint) Signin(ctx context.Context, username string, password string) (token string, err error) {
	request := signinRequest{
		Username: username,
		Password: password,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	resp := response.(signinResponse)
	return resp.Token, nil
}

type verifyTokenRequest struct {
	Token string `schema:"token"`
}

type verifyTokenResponse struct {
	UserID string `json:"user_id"`
}

func MakeVerifyTokenEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(verifyTokenRequest)
		userID, err := s.VerifyToken(ctx, req.Token)
		return verifyTokenResponse{UserID: userID}, err
	}
}

func (e VerifyTokenEndpoint) VerifyToken(ctx context.Context, token string) (userID string, err error) {
	request := verifyTokenRequest{
		Token: token,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	resp := response.(verifyTokenResponse)
	return resp.UserID, nil
}
