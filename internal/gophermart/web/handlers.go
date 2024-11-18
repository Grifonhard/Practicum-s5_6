package web

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/logger"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/repository"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/services/auth"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/services/order"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/services/storage"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/model"
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
		if err != nil {
			if errors.Is(err, storage.ErrUserExist) {
				c.JSON(http.StatusConflict, gin.H{"error": storage.ErrUserExist.Error()})
				return
			}
			logger.Error("fail registed user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return
		}

		c.SetCookie(
			"auth_token",      // Имя cookie
			token,             // Значение cookie
			auth.EXPIREDAT*60, // Время жизни в секундах
			"/",               // Путь, где cookie будет доступен
			"",                // Домен 
			false,             // Secure (использовать только HTTPS, если true)
			true,              // HttpOnly (доступно только для HTTP запросов, не для JavaScript)
		)

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
		if err != nil {
			if errors.Is(err, auth.ErrWrongPassword) || errors.Is(err, storage.ErrUserNotExist) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect login password pair"})
				return
			}
			logger.Error("fail auth: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return
		}

		c.SetCookie(
			"auth_token",      // Имя cookie
			token,             // Значение cookie
			auth.EXPIREDAT*60, // Время жизни в секундах
			"/",               // Путь, где cookie будет доступен
			"",                // Домен 
			false,             // Secure (использовать только HTTPS, если true)
			true,              // HttpOnly (доступно только для HTTP запросов, не для JavaScript)
		)

		c.JSON(http.StatusOK, gin.H{"message": "User login and authenticated"})
	}
}

func AddOrder(m *order.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Content-Type") != "text/plain" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type must be text/plain"})
			return
		}

		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username not found in context"})
			return
		}

		usernameStr, ok := username.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username type assertion failed"})
			return
		}

		var orderID string
		if err := c.ShouldBind(&orderID); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid request format"})
			return
		}

		orderIDInt, err := strconv.Atoi(orderID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "orderID must be a number"})
			return
		}

		err = m.AddOrder(usernameStr, orderIDInt)
		if err != nil {
			if errors.Is(err, storage.ErrOrderExistThis) {
				c.JSON(http.StatusOK, "success")
				return
			}
			if errors.Is(err, storage.ErrOrderExistThis) {
				c.JSON(http.StatusConflict, gin.H{"error": "no orders found"})
				return
			}
			logger.Error("fail add order: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add order"})
			return
		}

		c.JSON(http.StatusAccepted, "success")
	}
}

func ListOrders(m *order.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username not found in context"})
			return
		}

		usernameStr, ok := username.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username type assertion failed"})
			return
		}

		orders, err := m.ListOrders(usernameStr)
		if err != nil {
			if errors.Is(err, repository.ErrOrdersNotFound) {
				c.JSON(http.StatusNoContent, orders)
				return
			}
			logger.Error("fail list orders: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		c.JSON(http.StatusOK, orders)		
	}
}

func Balance(m *order.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username not found in context"})
			return
		}

		usernameStr, ok := username.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username type assertion failed"})
			return
		}

		balance, err := m.Balance(usernameStr)
		if err != nil {
			logger.Error("fail get balance: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		c.JSON(http.StatusOK, balance)
	}
}

func Withdraw(m *order.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username not found in context"})
			return
		}

		usernameStr, ok := username.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username type assertion failed"})
			return
		}

		var req model.WithdrawRequest
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid request format"})
			return
		}

		err := m.Withdraw(usernameStr, req.Order, req.Sum)
		if err != nil {
			if errors.Is(err, order.ErrNotEnoughBalance) {
				c.JSON(http.StatusPaymentRequired, gin.H{"error": "not anough points"})
				return
			}
			if errors.Is(err, order.ErrAlreadyDebited) || errors.Is(err, order.ErrTooMuchTransact) || errors.Is(err, repository.ErrOrderNotFound) {
				c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "wrong order number"})
				return
			}
			logger.Error("fail withdraw: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		c.JSON(http.StatusOK, "success")
	}
}

func Withdrawals(m *order.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username not found in context"})
			return
		}

		usernameStr, ok := username.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username type assertion failed"})
			return
		}

		ws, err := m.Withdrawls(usernameStr)
		if err != nil {
			if errors.Is(err, repository.ErrTransNotFound) {
				c.JSON(http.StatusNoContent, ws)
			}
			logger.Error("fail withdrawals: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		c.JSON(http.StatusOK, ws)
	}
}
