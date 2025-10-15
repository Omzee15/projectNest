package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/internal/services"
	"lucid-lists-backend/internal/utils"
	"lucid-lists-backend/pkg/logger"
)

type CanvasHandler struct {
	canvasService *services.CanvasService
}

func NewCanvasHandler(canvasService *services.CanvasService) *CanvasHandler {
	return &CanvasHandler{
		canvasService: canvasService,
	}
}

// GetCanvas handles GET /api/projects/:uid/canvas
func (h *CanvasHandler) GetCanvas(c *gin.Context) {
	uidParam := c.Param("uid")

	logger.WithComponent("canvas-handler").
		WithFields(map[string]interface{}{"project_uid": uidParam}).
		Info("Getting canvas for project")

	projectUID, err := uuid.Parse(uidParam)
	if err != nil {
		logger.WithComponent("canvas-handler").
			WithFields(map[string]interface{}{"invalid_uid": uidParam}).
			Warn("Invalid project UID format")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project UID format")
		return
	}

	canvas, err := h.canvasService.GetCanvas(c.Request.Context(), projectUID)
	if err != nil {
		logger.WithComponent("canvas-handler").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Failed to get canvas")

		if err.Error() == "project not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Project not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve canvas")
		}
		return
	}

	logger.WithComponent("canvas-handler").
		WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
		Info("Successfully retrieved canvas")

	utils.SuccessResponse(c, canvas, "")
}

// UpdateCanvas handles POST/PUT /api/projects/:uid/canvas
func (h *CanvasHandler) UpdateCanvas(c *gin.Context) {
	uidParam := c.Param("uid")
	var req models.CanvasRequest

	logger.WithComponent("canvas-handler").
		WithFields(map[string]interface{}{"project_uid": uidParam}).
		Info("Updating canvas for project")

	projectUID, err := uuid.Parse(uidParam)
	if err != nil {
		logger.WithComponent("canvas-handler").
			WithFields(map[string]interface{}{"invalid_uid": uidParam}).
			Warn("Invalid project UID format")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project UID format")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithComponent("canvas-handler").
			WithFields(map[string]interface{}{"error": err.Error()}).
			Warn("Invalid request body for update canvas")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.canvasService.UpdateCanvas(c.Request.Context(), projectUID, &req)
	if err != nil {
		logger.WithComponent("canvas-handler").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Failed to update canvas")

		if err.Error() == "project not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Project not found")
		} else if err.Error() == "invalid JSON format" {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON format in canvas state")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update canvas")
		}
		return
	}

	logger.WithComponent("canvas-handler").
		WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
		Info("Successfully updated canvas")

	utils.SuccessResponse(c, gin.H{"message": "Canvas updated successfully"}, "Canvas updated successfully")
}

// DeleteCanvas handles DELETE /api/projects/:uid/canvas
func (h *CanvasHandler) DeleteCanvas(c *gin.Context) {
	uidParam := c.Param("uid")

	logger.WithComponent("canvas-handler").
		WithFields(map[string]interface{}{"project_uid": uidParam}).
		Info("Deleting canvas for project")

	projectUID, err := uuid.Parse(uidParam)
	if err != nil {
		logger.WithComponent("canvas-handler").
			WithFields(map[string]interface{}{"invalid_uid": uidParam}).
			Warn("Invalid project UID format")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project UID format")
		return
	}

	err = h.canvasService.DeleteCanvas(c.Request.Context(), projectUID)
	if err != nil {
		logger.WithComponent("canvas-handler").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Failed to delete canvas")

		if err.Error() == "project not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Project not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete canvas")
		}
		return
	}

	logger.WithComponent("canvas-handler").
		WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
		Info("Successfully deleted canvas")

	utils.SuccessResponse(c, gin.H{"message": "Canvas deleted successfully"}, "Canvas deleted successfully")
}
