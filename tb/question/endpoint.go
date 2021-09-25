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

type AddRequest struct {
	Question           model.Question `json:"question"`
	CorrectOptionIndex int            `json:"correct_option_index"`
}

type AddResponse struct {
	Id string `json:"id"`
}

func MakeAddEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddRequest)
		id, err := s.Add(ctx, req.Question, req.CorrectOptionIndex)
		return AddResponse{Id: id}, err
	}
}

func (e AddEndpoint) Add(ctx context.Context, q model.Question, correctOptionIndex int) (qid string, err error) {
	request := AddRequest{
		Question:           q,
		CorrectOptionIndex: correctOptionIndex,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	resp := response.(AddResponse)
	return resp.Id, nil
}

type GetRequest struct {
	Id string `schema:"id"`
}

type GetResponse struct {
	Question model.Question `json:"question"`
}

func MakeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRequest)
		question, err := s.Get(ctx, req.Id)
		return GetResponse{Question: question}, err
	}
}

func (e GetEndpoint) Get(ctx context.Context, id string) (question model.Question, err error) {
	request := GetRequest{
		Id: id,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	resp := response.(GetResponse)
	return resp.Question, nil
}

type ListRequest struct {
	Ids []string `schema:"ids"`
}

type ListResponse struct {
	Questions []model.Question `json:"questions"`
}

func MakeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListRequest)
		questions, err := s.List(ctx, req.Ids)
		return ListResponse{Questions: questions}, err
	}
}

func (e ListEndpoint) List(ctx context.Context, ids []string) (questions []model.Question, err error) {
	request := ListRequest{
		Ids: ids,
	}
	response, err := e(ctx, request)
	if err != nil {
		return
	}
	resp := response.(ListResponse)
	return resp.Questions, nil
}
