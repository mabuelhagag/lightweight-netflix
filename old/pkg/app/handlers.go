package app

import (
	"lightweight-netflix/pkg/api"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) ApiStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		response := map[string]string{
			"status": "success",
			"data":   "Lightweight Netflix API is running smoothly",
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var newUser api.NewUserRequest
		err := c.ShouldBindJSON(&newUser)

		if err != nil {
			log.Printf("handler error: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		err = s.userService.New(newUser)

		if err != nil {
			log.Printf("service error: %v", err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		response := map[string]string{
			"status": "success",
			"data":   "new user created",
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) CreateWeightEntry() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var newWeight api.NewWeightRequest
		err := c.ShouldBindJSON(&newWeight)

		if err != nil {
			log.Printf("handler error: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		err = s.weightService.New(newWeight)

		if err != nil {
			log.Printf("service error: %v", err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		response := map[string]string{
			"status": "success",
			"data":   "new user created",
		}

		c.JSON(http.StatusOK, response)
	}
}
