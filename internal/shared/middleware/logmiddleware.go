package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		event := log.Info()

		switch {
		case status >= 500:
			event = log.Error()
		case status >= 400:
			event = log.Warn()
		}

		event.
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", status).
			Int64("duration_ms", duration.Milliseconds()).
			Str("client_ip", c.ClientIP()).
			Msg("HTTP request")
	}
}
