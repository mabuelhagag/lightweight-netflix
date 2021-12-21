package userrepo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// Repo Interface
type Repo interface {
	CreateUser(user *User) (*User, error)
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

// User struct
type User struct {
	FullName string `bson:"name"`
	Age      uint8  `bson:"page_count"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

func (b *userRepo) CreateUser(user *User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // releases resources if CreateUser completes before timeout elapses
	collection := b.db.Database("lw-netflix").Collection("users")
	var result bson.M
	collection.FindOne(context.TODO(), bson.D{{"email", user.Email}}).Decode(&result)
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
