package handlers

import (
	"net/http"

	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/internal/services"
	"lucid-lists-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserSettingsHandler struct {
	service *services.UserSettingsService
}

func NewUserSettingsHandler(service *services.UserSettingsService) *UserSettingsHandler {
	return &UserSettingsHandler{
		service: service,
	}
}

// GetUserSettings retrieves the current user's settings
func (h *UserSettingsHandler) GetUserSettings(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	settings, err := h.service.GetUserSettings(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user settings: "+err.Error())
		return
	}

	utils.SuccessResponse(c, settings, "Settings retrieved successfully")
}

// UpdateUserSettings updates the current user's settings
func (h *UserSettingsHandler) UpdateUserSettings(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var request models.UserSettingsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	settings, err := h.service.UpdateUserSettings(c.Request.Context(), userID, &request)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to update user settings: "+err.Error())
		return
	}

	utils.SuccessResponse(c, settings, "Settings updated successfully")
}

// ResetUserSettings resets the current user's settings to defaults
func (h *UserSettingsHandler) ResetUserSettings(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Initialize default settings (this will overwrite existing settings)
	err := h.service.InitializeDefaultSettings(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to reset settings: "+err.Error())
		return
	}

	// Get the reset settings to return
	settings, err := h.service.GetUserSettings(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve reset settings: "+err.Error())
		return
	}

	utils.SuccessResponse(c, settings, "Settings reset to defaults successfully")
}
