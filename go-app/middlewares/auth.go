package middlewares

import (
	"go-app/controllers"

	"go-app/definitions/user"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Authorize validates token and authorizes users
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			controllers.HTTPRes(c, http.StatusForbidden, "No Authorization header provided", nil)
			c.Abort()
			return
		}

		extractedToken := strings.Split(clientToken, "Bearer ")

		if len(extractedToken) == 2 {
			clientToken = strings.TrimSpace(extractedToken[1])
		} else {
			controllers.HTTPRes(c, http.StatusForbidden, "Incorrect Format of Authorization Token", nil)
			c.Abort()
			return
		}

		claims, err := user.AppJwtWrapper.ValidateToken(clientToken)
		if err != nil {
			controllers.HTTPRes(c, http.StatusUnauthorized, "Error while validating token", err.Error())
			c.Abort()
			return
		}

		c.Set("email", claims.Email)

		c.Next()

	}
}
