package lesson

import (
	"context"

	"sample/tb/model"

	"github.com/go-kit/kit/endpoint"
)

type AddEndpoint endpoint.Endpoint
type GetEndpoint endpoint.Endpoint
type Get1Endpoint endpoint.Endpoint

type Endpoint struct {
	AddEndpoint
	GetEndpoint
	Get1Endpoint
}

type AddRequest struct {
	Name        string   `json:"name"`
	QuestionIds []string `json:"question_ids"`
}

type AddResponse struct {
	Id  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

func (r AddResponse) Error() error {return r.Err}

func MakeAddEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddRequest)
		id, err := s.Add(ctx, req.Name, req.QuestionIds)
		return AddResponse{Id: id, Err: err}, nil
	}
}

func (e AddEndpoint) Add(ctx context.Context, name string, qids []string) (id string, err error) {
	request := AddRequest{
		Name:        name,
		QuestionIds: qids,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	resp := response.(AddResponse)
	return resp.Id, resp.Err
}

type GetRequest struct {
	Id string `schema:"id"`
}

type GetResponse struct {
	Lesson model.Lesson `json:"lesson"`
	Err    error        `json:"error,omitempty"`
}

func (r GetResponse) Error() error {return r.Err}

func MakeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRequest)
		lesson, err := s.Get(ctx, req.Id)
		return GetResponse{Lesson: lesson, Err: err}, nil
	}
}

func (e GetEndpoint) Get(ctx context.Context, id string) (lesson model.Lesson, err error) {
	request := GetRequest{
		Id: id,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	resp := response.(GetResponse)
	return resp.Lesson, resp.Err
}

type Get1Request struct {
	Id string `schema:"id"`
}

type Get1Response struct {
	Lesson GetLessonRes `json:"lesson"`
	Err    error        `json:"error,omitempty"`
}

func (r Get1Response) Error() error {return r.Err}

func MakeGet1Endpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(Get1Request)
		lesson, err := s.Get1(ctx, req.Id)
		return Get1Response{Lesson: lesson, Err: err}, nil
	}
}

func (e Get1Endpoint) Get1(ctx context.Context, id string) (lesson GetLessonRes, err error) {
	request := Get1Request{
		Id: id,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	resp := response.(Get1Response)
	return resp.Lesson, resp.Err
}
