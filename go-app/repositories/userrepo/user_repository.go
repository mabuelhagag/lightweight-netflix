package userrepo

import (
	"context"
	"errors"
	"go-app/definitions/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Repo Interface
type Repo interface {
	CreateUser(user *user.User) (*user.User, error)
	CheckPassword(user *user.LoginInfoInput) error
	GetUser(email string) (*user.User, error)
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

func (b *userRepo) GetUser(email string) (*user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // releases resources if CreateUser completes before timeout elapses
	collection := b.db.Database("lw-netflix").Collection("users")
	var result *user.User
	if err := collection.FindOne(ctx, bson.D{{"email", email}}).Decode(&result); err != nil {
		return nil, errors.New("Unable to get user")
	}
	userInstance := user.User{
		ID:       result.ID,
		FullName: result.FullName,
		Age:      result.Age,
		Email:    result.Email,
	}
	return &userInstance, nil
}

func (b *userRepo) CreateUser(user *user.User) (*user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // releases resources if CreateUser completes before timeout elapses
	collection := b.db.Database("lw-netflix").Collection("users")
	var result bson.M
	if err := collection.FindOne(ctx, bson.D{{"email", user.Email}}).Decode(&result); err != nil {
		if err != nil {
			if err != mongo.ErrNoDocuments {
				return nil, errors.New("Unable to check for user existance")
			}
		}

	}
	if result != nil {
		return nil, errors.New("User already exists")
	}

	bytes, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(bytes)
	_, err := collection.InsertOne(ctx, *user)

	if err != nil {
		panic(err)
	}
	return user, nil
}

func (b *userRepo) CheckPassword(input *user.LoginInfoInput) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // releases resources if CreateUser completes before timeout elapses

	collection := b.db.Database("lw-netflix").Collection("users")

	var result user.User
	err = collection.FindOne(ctx, bson.D{{"email", input.Email}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("User not found")
		}
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(input.Password))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("Email and password combination is incorrect")
		}
	}
	return nil
}
