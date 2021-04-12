package auth

import (
	"context"
	"time"

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

func (s *logSvc) Signin(ctx context.Context, username string, password string) (token string, err error) {
	defer func(t time.Time) {
		s.logger.Log(
			"ts", t,
			"method", "Signin",
			"username", username,
			"took", time.Since(t),
			"err", err,
		)
	}(time.Now())
	return s.Service.Signin(ctx, username, password)
}

func (s *logSvc) VerifyToken(ctx context.Context, token string) (userId string, err error) {
	defer func(t time.Time) {
		s.logger.Log(
			"ts", t,
			"method", "VerifyToken",
			"time", time.Since(t),
			"err", err,
		)
	}(time.Now())
	return s.Service.VerifyToken(ctx, token)
}
