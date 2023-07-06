package response

import (
	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, statusCode int, message string, status bool, data any) {
	c.JSON(statusCode, gin.H{
		"message":    message,
		"statusCode": statusCode,
		"status":     status,
		"data":       data,
	})
}
