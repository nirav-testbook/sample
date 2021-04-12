package model

import (
	"context"
	"errors"
)

var (
	ErrUserNotFound      = errors.New("User not found")
	ErrUserAlreadyExists = errors.New("User already exists")
)

type User struct {
	Id       string `json:"id" bson:"_id"`
	Name     string `json:"name" bson:"name"`
	Username string `json:"username" bson:"username"`
	Password string `json:"-" bson:"password"`
}

type UserRepo interface {
	Add(ctx context.Context, acc User) (err error)
	Get(ctx context.Context, id string) (user User, err error)
	Get1(ctx context.Context, username string) (user User, err error)
}
