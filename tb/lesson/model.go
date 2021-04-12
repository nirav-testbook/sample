package lesson

import "sample/tb/model"

type GetLessonQuestionRes struct {
	Id      string         `json:"id" bson:"_id"`
	Text    string         `json:"text" bson:"text"`
	Options []model.Option `json:"options" bson:"options"`
}

type GetLessonRes struct {
	Id        string                 `json:"id"`
	Name      string                 `json:"name"`
	Questions []GetLessonQuestionRes `json:"questions"`
}
