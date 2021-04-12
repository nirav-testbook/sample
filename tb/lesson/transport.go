package lesson

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

	get1 := kithttp.NewServer(
		MakeGet1Endpoint(s),
		DecodeGet1Request,
		chttp.EncodeJSONResponse,
		opts...,
	)

	r.Handle("/lesson", add).Methods(http.MethodPost)
	r.Handle("/lesson", get).Methods(http.MethodGet)
	r.Handle("/lesson/1", get1).Methods(http.MethodGet)

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

func DecodeGet1Request(ctx context.Context, r *http.Request) (interface{}, error) {
	var req get1Request
	err := schema.NewDecoder().Decode(&req, r.URL.Query())
	return req, err
}

func DecodeGet1Response(ctx context.Context, r *http.Response) (interface{}, error) {
	var resp get1Response
	err := chttp.DecodeResponse(ctx, r, &resp)
	return resp, err
}
