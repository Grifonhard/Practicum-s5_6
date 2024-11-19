package web

import (
	"net/http"

	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/repository"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/services/auth"
	"github.com/gin-gonic/gin"
)

const (
	USERID = "username"
	COOKIEAUTH = "auth_token"
)

func Authentication(am *auth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {

		authCookie, err := c.Cookie(COOKIEAUTH)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication cookie required"})
			c.Abort()
			return
		}

		userID, err := am.Authentication(authCookie)
		if err == auth.ErrInvalidToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if err == repository.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
			c.Abort()
			return
		}

		c.Set(USERID, userID)

		c.Next()
	}
}
