package lesson

import (
	"net/http"

	"sample/common/auth/token"
	chttp "sample/common/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHandler(s Service) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerBefore(token.HTTPTokenToContext),
		kithttp.ServerErrorEncoder(chttp.EncodeError),
	}

	add := kithttp.NewServer(
		MakeAddEndpoint(s),
		chttp.DecodeJsonReqOf(AddRequest{}),
		chttp.EncodeJsonResp,
		opts...,
	)

	get := kithttp.NewServer(
		MakeGetEndpoint(s),
		chttp.DecodeQueryReqOf(GetRequest{}),
		chttp.EncodeJsonResp,
		opts...,
	)

	get1 := kithttp.NewServer(
		MakeGet1Endpoint(s),
		chttp.DecodeQueryReqOf(Get1Request{}),
		chttp.EncodeJsonResp,
		opts...,
	)

	list := kithttp.NewServer(
		MakeListEndpoint(s),
		chttp.DecodeQueryReqOf(ListRequest{}),
		chttp.EncodeJsonResp,
		opts...,
	)

	r.Handle("/lesson", add).Methods(http.MethodPost)
	r.Handle("/lesson", get).Methods(http.MethodGet)
	r.Handle("/lesson/1", get1).Methods(http.MethodGet)
	r.Handle("/lesson/all", list).Methods(http.MethodGet)

	return r
}
