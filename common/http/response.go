package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	kithttp "github.com/go-kit/kit/transport/http"
)

type Errorer interface {
	Error() error
}

type ErrResp struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
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
		Error:   err.Error(),
	})
}

func EncodeJsonResp(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if e, ok := response.(Errorer); ok && e.Error() != nil {
		EncodeError(ctx, e.Error(), w)
		return nil
	}
	return json.NewEncoder(w).Encode(Resp{
		Success: true,
		Data:    response,
	})
}

func DecodeJsonRespOf(resp interface{}) kithttp.DecodeResponseFunc {
	return func(ctx context.Context, r *http.Response) (interface{}, error) {
		obj := reflect.New(reflect.TypeOf(resp))
		err := decodeResponse(ctx, r, obj.Interface())
		return obj.Elem().Interface(), err
	}
}

func decodeResponse(ctx context.Context, r *http.Response, resp interface{}) error {
	enc := json.NewDecoder(r.Body)
	if r.StatusCode != 200 {
		var e ErrResp
		err := enc.Decode(&e)
		if err != nil {
			return err
		}
		return errors.New(e.Error)
	}
	var data EncodeResp
	err := enc.Decode(&data)
	if err != nil {
		return err
	}
	return json.Unmarshal(data.Data, &resp)
}
