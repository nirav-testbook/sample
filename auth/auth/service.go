package auth

import (
	"context"
	"time"

	"sample/auth/model"
	"sample/common/err"
	"sample/common/sid"
	"sample/user/user"
)

var (
	sessionDuration = 365 * 24 * time.Hour
)

var (
	errInvalidArgument = err.New(301, "invalid argument")
	errInvalidPassword = err.New(302, "invalid password")
	errInvalidToken    = err.New(303, "invalid token")
)

type Service interface {
	Signin(ctx context.Context, username string, password string) (token string, err error)
	VerifyToken(ctx context.Context, token string) (userID string, err error)
}

type service struct {
	sessionRepo model.SessionRepo
	userSvc     user.Service
}

func NewService(sessionRepo model.SessionRepo, userSvc user.Service) Service {
	return &service{
		sessionRepo: sessionRepo,
		userSvc:     userSvc,
	}
}

func (s *service) Signin(ctx context.Context, username string, password string) (token string, err error) {
	if len(username) < 1 || len(password) < 1 {
		err = errInvalidArgument
		return
	}

	uid, err := s.userSvc.CheckPassword(ctx, username, password)
	if err != nil {
		return
	}

	t, err := sid.New(24)
	if err != nil {
		return
	}

	session := model.Session{
		Token:      t,
		UserId:     uid,
		ExpiryDate: time.Now().Add(sessionDuration),
	}

	err = s.sessionRepo.Add(ctx, session)
	if err != nil {
		return
	}

	return t, nil
}

func (s *service) VerifyToken(ctx context.Context, token string) (userId string, err error) {
	if len(token) < 1 {
		err = errInvalidArgument
		return
	}

	session, err := s.sessionRepo.Get(ctx, token)
	if err != nil {
		return "", errInvalidToken
	}

	if session.ExpiryDate.Before(time.Now()) {
		return "", errInvalidToken
	}

	return session.UserId, nil
}
