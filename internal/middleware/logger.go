package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Logger middleware para logging de solicitudes HTTP usando zerolog
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()
		latency := time.Since(start)
		if raw != "" {
			path = path + "?" + raw
		}
		logEvent := log.Info()
		if len(c.Errors) > 0 {
			logEvent = log.Error().Str("error", c.Errors.String())
		}

		logEvent.
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", c.Writer.Status()).
			Str("ip", c.ClientIP()).
			Dur("latency", latency).
			Str("user_agent", c.Request.UserAgent()).
			Msg("HTTP Request")
	}
}
