package model

import (
	"context"
	"errors"
)

var (
	ErrQuestionNotFound      = errors.New("Question not found")
	ErrQuestionAlreadyExists = errors.New("Question already exists")
)

type Option struct {
	Id    string `json:"id" bson:"_id"`
	Text  string `json:"text" bson:"text"`
	Order int    `json:"order" bson:"order"`
}

type Question struct {
	Id              string   `json:"id" bson:"_id"`
	Text            string   `json:"text" bson:"text"`
	Options         []Option `json:"options" bson:"options"`
	CorrectOptionId string   `json:"correct_option_id" bson:"correct_option_id"`
}

type QuestionRepo interface {
	Add(ctx context.Context, q Question) (err error)
	Get(ctx context.Context, id string) (question Question, err error)
	List(ctx context.Context, ids []string) (questions []Question, err error)
}
