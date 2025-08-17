package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Recovery middleware for handling panics
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": err,
			})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}
