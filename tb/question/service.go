package question

import (
	"context"

	"sample/common/err"
	"sample/common/id"
	"sample/tb/model"
)

var (
	errInvalidArgument = err.New(101, "invalid argument")
)

type Service interface {
	Add(ctx context.Context, q model.Question, correctOptionIndex int) (id string, err error)
	Get(ctx context.Context, id string) (question model.Question, err error)
	List(ctx context.Context, ids []string) (questions []model.Question, err error)
}

type service struct {
	questionRepo model.QuestionRepo
}

func NewService(questionRepo model.QuestionRepo) Service {
	return &service{
		questionRepo: questionRepo,
	}
}

func (s *service) Add(ctx context.Context, q model.Question, correctOptionIndex int) (qid string, err error) {
	if len(q.Text) < 1 || len(q.Options) < 1 {
		return "", errInvalidArgument
	}

	q.Id = id.New()
	for i := range q.Options {
		q.Options[i].Id = id.New()
		if i == correctOptionIndex {
			q.CorrectOptionId = q.Options[i].Id
		}
	}

	err = s.questionRepo.Add(ctx, q)
	if err != nil {
		return
	}

	return q.Id, nil
}

func (s *service) Get(ctx context.Context, id string) (q model.Question, err error) {
	if len(id) < 1 {
		err = errInvalidArgument
		return
	}
	return s.questionRepo.Get(ctx, id)
}

func (s *service) List(ctx context.Context, ids []string) (questions []model.Question, err error) {
	if len(ids) < 1 {
		err = errInvalidArgument
		return
	}
	return s.questionRepo.List(ctx, ids)
}
