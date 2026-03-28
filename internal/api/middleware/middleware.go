// Package middleware holds Gin middleware (CORS, etc.).
package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// CORSMiddleware allows browser clients from the specified origin and handles OPTIONS requests.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Cache-Control")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		reqLogger := log.With().
			Str("request_id", requestID).
			Logger()

		c.Set("logger", reqLogger)
		c.Header("X-Request-ID", requestID)

		ip := c.ClientIP()
		ua := c.Request.UserAgent()

		c.Next()

		status := c.Writer.Status()

		var evt *zerolog.Event
		switch {
		case status >= 500:
			evt = reqLogger.Error()
		case status >= 400:
			evt = reqLogger.Warn()
		default:
			evt = reqLogger.Info()
		}

		evt.
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", status).
			Int64("latency_ms", time.Since(start).Milliseconds()).
			Str("ip", ip).
			Str("user_agent", ua).
			Msg("http request")
	}
}
