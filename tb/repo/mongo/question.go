package mongo

import (
	"context"
	"sample/tb/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type questionRepo struct {
	c *mongo.Collection
}

func NewQuestionRepo(db *mongo.Database) (*questionRepo, error) {
	c := db.Collection("Question")
	return &questionRepo{
		c: c,
	}, nil
}

func (repo *questionRepo) Add(ctx context.Context, a model.Question) (err error) {
	_, err = repo.c.InsertOne(ctx, a)
	if mongo.IsDuplicateKeyError(err) {
		err = model.ErrQuestionAlreadyExists
	}
	return
}

func (repo *questionRepo) Get(ctx context.Context, id string) (q model.Question, err error) {
	err = repo.c.FindOne(ctx, bson.M{"_id": id}).Decode(&q)
	if err == mongo.ErrNoDocuments {
		err = model.ErrQuestionNotFound
	}
	return
}

func (repo *questionRepo) List(ctx context.Context, ids []string) (qs []model.Question, err error) {
	c, err := repo.c.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return
	}
	defer c.Close(ctx)
	err = c.All(ctx, &qs)
	return
}
