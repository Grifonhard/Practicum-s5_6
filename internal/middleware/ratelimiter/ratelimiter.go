package ratelimiter

import (
	"fmt"
	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	requestsLimit = 60
	rate          = time.Minute
)

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.Header("Retry-After", time.Until(info.ResetTime).String())
	format := fmt.Sprintf("Too many requests. Try again in %s", time.Until(info.ResetTime).String())
	c.String(http.StatusTooManyRequests, format)
}

func NewRateLimiter() gin.HandlerFunc {
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  rate,
		Limit: requestsLimit,
	})
	rateLimiter := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})

	return rateLimiter
}
