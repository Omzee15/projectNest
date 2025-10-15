package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"lucid-lists-backend/internal/services"
	"lucid-lists-backend/internal/utils"
	"lucid-lists-backend/pkg/logger"
)

// Authentication validates JWT tokens and sets user context
func Authentication(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.WithComponent("auth-middleware").
				Warn("No authorization header provided")
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		// Check if token has Bearer prefix
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logger.WithComponent("auth-middleware").
				Warn("Invalid authorization header format")
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Validate token
		claims, err := authService.ValidateToken(token)
		if err != nil {
			logger.WithComponent("auth-middleware").
				WithFields(map[string]interface{}{
					"error": err.Error(),
				}).
				Warn("Invalid or expired token")
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("user_uid", claims.UserUID)
		c.Set("user_email", claims.Email)
		c.Set("user_name", claims.Name)

		logger.WithComponent("auth-middleware").
			WithFields(map[string]interface{}{
				"user_uid": claims.UserUID.String(),
				"email":    claims.Email,
			}).
			Info("Authentication successful")

		c.Next()
	}
}
