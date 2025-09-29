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

type ProjectHandler struct {
	projectService *services.ProjectService
}

func NewProjectHandler(projectService *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// GetProjects handles GET /api/projects
func (h *ProjectHandler) GetProjects(c *gin.Context) {
	logger.WithComponent("project-handler").Info("Getting all projects")

	projects, err := h.projectService.GetAllProjects(c.Request.Context())
	if err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{"error": err.Error()}).
			Error("Failed to get projects")
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve projects")
		return
	}

	logger.WithComponent("project-handler").
		WithFields(map[string]interface{}{"count": len(projects)}).
		Info("Successfully retrieved projects")

	utils.SuccessResponse(c, projects, "")
}

// GetProject handles GET /api/projects/:uid
func (h *ProjectHandler) GetProject(c *gin.Context) {
	uidParam := c.Param("uid")

	logger.WithComponent("project-handler").
		WithFields(map[string]interface{}{"project_uid": uidParam}).
		Info("Getting project with lists and tasks")

	projectUID, err := uuid.Parse(uidParam)
	if err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{"invalid_uid": uidParam}).
			Warn("Invalid project UID format")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project UID format")
		return
	}

	project, err := h.projectService.GetProjectWithLists(c.Request.Context(), projectUID)
	if err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Failed to get project")
		utils.ErrorResponse(c, http.StatusNotFound, "Project not found")
		return
	}

	logger.WithComponent("project-handler").
		WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
		Info("Successfully retrieved project")

	utils.SuccessResponse(c, project, "")
}

// CreateProject handles POST /api/projects
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req models.ProjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{"error": err.Error()}).
			Warn("Invalid request body for create project")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	logger.WithComponent("project-handler").
		WithFields(map[string]interface{}{"project_name": req.Name}).
		Info("Creating new project")

	project, err := h.projectService.CreateProject(c.Request.Context(), &req)
	if err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{
				"project_name": req.Name,
				"error":        err.Error(),
			}).
			Error("Failed to create project")
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create project")
		return
	}

	logger.WithComponent("project-handler").
		WithFields(map[string]interface{}{
			"project_uid":  project.ProjectUID.String(),
			"project_name": project.Name,
		}).
		Info("Successfully created project")

	utils.CreatedResponse(c, project, "Project created successfully")
}

// UpdateProject handles PUT /api/projects/:uid
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	uidParam := c.Param("uid")

	logger.WithComponent("project-handler").
		WithFields(map[string]interface{}{"project_uid": uidParam}).
		Info("Updating project")

	projectUID, err := uuid.Parse(uidParam)
	if err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{"invalid_uid": uidParam}).
			Warn("Invalid project UID format")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project UID format")
		return
	}

	var req models.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{"error": err.Error()}).
			Warn("Invalid request body for update project")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	project, err := h.projectService.UpdateProject(c.Request.Context(), projectUID, &req)
	if err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Failed to update project")
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update project")
		return
	}

	logger.WithComponent("project-handler").
		WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
		Info("Successfully updated project")

	utils.SuccessResponse(c, project, "Project updated successfully")
}

// DeleteProject handles DELETE /api/projects/:uid
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	uidParam := c.Param("uid")

	logger.WithComponent("project-handler").
		WithFields(map[string]interface{}{"project_uid": uidParam}).
		Info("Deleting project")

	projectUID, err := uuid.Parse(uidParam)
	if err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{"invalid_uid": uidParam}).
			Warn("Invalid project UID format")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project UID format")
		return
	}

	err = h.projectService.DeleteProject(c.Request.Context(), projectUID)
	if err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Failed to delete project")
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete project")
		return
	}

	logger.WithComponent("project-handler").
		WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
		Info("Successfully deleted project")

	utils.SuccessResponse(c, nil, "Project deleted successfully")
}

// PartialUpdateProject handles PATCH /api/projects/:uid
func (h *ProjectHandler) PartialUpdateProject(c *gin.Context) {
	uidParam := c.Param("uid")

	projectUID, err := uuid.Parse(uidParam)
	if err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{"invalid_uid": uidParam}).
			Warn("Invalid project UID format")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project UID format")
		return
	}

	var updates models.ProjectUpdateRequest
	if err := c.ShouldBindJSON(&updates); err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{"error": err.Error()}).
			Warn("Invalid request body for partial update project")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	logger.WithComponent("project-handler").
		WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
		Info("Partially updating project")

	project, err := h.projectService.PartialUpdateProject(c.Request.Context(), projectUID, &updates)
	if err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Failed to partial update project")

		if err.Error() == "Project not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Project not found")
		} else if err.Error() == "No fields to update" {
			utils.ErrorResponse(c, http.StatusBadRequest, "No fields to update")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update project")
		}
		return
	}

	logger.WithComponent("project-handler").
		WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
		Info("Successfully partial updated project")

	utils.SuccessResponse(c, project, "Project updated successfully")
}
