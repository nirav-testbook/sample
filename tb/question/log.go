package question

import (
	"context"
	"time"

	"sample/tb/model"

	"github.com/go-kit/kit/log"
)

type logSvc struct {
	Service
	logger log.Logger
}

func NewLogService(s Service, logger log.Logger) Service {
	return &logSvc{
		Service: s,
		logger:  logger,
	}
}

func (s *logSvc) Add(ctx context.Context, q model.Question, correctOptionIndex int) (qid string, err error) {
	defer func(t time.Time) {
		s.logger.Log(
			"ts", t,
			"method", "Add",
			"question", q,
			"correct_option_index", correctOptionIndex,
			"took", time.Since(t),
			"err", err,
		)
	}(time.Now())
	return s.Service.Add(ctx, q, correctOptionIndex)
}

func (s *logSvc) Get(ctx context.Context, id string) (question model.Question, err error) {
	defer func(t time.Time) {
		s.logger.Log(
			"ts", t,
			"method", "Get",
			"id", id,
			"took", time.Since(t),
			"err", err,
		)
	}(time.Now())
	return s.Service.Get(ctx, id)
}

func (s *logSvc) List(ctx context.Context, ids []string) (questions []model.Question, err error) {
	defer func(t time.Time) {
		s.logger.Log(
			"ts", t,
			"method", "List",
			"ids", ids,
			"took", time.Since(t),
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx, ids)
}
