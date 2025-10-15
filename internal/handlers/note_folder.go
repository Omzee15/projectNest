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

type NoteFolderHandler struct {
	folderService *services.NoteFolderService
}

func NewNoteFolderHandler(folderService *services.NoteFolderService) *NoteFolderHandler {
	return &NoteFolderHandler{
		folderService: folderService,
	}
}

// GetFolders retrieves all folders for a project
func (h *NoteFolderHandler) GetFolders(c *gin.Context) {
	projectUIDStr := c.Param("uid")
	projectUID, err := uuid.Parse(projectUIDStr)
	if err != nil {
		logger.WithComponent("note-folder-handler").
			WithFields(map[string]interface{}{
				"project_uid": projectUIDStr,
				"error":       err.Error(),
			}).
			Error("Invalid project UID")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project UID")
		return
	}

	logger.WithComponent("note-folder-handler").
		WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
		Info("Getting folders for project")

	folders, err := h.folderService.GetFoldersByProject(c.Request.Context(), projectUID)
	if err != nil {
		logger.WithComponent("note-folder-handler").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Failed to get folders")
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get folders")
		return
	}

	utils.SuccessResponse(c, folders, "Folders retrieved successfully")
}

// CreateFolder creates a new folder
func (h *NoteFolderHandler) CreateFolder(c *gin.Context) {
	projectUIDStr := c.Param("uid")
	projectUID, err := uuid.Parse(projectUIDStr)
	if err != nil {
		logger.WithComponent("note-folder-handler").
			WithFields(map[string]interface{}{
				"project_uid": projectUIDStr,
				"error":       err.Error(),
			}).
			Error("Invalid project UID for folder creation")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project UID")
		return
	}

	var request models.NoteFolderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.WithComponent("note-folder-handler").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Invalid request body")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := utils.ValidateStruct(request); err != nil {
		logger.WithComponent("note-folder-handler").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Validation failed")
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed")
		return
	}

	logger.WithComponent("note-folder-handler").
		WithFields(map[string]interface{}{
			"project_uid": projectUID.String(),
			"name":        request.Name,
		}).
		Info("Creating folder for project")

	folder, err := h.folderService.CreateFolder(c.Request.Context(), projectUID, request)
	if err != nil {
		logger.WithComponent("note-folder-handler").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"name":        request.Name,
				"error":       err.Error(),
			}).
			Error("Failed to create folder")
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create folder")
		return
	}

	utils.CreatedResponse(c, folder, "Folder created successfully")
}

// UpdateFolder updates an existing folder
func (h *NoteFolderHandler) UpdateFolder(c *gin.Context) {
	folderUIDStr := c.Param("uid")
	folderUID, err := uuid.Parse(folderUIDStr)
	if err != nil {
		logger.WithComponent("note-folder-handler").
			WithFields(map[string]interface{}{
				"folder_uid": folderUIDStr,
				"error":      err.Error(),
			}).
			Error("Invalid folder UID")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid folder UID")
		return
	}

	var request models.NoteFolderUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.WithComponent("note-folder-handler").
			WithFields(map[string]interface{}{
				"folder_uid": folderUID.String(),
				"error":      err.Error(),
			}).
			Error("Invalid request body")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := utils.ValidateStruct(request); err != nil {
		logger.WithComponent("note-folder-handler").
			WithFields(map[string]interface{}{
				"folder_uid": folderUID.String(),
				"error":      err.Error(),
			}).
			Error("Validation failed")
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed")
		return
	}

	logger.WithComponent("note-folder-handler").
		WithFields(map[string]interface{}{"folder_uid": folderUID.String()}).
		Info("Updating folder")

	folder, err := h.folderService.UpdateFolder(c.Request.Context(), folderUID, request)
	if err != nil {
		logger.WithComponent("note-folder-handler").
			WithFields(map[string]interface{}{
				"folder_uid": folderUID.String(),
				"error":      err.Error(),
			}).
			Error("Failed to update folder")
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update folder")
		return
	}

	utils.SuccessResponse(c, folder, "Folder updated successfully")
}

// DeleteFolder deletes a folder
func (h *NoteFolderHandler) DeleteFolder(c *gin.Context) {
	folderUIDStr := c.Param("uid")
	folderUID, err := uuid.Parse(folderUIDStr)
	if err != nil {
		logger.WithComponent("note-folder-handler").
			WithFields(map[string]interface{}{
				"folder_uid": folderUIDStr,
				"error":      err.Error(),
			}).
			Error("Invalid folder UID")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid folder UID")
		return
	}

	logger.WithComponent("note-folder-handler").
		WithFields(map[string]interface{}{"folder_uid": folderUID.String()}).
		Info("Deleting folder")

	err = h.folderService.DeleteFolder(c.Request.Context(), folderUID)
	if err != nil {
		logger.WithComponent("note-folder-handler").
			WithFields(map[string]interface{}{
				"folder_uid": folderUID.String(),
				"error":      err.Error(),
			}).
			Error("Failed to delete folder")
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete folder")
		return
	}

	utils.SuccessResponse(c, nil, "Folder deleted successfully")
}
