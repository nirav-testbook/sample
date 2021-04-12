package mongo

import (
	"context"
	"sample/auth/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type sessionRepo struct {
	c *mongo.Collection
}

func NewSessionRepo(db *mongo.Database) (*sessionRepo, error) {
	c := db.Collection("Session")
	return &sessionRepo{
		c: c,
	}, nil
}

func (repo *sessionRepo) Add(ctx context.Context, s model.Session) (err error) {
	_, err = repo.c.InsertOne(ctx, s)
	if mongo.IsDuplicateKeyError(err) {
		err = model.ErrSessionAlreadyExists
	}
	return
}

func (repo *sessionRepo) Get(ctx context.Context, token string) (l model.Session, err error) {
	err = repo.c.FindOne(ctx, bson.M{"_id": token}).Decode(&l)
	if err == mongo.ErrNoDocuments {
		err = model.ErrSessionNotFound
	}
	return
}
