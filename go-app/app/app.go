package app

import (
	"context"
	"go-app/configs"
	"go-app/repositories/userrepo"
	"go-app/services/userservice"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-app/controllers"
)

var (
	r = gin.Default()
)

// Run is the App Entry Point
func Run() {

	/*
		====== Setup configs ============
	*/
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	config := configs.GetConfig()

	// Set client options
	clientOptions := options.Client().ApplyURI(config.MongoDB.URI) // use env variables
	// Connect to MongoDB
	mongoDB, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		panic(err)
	}

	/*
		====== Setup repositories =======
	*/
	userRepo := userrepo.NewUserRepo(mongoDB)
	/*
		====== Setup services ===========
	*/
	userService := userservice.NewUserService(userRepo)
	/*
		====== Setup controllers ========
	*/
	userCtl := controllers.NewUserController(userService)

	/*
		======== Routes ============
	*/

	// API Home
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to your App on Docker",
		})
	})

	/*
		===== User Routes =====
	*/
	r.POST("/users", userCtl.RegisterUser)

	err = r.Run()
	if err != nil {
		panic(err)
	}
}
