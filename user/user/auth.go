package user

/*
import (
	"context"

	"sample/auth"
	"sample/common/auth/token"
	"sample/common/err"
	"sample/user/model"
)

var (
	errAccessTokenNotFound = err.New(1, "access token not found")
)

type authSvc struct {
	Service
	authService auth.Service
}

func NewAuthService(s Service, authService auth.Service) Service {
	return &authSvc{
		Service:     s,
		authService: authService,
	}
}

func (s *authSvc) verifyToken(ctx context.Context) (userID string, err error) {
	token, ok := ctx.Value(token.ContextKey).(string)
	if !ok || len(token) < 1 {
		return "", errAccessTokenNotFound
	}

	return s.authService.VerifyToken(ctx, token)
}

func (s *authSvc) Get(ctx context.Context, username string) (account model.Account, err error) {
	_, err = s.verifyToken(ctx)
	if err != nil {
		return
	}

	return s.Service.Get(ctx, username)
}

func (s *authSvc) Get1(ctx context.Context, id string) (account model.Account, err error) {
	_, err = s.verifyToken(ctx)
	if err != nil {
		return
	}

	return s.Service.Get1(ctx, id)
}

func (s *authSvc) Add(ctx context.Context, name string, username string, password string) (err error) {
	_, err = s.verifyToken(ctx)
	if err != nil {
		return
	}

	return s.Service.Add(ctx, name, username, password)
}

func (s *authSvc) Update(ctx context.Context, username string, name string) (err error) {
	_, err = s.verifyToken(ctx)
	if err != nil {
		return
	}

	return s.Service.Update(ctx, username, name)
}

func (s *authSvc) List(ctx context.Context) (account []model.Account, err error) {
	_, err = s.verifyToken(ctx)
	if err != nil {
		return
	}

	return s.Service.List(ctx)
}
*/
