package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/schema"
)

func EncodeQueryReq(_ context.Context, r *http.Request, request interface{}) error {
	v := url.Values{}
	err := schema.NewEncoder().Encode(request, v)
	if err != nil {
		return err
	}
	r.URL.RawQuery = v.Encode()
	return nil
}

func DecodeQueryReqOf(req interface{}) kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		obj := reflect.New(reflect.TypeOf(req))
		dec := schema.NewDecoder()
		dec.IgnoreUnknownKeys(true)
		err := dec.Decode(obj.Interface(), r.URL.Query())
		return obj.Elem().Interface(), err
	}
}

func EncodeJsonReq(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func DecodeJsonReqOf(req interface{}) kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		obj := reflect.New(reflect.TypeOf(req))
		err := json.NewDecoder(r.Body).Decode(obj.Interface())
		return obj.Elem().Interface(), err
	}
}
