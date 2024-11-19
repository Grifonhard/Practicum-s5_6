package web

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/logger"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/model"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/repository"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/services/auth"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/services/order"
	"github.com/gin-gonic/gin"
)

type RegRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Registration(m *auth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {

		logger.Debug("handlers Registration %+v", c.Request)

		var req RegRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		token, err := m.Registration(req.Login, req.Password)
		if err != nil {
			if errors.Is(err, repository.ErrUserExist) {
				c.JSON(http.StatusConflict, gin.H{"error": repository.ErrUserExist.Error()})
				return
			}
			logger.Error("fail registed user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return
		}

		c.SetCookie(
			COOKIEAUTH,        // Имя cookie
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

		logger.Debug("handlers Login %+v", c.Request)

		var req RegRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		token, err := m.Login(req.Login, req.Password)
		if err != nil {
			if errors.Is(err, auth.ErrWrongPassword) || errors.Is(err, repository.ErrUserNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect login password pair"})
				return
			}
			logger.Error("fail auth: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return
		}

		c.SetCookie(
			COOKIEAUTH,        // Имя cookie
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

		logger.Debug("handlers AddOrder %+v", c.Request)

		if c.GetHeader("Content-Type") != "text/plain" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type must be text/plain"})
			return
		}

		userIDinterface, exists := c.Get(USERID)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username not found in context"})
			return
		}

		userID, ok := userIDinterface.(int)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username type assertion failed"})
			return
		}
		rawData, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "failed to read request body"})
			return
		}
	
		// Преобразуем тело запроса (число) в int
		orderID, err := strconv.Atoi(string(rawData))
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid request format"})
			return
		}

		err = m.AddOrder(userID, orderID)
		if err != nil {
			if errors.Is(err, order.ErrOrderExistThis) {
				c.JSON(http.StatusOK, "success")
				return
			}
			if errors.Is(err, order.ErrOrderExistThis) {
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

		logger.Debug("handlers ListOrders %+v", c.Request)

		userIDinterface, exists := c.Get(USERID)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username not found in context"})
			return
		}

		userID, ok := userIDinterface.(int)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username type assertion failed"})
			return
		}

		orders, err := m.ListOrders(userID)
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

		logger.Debug("handlers Balance %+v", c.Request)

		userIDinterface, exists := c.Get(USERID)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username not found in context"})
			return
		}

		userID, ok := userIDinterface.(int)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username type assertion failed"})
			return
		}

		balance, err := m.Balance(userID)
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

		logger.Debug("handlers Withdraw %+v", c.Request)

		userIDinterface, exists := c.Get(USERID)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username not found in context"})
			return
		}

		userID, ok := userIDinterface.(int)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username type assertion failed"})
			return
		}

		var req model.WithdrawRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "cann't parse json"})
			return
		}

		err := m.Withdraw(userID, req.Order, req.Sum)
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

		logger.Debug("handlers Withdrawals %+v", c.Request)

		userIDinterface, exists := c.Get(USERID)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username not found in context"})
			return
		}

		userID, ok := userIDinterface.(int)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username type assertion failed"})
			return
		}

		ws, err := m.Withdrawls(userID)
		if err != nil {
			if errors.Is(err, repository.ErrTransNotFound) {
				c.JSON(http.StatusNoContent, ws)
				return
			}
			logger.Error("fail withdrawals: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		c.JSON(http.StatusOK, ws)
	}
}
