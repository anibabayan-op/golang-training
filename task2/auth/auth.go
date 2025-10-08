package middleware

import (
	"golang-training/task2/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		user, err := userService.GetCurrentUser(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("currentUser", user)
		c.Set("token", token)
		c.Next()
	}
}
