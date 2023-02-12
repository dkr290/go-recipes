package handlers

import (
	"os"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-API-KEY") != os.Getenv("X-API-KEY") {
			c.AbortWithStatus(401)

		}
		c.Next()

	}
}
