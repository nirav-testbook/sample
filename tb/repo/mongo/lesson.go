package mongo

import (
	"context"
	"sample/tb/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type lessonRepo struct {
	c *mongo.Collection
}

func NewLessonRepo(db *mongo.Database) (*lessonRepo, error) {
	c := db.Collection("Lesson")
	return &lessonRepo{
		c: c,
	}, nil
}

func (repo *lessonRepo) Add(ctx context.Context, l model.Lesson) (err error) {
	_, err = repo.c.InsertOne(ctx, l)
	if mongo.IsDuplicateKeyError(err) {
		err = model.ErrLessonAlreadyExists
	}
	return
}

func (repo *lessonRepo) Get(ctx context.Context, id string) (l model.Lesson, err error) {
	err = repo.c.FindOne(ctx, bson.M{"_id": id}).Decode(&l)
	if err == mongo.ErrNoDocuments {
		err = model.ErrLessonNotFound
	}
	return
}
