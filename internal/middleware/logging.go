package middleware

import (
	"bytes"
	"io"
	"time"

	"lucid-lists-backend/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// RequestLogging provides brief, structured request logging with debugging support
func RequestLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Skip OPTIONS requests - too verbose
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Read request body for debugging (non-destructively)
		var bodyBytes []byte
		if c.Request.Body != nil && (c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH") {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Brief request logging
		logger.WithComponent("http").
			WithFields(map[string]interface{}{
				"method": c.Request.Method,
				"path":   c.Request.URL.Path,
			}).
			Info("→ Request")

		// Process request
		c.Next()

		// Brief response logging
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		logFields := map[string]interface{}{
			"status":   statusCode,
			"duration": duration.Truncate(time.Microsecond).String(),
		}

		// Add request body to logs if request failed for debugging
		if statusCode >= 400 && len(bodyBytes) > 0 && len(bodyBytes) < 1000 {
			logFields["request_body"] = string(bodyBytes)
		}

		logEntry := logger.WithComponent("http").WithFields(logFields)

		if statusCode >= 500 {
			logEntry.Error("← Server Error")
		} else if statusCode >= 400 {
			logEntry.Warn("← Client Error")
		} else {
			logEntry.Info("← Success")
		}
	}
}

// CORSWithLogging provides CORS middleware with logging
func CORSWithLogging(allowedOrigins []string) gin.HandlerFunc {
	config := cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * 3600, // 12 hours
	}

	logger.WithComponent("cors").
		WithFields(map[string]interface{}{
			"allowed_origins": config.AllowOrigins,
			"allowed_methods": config.AllowMethods,
		}).
		Info("CORS middleware initialized")

	return cors.New(config)
}
