package question

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

	list := kithttp.NewServer(
		MakeListEndpoint(s),
		DecodeListRequest,
		chttp.EncodeJSONResponse,
		opts...,
	)

	r.Handle("/question", add).Methods(http.MethodPost)
	r.Handle("/question", get).Methods(http.MethodGet)
	r.Handle("/question/list", list).Methods(http.MethodGet)

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

func DecodeListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req listRequest
	err := schema.NewDecoder().Decode(&req, r.URL.Query())
	return req, err
}

func DecodeListResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var resp listResponse
	err := chttp.DecodeResponse(ctx, r, &resp)
	return resp, err
}
