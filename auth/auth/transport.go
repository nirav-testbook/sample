package auth

import (
	"context"
	"encoding/json"
	"net/http"

	chttp "sample/common/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func NewHandler(s Service) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(chttp.EncodeError),
	}

	signin := kithttp.NewServer(
		MakeSigninEndpoint(s),
		DecodeSigninRequest,
		chttp.EncodeJsonResp,
		opts...,
	)

	verifyToken := kithttp.NewServer(
		MakeVerifyTokenEndpoint(s),
		DecodeVerifyTokenRequest,
		chttp.EncodeJsonResp,
		opts...,
	)

	r.Handle("/auth/signin", signin).Methods(http.MethodPost)
	r.Handle("/auth/verify", verifyToken).Methods(http.MethodGet)

	return r
}

func DecodeSigninRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req signinRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func DecodeSigninResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var resp signinResponse
	err := chttp.DecodeResponse(ctx, r, &resp)
	return resp, err
}

func DecodeVerifyTokenRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req verifyTokenRequest
	err := schema.NewDecoder().Decode(&req, r.URL.Query())
	return req, err
}

func DecodeVerifyTokenResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var resp verifyTokenResponse
	err := chttp.DecodeResponse(ctx, r, &resp)
	return resp, err
}
