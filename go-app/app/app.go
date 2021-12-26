package app

import (
	"context"
	"go-app/configs"
	"go-app/middlewares"
	"go-app/repositories/moviesrepo"
	"go-app/repositories/userrepo"
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
	moviesRepo := moviesrepo.NewMoviesRepo(mongoDB)

	/*
		====== Setup controllers ========
	*/
	userCtl := controllers.NewUserController(userRepo)
	moviesCtl := controllers.NewMoviesController(moviesRepo, userRepo)

	/*
		======== Routes ============
	*/
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	// API Home
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to your App on Docker",
		})
	})

	/*
		===== User Routes =====
	*/

	authorized := func(c *gin.Context) {
		controllers.HTTPRes(c, http.StatusOK, "Authorized", nil)
	}

	r.POST("/users/register/", userCtl.RegisterUser)
	r.POST("/users/login/", userCtl.LoginUser)
	movies := r.Group("/movies/").Use(middlewares.Authorize())
	{
		movies.GET("", authorized)
		movies.GET("sort/:by/:direction/", authorized)
		movies.GET("watched/", authorized)
	}
	movie := r.Group("/movie/").Use(middlewares.Authorize())
	{
		movie.POST("add/", moviesCtl.AddMovie)
		movie.GET("info/:id/", authorized)
		movie.PUT("info/:id/", moviesCtl.UploadCover)
		movie.POST("info/:id/", moviesCtl.UpdateMovie)
		movie.DELETE("info/:id/", moviesCtl.DeleteMovie)
		movie.GET("watch/:id/", moviesCtl.WatchMovie)
		movie.POST("review/:id/", authorized)
	}
	err = r.Run()
	if err != nil {
		panic(err)
	}
}
