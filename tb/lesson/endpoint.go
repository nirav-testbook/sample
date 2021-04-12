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

type addRequest struct {
	Name        string   `json:"name"`
	QuestionIds []string `json:"question_ids"`
}

type addResponse struct {
	Id string `json:"id"`
}

func MakeAddEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addRequest)
		id, err := s.Add(ctx, req.Name, req.QuestionIds)
		return addResponse{Id: id}, err
	}
}

func (e AddEndpoint) Add(ctx context.Context, name string, qids []string) (id string, err error) {
	request := addRequest{
		Name:        name,
		QuestionIds: qids,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	resp := response.(addResponse)
	return resp.Id, nil
}

type getRequest struct {
	Id string `schema:"id"`
}

type getResponse struct {
	Lesson model.Lesson `json:"lesson"`
}

func MakeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getRequest)
		lesson, err := s.Get(ctx, req.Id)
		return getResponse{Lesson: lesson}, err
	}
}

func (e GetEndpoint) Get(ctx context.Context, id string) (lesson model.Lesson, err error) {
	request := getRequest{
		Id: id,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	resp := response.(getResponse)
	return resp.Lesson, nil
}

type get1Request struct {
	Id string `schema:"id"`
}

type get1Response struct {
	Lesson GetLessonRes `json:"lesson"`
}

func MakeGet1Endpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(get1Request)
		lesson, err := s.Get1(ctx, req.Id)
		return get1Response{Lesson: lesson}, err
	}
}

func (e Get1Endpoint) Get1(ctx context.Context, id string) (lesson GetLessonRes, err error) {
	request := get1Request{
		Id: id,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	resp := response.(get1Response)
	return resp.Lesson, nil
}
