package user

import (
	"context"
	"errors"

	"sample/common/id"
	"sample/user/model"

	"golang.org/x/crypto/bcrypt"
)

var (
	errInvalidArgument       = errors.New("Invalid argument")
	errInvalidPassword       = errors.New("Invalid password")
	errUsernameAlreadyExists = errors.New("Username already exists")
)

type Service interface {
	Add(ctx context.Context, name string, username string, password string) (id string, err error)
	Get(ctx context.Context, username string) (user model.User, err error)
	CheckPassword(ctx context.Context, username string, password string) (id string, err error)
}

type service struct {
	userRepo model.UserRepo
}

func NewService(userRepo model.UserRepo) Service {
	return &service{
		userRepo: userRepo,
	}
}

func (s *service) Add(ctx context.Context, name string, username string, password string) (userId string, err error) {
	if len(name) < 1 || len(username) < 1 || len(password) < 8 {
		return "", errInvalidArgument
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return
	}

	user := model.User{
		Id:       id.New(),
		Name:     name,
		Username: username,
		Password: string(hash),
	}

	err = s.userRepo.Add(ctx, user)
	if err != nil {
		if err == model.ErrUserAlreadyExists {
			err = errUsernameAlreadyExists
		}
		return
	}

	return user.Id, nil
}

func (s *service) Get(ctx context.Context, id string) (user model.User, err error) {
	if len(id) < 1 {
		err = errInvalidArgument
		return
	}
	return s.userRepo.Get(ctx, id)
}

func (s *service) CheckPassword(ctx context.Context, username string, password string) (id string, err error) {
	if len(username) < 1 {
		err = errInvalidArgument
		return
	}

	u, err := s.userRepo.Get1(ctx, username)
	if err != nil {
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			err = errInvalidPassword
		}
		return
	}
	return u.Id, nil

}
