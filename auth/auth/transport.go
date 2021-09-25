package auth

import (
	"net/http"

	chttp "sample/common/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHandler(s Service) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(chttp.EncodeError),
	}

	signin := kithttp.NewServer(
		MakeSigninEndpoint(s),
		chttp.DecodeJsonReqOf(SigninRequest{}),
		chttp.EncodeJsonResp,
		opts...,
	)

	verifyToken := kithttp.NewServer(
		MakeVerifyTokenEndpoint(s),
		chttp.DecodeJsonReqOf(VerifyTokenRequest{}),
		chttp.EncodeJsonResp,
		opts...,
	)

	r.Handle("/auth/signin", signin).Methods(http.MethodPost)
	r.Handle("/auth/verify", verifyToken).Methods(http.MethodGet)

	return r
}
