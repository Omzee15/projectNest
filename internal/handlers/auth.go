package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/internal/services"
	"lucid-lists-backend/internal/utils"
	"lucid-lists-backend/pkg/logger"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var request models.RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.WithComponent("auth-handler").
			WithFields(map[string]interface{}{
				"error": err.Error(),
			}).
			Error("Invalid request body for registration")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	logger.WithComponent("auth-handler").
		WithFields(map[string]interface{}{
			"email": request.Email,
			"name":  request.Name,
		}).
		Info("Processing user registration")

	response, err := h.authService.Register(c.Request.Context(), request)
	if err != nil {
		logger.WithComponent("auth-handler").
			WithFields(map[string]interface{}{
				"email": request.Email,
				"error": err.Error(),
			}).
			Error("Failed to register user")

		if err.Error() == "user with this email already exists" {
			utils.ErrorResponse(c, http.StatusConflict, "User with this email already exists")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to register user")
		}
		return
	}

	logger.WithComponent("auth-handler").
		WithFields(map[string]interface{}{
			"user_uid": response.User.UserUID.String(),
			"email":    response.User.Email,
		}).
		Info("Successfully registered user")

	utils.SuccessResponse(c, response, "User registered successfully")
}

// Login handles POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var request models.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.WithComponent("auth-handler").
			WithFields(map[string]interface{}{
				"error": err.Error(),
			}).
			Error("Invalid request body for login")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	logger.WithComponent("auth-handler").
		WithFields(map[string]interface{}{
			"email": request.Email,
		}).
		Info("Processing user login")

	response, err := h.authService.Login(c.Request.Context(), request)
	if err != nil {
		logger.WithComponent("auth-handler").
			WithFields(map[string]interface{}{
				"email": request.Email,
				"error": err.Error(),
			}).
			Error("Failed to login user")

		if err.Error() == "invalid email or password" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to login user")
		}
		return
	}

	logger.WithComponent("auth-handler").
		WithFields(map[string]interface{}{
			"user_uid": response.User.UserUID.String(),
			"email":    response.User.Email,
		}).
		Info("Successfully logged in user")

	utils.SuccessResponse(c, response, "User logged in successfully")
}
