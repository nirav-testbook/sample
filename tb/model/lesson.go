package model

import (
	"context"
	"errors"
)

var (
	ErrLessonNotFound      = errors.New("Lesson not found")
	ErrLessonAlreadyExists = errors.New("Lesson already exists")
)

type Lesson struct {
	Id          string   `json:"id" bson:"_id"`
	Name        string   `json:"title" bson:"name"`
	QuestionIds []string `json:"question_ids" bson:"question_ids"`
}

type LessonRepo interface {
	Add(ctx context.Context, lesson Lesson) (err error)
	Get(ctx context.Context, id string) (lesson Lesson, err error)
}
