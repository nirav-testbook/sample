package question

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

	list := kithttp.NewServer(
		MakeListEndpoint(s),
		chttp.DecodeQueryReqOf(ListRequest{}),
		chttp.EncodeJsonResp,
		opts...,
	)

	r.Handle("/question", add).Methods(http.MethodPost)
	r.Handle("/question", get).Methods(http.MethodGet)
	r.Handle("/question/list", list).Methods(http.MethodGet)

	return r
}
