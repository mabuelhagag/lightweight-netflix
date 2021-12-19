package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"lightweight-netflix/pkg/api"
	"log"
)

type Storage interface {
	CreateUser(request api.NewUserRequest) error
}

type storage struct {
	db *mongo.Database
}

type user struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	FullName string             `bson:"fullName,omitempty"`
	Email    string             `bson:"email,omitempty"`
	Age      uint8              `bson:"age,omitempty"`
	Password string             `bson:"password,omitempty"`
}

func NewStorage(db *mongo.Database) Storage {
	return &storage{db: db}
}

func (s *storage) CreateUser(request api.NewUserRequest) error {
	collection := s.db.Collection("users")
	user := user{FullName: request.FullName, Email: request.Email, Age: request.Age, Password: request.Password}
	insertResult, err := collection.InsertOne(context.TODO(), user)

	if err != nil {
		log.Printf("this was the error: %v", err.Error())
		return err
	}
	fmt.Println(insertResult.InsertedID)

	return nil
}
