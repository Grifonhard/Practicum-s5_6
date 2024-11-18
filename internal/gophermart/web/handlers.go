package web

import (
	"net/http"

	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/auth"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/order/storage"
	"github.com/gin-gonic/gin"
)

type RegRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Registration(m *auth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegRequest
		
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		token, err := m.Registration(req.Login, req.Password)
		if err == storage.ErrUserExist {
			c.JSON(http.StatusConflict, gin.H{"error": storage.ErrUserExist.Error()})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return
		}

		// c.Header("Authorization", token)
		// TODO передаём в тело
		c.JSON(http.StatusOK, gin.H{"message": "User registered and authenticated"})
	}
}

func Login(m *auth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegRequest
		
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		token, err := m.Login(req.Login, req.Password)
		if err == auth.ErrWrongPassword || err == storage.ErrUserNotExist {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect login password pair"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return
		}

		// c.Header("Authorization", token)
		// TODO передаём в тело
		c.JSON(http.StatusOK, gin.H{"message": "User registered and authenticated"})
	}
}