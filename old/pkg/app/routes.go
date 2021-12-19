package app

import "github.com/gin-gonic/gin"

func (s *Server) Routes() *gin.Engine {
	router := s.router

	// group all routes under api/api
	v1 := router.Group("api/v1")
	{
		v1.GET("/status/", s.ApiStatus())
		// prefix the user routes
		user := v1.Group("/users/")
		{
			user.POST("register/", s.ApiStatus())
			user.POST("login/", s.ApiStatus())
		}

		// prefix the movies routes
		movies := v1.Group("/movies/")
		{
			movies.GET("", s.ApiStatus())
			movies.GET("sort/:by/:direction/", s.ApiStatus())
			movies.GET("watched/", s.ApiStatus())
			movies.POST("add/", s.ApiStatus())
		}
		// I tried to add this group under movies group but there seem to be a limitation in gin:
		// https://stackoverflow.com/questions/55973623/how-can-i-have-a-parameter-with-slashes-in-gin-url
		movie := v1.Group("movie/")
		{
			movie.GET("info/:id/", s.ApiStatus())
			movie.DELETE("info/:id/", s.ApiStatus())
			movie.GET("watch/:id/", s.ApiStatus())
			movie.POST("review/:id/", s.ApiStatus())
		}

	}

	return router
}
