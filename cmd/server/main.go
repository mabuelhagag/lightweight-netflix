package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"lightweight-netflix/pkg/api"
	"lightweight-netflix/pkg/app"
	"lightweight-netflix/pkg/repository"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "this is the startup error: %s\n", err)
		os.Exit(1)
	}
}

// func run will be responsible for setting up db connections, routers etc
func run() error {
	// setup database connection
	db, err := setupDatabase()
	if err != nil {
		return err
	}

	// create storage dependency
	storage := repository.NewStorage(db)

	if err != nil {
		return err
	}

	// create router dependency
	router := gin.Default()
	router.Use(cors.Default())

	// create user service
	userService := api.NewUserService(storage)

	// create weight service
	weightService := api.NewWeightService(storage)

	server := app.NewServer(router, userService, weightService)

	// start the server
	err = server.Run()

	if err != nil {
		return err
	}

	return nil
}

func setupDatabase() (*mongo.Database, error) {
	DBUsername := os.Getenv("MONGO_INITDB_DATABASE")
	DBPassword := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	DB := os.Getenv("MONGO_INITDB_DATABASE")

	uri := fmt.Sprintf("mongodb://%s:%s@db:27017", DBUsername, DBPassword)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	db := client.Database(DB)
	return db, err
}
