package handler

import (
	"fmt"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/service"
	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"time"
)

const (
	retryDelay = "60"
)

type Handler struct {
	Router         *gin.Engine
	OrderService   *service.OrderService
	AccrualService *service.AccrualService
}

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.Header("Retry-After", retryDelay)
	format := fmt.Sprintf("Too many requests. Try again in %s", time.Until(info.ResetTime).String())
	c.String(429, format)
}

func NewHandler(
	router *gin.Engine,
	orderService *service.OrderService,
	accrualService *service.AccrualService,
) {
	h := &Handler{
		Router:         router,
		OrderService:   orderService,
		AccrualService: accrualService,
	}
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Second,
		Limit: 5,
	})
	rateLimiter := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})

	g := h.Router.Group("/api")

	g.GET("/orders/:number", rateLimiter, h.GetOrderHandler)
	g.POST("/orders", h.OrderRegistrationHandler)
	g.POST("/goods", h.AccrualRegistrationHandler)
	g.GET("/orders", h.GetOrdersHandler)
}
