package userrepo

import (
	"context"
	"go-app/domain/user"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// Repo Interface
type Repo interface {
	CreateUser(user *user.User) (*user.User, error)
}

type userRepo struct {
	db *mongo.Client
}

// NewUserRepo will instantiate User Repository
func NewUserRepo(db *mongo.Client) Repo {
	return &userRepo{
		db: db,
	}
}

func (b *userRepo) CreateUser(user *user.User) (*user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // releases resources if CreateUser completes before timeout elapses
	collection := b.db.Database("lt-netflix").Collection("users")
	_, err := collection.InsertOne(ctx, *user)

	if err != nil {
		panic(err)
	}
	return user, nil
}
