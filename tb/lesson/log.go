package lesson

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

func (s *logSvc) Add(ctx context.Context, name string, qids []string) (id string, err error) {
	defer func(t time.Time) {
		s.logger.Log(
			"ts", t,
			"method", "Add",
			"name", name,
			"qids", qids,
			"id", id,
			"took", time.Since(t),
			"err", err,
		)
	}(time.Now())
	return s.Service.Add(ctx, name, qids)
}

func (s *logSvc) Get(ctx context.Context, id string) (lesson model.Lesson, err error) {
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

func (s *logSvc) Get1(ctx context.Context, id string) (lesson GetLessonRes, err error) {
	defer func(t time.Time) {
		s.logger.Log(
			"ts", t,
			"method", "Get1",
			"id", id,
			"took", time.Since(t),
			"err", err,
		)
	}(time.Now())
	return s.Service.Get1(ctx, id)
}
