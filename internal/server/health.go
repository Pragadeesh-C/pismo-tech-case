package server

import (
	"github.com/gin-gonic/gin"
)

// HealthHandler returns a simple liveness probe for orchestrators and load balancers.
func HealthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	}
}
