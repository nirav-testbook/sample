package lesson

import (
	"context"

	"sample/common/err"
	"sample/common/id"
	"sample/tb/model"
)

var (
	errInvalidArgument = err.New(101, "invalid argument")
)

type Service interface {
	Add(ctx context.Context, name string, qids []string) (id string, err error)
	Get(ctx context.Context, id string) (lessons model.Lesson, err error)
	Get1(ctx context.Context, id string) (lesson GetLessonRes, err error)
	List(ctx context.Context) (lesson []model.Lesson, err error)
}

type service struct {
	lessonRepo   model.LessonRepo
	questionRepo model.QuestionRepo
}

func NewService(lessonRepo model.LessonRepo, questionRepo model.QuestionRepo) Service {
	return &service{
		lessonRepo:   lessonRepo,
		questionRepo: questionRepo,
	}
}

func (s *service) Add(ctx context.Context, name string, qids []string) (lessonId string, err error) {
	if len(name) < 1 || len(qids) < 1 {
		return "", errInvalidArgument
	}

	lesson := model.Lesson{
		Id:          id.New(),
		QuestionIds: qids,
	}

	err = s.lessonRepo.Add(ctx, lesson)
	if err != nil {
		return
	}

	return lesson.Id, nil
}

func (s *service) Get(ctx context.Context, id string) (lesson model.Lesson, err error) {
	if len(id) < 1 {
		err = errInvalidArgument
		return
	}

	return s.lessonRepo.Get(ctx, id)
}

func (s *service) Get1(ctx context.Context, id string) (lesson GetLessonRes, err error) {
	if len(id) < 1 {
		err = errInvalidArgument
		return
	}

	l, err := s.lessonRepo.Get(ctx, id)
	if err != nil {
		return
	}

	questions, err := s.questionRepo.List(ctx, l.QuestionIds)
	if err != nil {
		return
	}

	lesson = GetLessonRes{
		Id:   l.Id,
		Name: l.Name,
	}
	for i := range questions {
		lesson.Questions = append(lesson.Questions, GetLessonQuestionRes{
			Id:      questions[i].Id,
			Text:    questions[i].Text,
			Options: questions[i].Options,
		})
	}
	return
}

func (s *service) List(ctx context.Context) (lessons []model.Lesson, err error) {
	return s.lessonRepo.List(ctx)
}
