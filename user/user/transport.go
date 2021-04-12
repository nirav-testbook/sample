package user

import (
	"context"
	"encoding/json"
	"net/http"

	"sample/common/auth/token"
	chttp "sample/common/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func NewHandler(s Service) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerBefore(token.HTTPTokenToContext),
		kithttp.ServerErrorEncoder(chttp.EncodeError),
	}

	add := kithttp.NewServer(
		MakeAddEndpoint(s),
		DecodeAddRequest,
		chttp.EncodeJSONResponse,
		opts...,
	)

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

	r.Handle("/user", add).Methods(http.MethodPost)
	r.Handle("/user", get).Methods(http.MethodGet)
	r.Handle("/user/check", checkPassword).Methods(http.MethodGet)

	return r
}

func DecodeAddRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req addRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func DecodeAddResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var resp addResponse
	err := chttp.DecodeResponse(ctx, r, &resp)
	return resp, err
}

func DecodeGetRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req getRequest
	err := schema.NewDecoder().Decode(&req, r.URL.Query())
	return req, err
}

func DecodeGetResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var resp getResponse
	err := chttp.DecodeResponse(ctx, r, &resp)
	return resp, err
}

func DecodeCheckPasswordRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req checkPasswordRequest
	err := schema.NewDecoder().Decode(&req, r.URL.Query())
	return req, err
}

func DecodeCheckPasswordResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var resp checkPasswordResponse
	err := chttp.DecodeResponse(ctx, r, &resp)
	return resp, err
}
