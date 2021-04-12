package mongo

import (
	"context"
	"sample/user/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepo struct {
	c *mongo.Collection
}

func NewUserRepo(db *mongo.Database) (*userRepo, error) {
	c := db.Collection("User")
	c.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{"username", 1}},
			Options: options.Index().SetBackground(true).SetUnique(true),
		},
	})
	return &userRepo{
		c: c,
	}, nil
}

func (repo *userRepo) Add(ctx context.Context, a model.User) (err error) {
	_, err = repo.c.InsertOne(ctx, a)
	if mongo.IsDuplicateKeyError(err) {
		err = model.ErrUserAlreadyExists
	}
	return
}

func (repo *userRepo) Get(ctx context.Context, id string) (a model.User, err error) {
	err = repo.c.FindOne(ctx, bson.M{"_id": id}).Decode(&a)
	if err == mongo.ErrNoDocuments {
		err = model.ErrUserNotFound
	}
	return
}

func (repo *userRepo) Get1(ctx context.Context, username string) (a model.User, err error) {
	err = repo.c.FindOne(ctx, bson.M{"username": username}).Decode(&a)
	if err == mongo.ErrNoDocuments {
		err = model.ErrUserNotFound
	}
	return
}
