package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
)

func EncodeJSONRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func EncodeSchemaRequest(_ context.Context, r *http.Request, request interface{}) error {
	v := url.Values{}
	err := schema.NewEncoder().Encode(request, v)
	if err != nil {
		return err
	}
	r.URL.RawQuery = v.Encode()
	return nil
}

type ErrResp struct {
	Success bool   `json:"success"`
	Error string `json:"error"`
}

type Resp struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type EncodeResp struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
}

func EncodeError(ctx context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(ErrResp{
		Success: false,
		Error: err.Error(),
	})
}

func EncodeJSONResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(Resp{
		Success: true,
		Data:    response,
	})
}

func DecodeResponse(ctx context.Context, r *http.Response, resp interface{}) error {
	enc := json.NewDecoder(r.Body)
	if r.StatusCode != 200 {
		var e ErrResp
		err := enc.Decode(&e)
		if err != nil {
			return err
		}
		return errors.New(e.Error)
	}
	var er EncodeResp
	err := enc.Decode(&er)
	if err != nil {
		return err
	}
	return json.Unmarshal(er.Data, &resp)
}
