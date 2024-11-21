package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func badRequestResponse(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err})
}
