package model

import (
	"context"
	"errors"
	"time"
)

var (
	ErrSessionNotFound      = errors.New("Session not found")
	ErrSessionAlreadyExists = errors.New("Session already exists")
)

type Session struct {
	Token      string    `json:"token" bson:"_id"`
	UserId     string    `json:"user_id" bson:"user_id"`
	ExpiryDate time.Time `json:"expiry_date" bson:"expiry_date"`
}

type SessionRepo interface {
	Add(ctx context.Context, s Session) error
	Get(ctx context.Context, token string) (s Session, err error)
}
