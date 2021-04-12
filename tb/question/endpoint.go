package question

import (
	"context"

	"sample/tb/model"

	"github.com/go-kit/kit/endpoint"
)

type AddEndpoint endpoint.Endpoint
type GetEndpoint endpoint.Endpoint
type ListEndpoint endpoint.Endpoint

type Endpoint struct {
	AddEndpoint
	GetEndpoint
	ListEndpoint
}

type addRequest struct {
	Question           model.Question `json:"question"`
	CorrectOptionIndex int            `json:"correct_option_index"`
}

type addResponse struct {
	Id string `json:"id"`
}

func MakeAddEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addRequest)
		id, err := s.Add(ctx, req.Question, req.CorrectOptionIndex)
		return addResponse{Id: id}, err
	}
}

func (e AddEndpoint) Add(ctx context.Context, q model.Question, correctOptionIndex int) (qid string, err error) {
	request := addRequest{
		Question:           q,
		CorrectOptionIndex: correctOptionIndex,
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
	Question model.Question `json:"question"`
}

func MakeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getRequest)
		question, err := s.Get(ctx, req.Id)
		return getResponse{Question: question}, err
	}
}

func (e GetEndpoint) Get(ctx context.Context, id string) (question model.Question, err error) {
	request := getRequest{
		Id: id,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	resp := response.(getResponse)
	return resp.Question, nil
}

type listRequest struct {
	Ids []string `schema:"ids"`
}

type listResponse struct {
	Questions []model.Question `json:"questions"`
}

func MakeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listRequest)
		questions, err := s.List(ctx, req.Ids)
		return listResponse{Questions: questions}, err
	}
}

func (e ListEndpoint) List(ctx context.Context, ids []string) (questions []model.Question, err error) {
	request := listRequest{
		Ids: ids,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	resp := response.(listResponse)
	return resp.Questions, nil
}
