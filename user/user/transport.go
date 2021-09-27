package user

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

	checkPassword := kithttp.NewServer(
		MakeCheckPasswordEndpoint(s),
		chttp.DecodeJsonReqOf(CheckPasswordRequest{}),
		chttp.EncodeJsonResp,
		opts...,
	)

	r.Handle("/user", add).Methods(http.MethodPost)
	r.Handle("/user", get).Methods(http.MethodGet)
	r.Handle("/user/check", checkPassword).Methods(http.MethodPost)

	return r
}
