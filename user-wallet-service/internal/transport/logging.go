package transport

import (
	"net/http"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		msg := "http request"
		if len(c.Errors) > 0 {
			logger.Error(msg,
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"status", status,
				"latency", latency.String(),
				"client", c.ClientIP(),
				"error", c.Errors.String(),
			)
			return
		}
		level := "info"
		if status >= http.StatusInternalServerError {
			level = "error"
		} else if status >= http.StatusBadRequest {
			level = "warn"
		}
		fields := []interface{}{
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", status,
			"latency", latency.String(),
			"client", c.ClientIP(),
		}

		switch level {
		case "error":
			logger.Error(msg, fields...)
		case "warn":
			logger.Warn(msg, fields...)
		case "info":
			logger.Info(msg, fields...)
		}
	}
}
