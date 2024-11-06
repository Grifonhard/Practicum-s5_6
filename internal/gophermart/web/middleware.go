package web

import (
	"net/http"

	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/auth"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/storage"
	"github.com/gin-gonic/gin"
)

func Authentication(am *auth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		err := am.Authentication(authHeader)
		if err == auth.ErrInvalidToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		} else if err == storage.ErrUserNotExist {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			c.Abort()
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
			c.Abort()
			return
		}

		c.Next()
	}
}