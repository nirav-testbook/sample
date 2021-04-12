package user

import (
	"context"
	"time"

	"sample/user/model"

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

func (s *logSvc) Add(ctx context.Context, name string, username string, password string) (id string, err error) {
	defer func(t time.Time) {
		s.logger.Log(
			"ts", t,
			"method", "Add",
			"name", name,
			"username", username,
			"id", id,
			"time", time.Since(t),
			"err", err,
		)
	}(time.Now())
	return s.Service.Add(ctx, name, username, password)
}

func (s *logSvc) Get(ctx context.Context, username string) (user model.User, err error) {
	defer func(t time.Time) {
		s.logger.Log(
			"ts", t,
			"method", "Get",
			"username", username,
			"took", time.Since(t),
			"err", err,
		)
	}(time.Now())
	return s.Service.Get(ctx, username)
}

func (s *logSvc) CheckPassword(ctx context.Context, username string, password string) (id string, err error) {
	defer func(t time.Time) {
		s.logger.Log(
			"ts", t,
			"method", "CheckPassword",
			"username", username,
			"time", time.Since(t),
			"err", err,
		)
	}(time.Now())
	return s.Service.CheckPassword(ctx, username, password)
}
