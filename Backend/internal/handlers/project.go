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

// Helper function to extract user information from authentication context
func getUserFromContext(c *gin.Context) (userID int, userUID uuid.UUID, userName string, err error) {
	userUIDInterface, exists := c.Get("user_uid")
	if !exists {
		err = utils.NewBadRequestError("User not authenticated")
		return
	}

	// The middleware sets user_uid as uuid.UUID type
	userUID, ok := userUIDInterface.(uuid.UUID)
	if !ok {
		err = utils.NewBadRequestError("Invalid user UID type")
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		err = utils.NewBadRequestError("User ID not found")
		return
	}

	// The middleware sets user_id as int type
	userID, ok = userIDInterface.(int)
	if !ok {
		err = utils.NewBadRequestError("Invalid user ID type")
		return
	}

	userNameInterface, exists := c.Get("user_name")
	if exists {
		userName = userNameInterface.(string)
	}

	return
}

// GetProjects handles GET /api/projects
func (h *ProjectHandler) GetProjects(c *gin.Context) {
	logger.WithComponent("project-handler").Info("Getting all projects")

	// Extract user information from authentication context
	userID, _, _, err := getUserFromContext(c)
	if err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{"error": err.Error()}).
			Error("Failed to get user from context")
		utils.ErrorResponse(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	projects, err := h.projectService.GetAllProjects(c.Request.Context(), userID)
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

	// Debug: Log the DBML content being returned
	if project.DbmlContent != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"dbml_length": len(*project.DbmlContent),
				"dbml_preview": func() string {
					if len(*project.DbmlContent) > 100 {
						return (*project.DbmlContent)[:100] + "..."
					}
					return *project.DbmlContent
				}(),
			}).
			Info("Returning project with DBML content")
	} else {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
			Info("Returning project without DBML content")
	}

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

	// Extract user information from authentication context
	userID, userUID, _, err := getUserFromContext(c)
	if err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{"error": err.Error()}).
			Error("Failed to get user from context")
		utils.ErrorResponse(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	logger.WithComponent("project-handler").
		WithFields(map[string]interface{}{"project_name": req.Name, "user_id": userID}).
		Info("Creating new project")

	project, err := h.projectService.CreateProject(c.Request.Context(), &req, userID, userUID)
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

	// Debug: Log the received updates
	logger.WithComponent("project-handler").
		WithFields(map[string]interface{}{
			"project_uid": projectUID.String(),
			"updates":     updates,
		}).
		Info("Received partial update request")

	if updates.DbmlContent != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{
				"project_uid":      projectUID.String(),
				"dbml_content_len": len(*updates.DbmlContent),
			}).
			Info("DBML content received")
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

// GetProjectProgress handles GET /api/projects/:uid/progress
func (h *ProjectHandler) GetProjectProgress(c *gin.Context) {
	uidParam := c.Param("uid")

	logger.WithComponent("project-handler").
		WithFields(map[string]interface{}{"project_uid": uidParam}).
		Info("Getting project progress")

	projectUID, err := uuid.Parse(uidParam)
	if err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{
				"project_uid": uidParam,
				"error":       err.Error(),
			}).
			Error("Invalid project UID format")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project UID format")
		return
	}

	progress, err := h.projectService.GetProjectProgress(c.Request.Context(), projectUID)
	if err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Failed to get project progress")

		if err.Error() == "Project not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Project not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve project progress")
		}
		return
	}

	logger.WithComponent("project-handler").
		WithFields(map[string]interface{}{
			"project_uid":     projectUID.String(),
			"total_tasks":     progress.TotalTasks,
			"completed_tasks": progress.CompletedTasks,
			"progress":        progress.Progress,
		}).
		Info("Successfully retrieved project progress")

	utils.SuccessResponse(c, progress, "")
}

// GetProjectsWithProgress handles GET /api/projects?include_progress=true
func (h *ProjectHandler) GetProjectsWithProgress(c *gin.Context) {
	includeProgress := c.Query("include_progress")

	if includeProgress == "true" {
		logger.WithComponent("project-handler").Info("Getting all projects with progress")

		// Extract user information from authentication context
		userID, _, _, err := getUserFromContext(c)
		if err != nil {
			logger.WithComponent("project-handler").
				WithFields(map[string]interface{}{"error": err.Error()}).
				Error("Failed to get user from context")
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authentication required")
			return
		}

		projects, err := h.projectService.GetAllProjectsWithProgress(c.Request.Context(), userID)
		if err != nil {
			logger.WithComponent("project-handler").
				WithFields(map[string]interface{}{"error": err.Error()}).
				Error("Failed to get projects with progress")
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve projects with progress")
			return
		}

		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{"count": len(projects)}).
			Info("Successfully retrieved projects with progress")

		utils.SuccessResponse(c, projects, "")
	} else {
		// Fallback to regular GetProjects method
		h.GetProjects(c)
	}
}

// AddProjectMember handles POST /api/projects/:uid/members
func (h *ProjectHandler) AddProjectMember(c *gin.Context) {
	uidParam := c.Param("uid")

	logger.WithComponent("project-handler").
		WithFields(map[string]interface{}{"project_uid": uidParam}).
		Info("Adding member to project")

	projectUID, err := uuid.Parse(uidParam)
	if err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{"error": err.Error()}).
			Error("Invalid project UID")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project UID format")
		return
	}

	// Extract requesting user information from authentication context
	_, requestingUserUID, _, err := getUserFromContext(c)
	if err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{"error": err.Error()}).
			Error("Failed to get user from context")
		utils.ErrorResponse(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req models.AddProjectMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{"error": err.Error()}).
			Error("Invalid request body")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Set default role to 'member' if not specified
	if req.Role == "" {
		req.Role = "member"
	}

	// Validate role
	if req.Role != "owner" && req.Role != "member" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid role. Must be 'owner' or 'member'")
		return
	}

	// Set default canWrite to false if not specified
	canWrite := false
	if req.CanWrite != nil {
		canWrite = *req.CanWrite
	}

	// Owners always have write access
	if req.Role == "owner" {
		canWrite = true
	}

	err = h.projectService.AddMemberByEmail(c.Request.Context(), projectUID, req.Email, req.Role, canWrite, requestingUserUID)
	if err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"email":       req.Email,
				"error":       err.Error(),
			}).
			Error("Failed to add member to project")

		if err.Error() == "user not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "User with this email not found")
			return
		}
		if err.Error() == "user already a member" {
			utils.ErrorResponse(c, http.StatusConflict, "User is already a member of this project")
			return
		}
		if err.Error() == "Only project owners can grant write access" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to add member to project")
		return
	}

	logger.WithComponent("project-handler").
		WithFields(map[string]interface{}{
			"project_uid": projectUID.String(),
			"email":       req.Email,
			"role":        req.Role,
			"can_write":   canWrite,
		}).
		Info("Successfully added member to project")

	utils.SuccessResponse(c, gin.H{
		"message":   "Member added successfully",
		"email":     req.Email,
		"role":      req.Role,
		"can_write": canWrite,
	}, "Member added successfully")
}

// GetProjectMembers handles GET /api/projects/:uid/members
func (h *ProjectHandler) GetProjectMembers(c *gin.Context) {
	uidParam := c.Param("uid")

	logger.WithComponent("project-handler").
		WithFields(map[string]interface{}{"project_uid": uidParam}).
		Info("Getting project members")

	projectUID, err := uuid.Parse(uidParam)
	if err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{"error": err.Error()}).
			Error("Invalid project UID")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project UID format")
		return
	}

	members, err := h.projectService.GetProjectMembers(c.Request.Context(), projectUID)
	if err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Failed to get project members")
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve project members")
		return
	}

	logger.WithComponent("project-handler").
		WithFields(map[string]interface{}{
			"project_uid": projectUID.String(),
			"count":       len(members),
		}).
		Info("Successfully retrieved project members")

	utils.SuccessResponse(c, members, "")
}

// UpdateMemberWriteAccess handles PUT /api/projects/:uid/members/:memberUid/write-access
func (h *ProjectHandler) UpdateMemberWriteAccess(c *gin.Context) {
	projectUIDParam := c.Param("uid")
	memberUIDParam := c.Param("memberUid")

	logger.WithComponent("project-handler").
		WithFields(map[string]interface{}{
			"project_uid": projectUIDParam,
			"member_uid":  memberUIDParam,
		}).
		Info("Updating member write access")

	projectUID, err := uuid.Parse(projectUIDParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project UID")
		return
	}

	memberUID, err := uuid.Parse(memberUIDParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid member UID")
		return
	}

	// Extract user information from authentication context
	_, requestingUserUID, _, err := getUserFromContext(c)
	if err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{"error": err.Error()}).
			Error("Failed to get user from context")
		utils.ErrorResponse(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	type UpdateMemberAccessRequest struct {
		CanWrite bool `json:"can_write"`
	}

	var req UpdateMemberAccessRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{"error": err.Error()}).
			Error("Invalid request body")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	err = h.projectService.UpdateMemberWriteAccess(c.Request.Context(), projectUID, memberUID, req.CanWrite, requestingUserUID)
	if err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"member_uid":  memberUID.String(),
				"error":       err.Error(),
			}).
			Error("Failed to update member write access")

		if err.Error() == "Project not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Project not found")
			return
		}
		if err.Error() == "Only project owners can change member access" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		if err.Error() == "Cannot change write access for project owners" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		if err.Error() == "Member not found in project" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update member write access")
		return
	}

	logger.WithComponent("project-handler").
		WithFields(map[string]interface{}{
			"project_uid": projectUID.String(),
			"member_uid":  memberUID.String(),
			"can_write":   req.CanWrite,
		}).
		Info("Successfully updated member write access")

	utils.SuccessResponse(c, gin.H{
		"message":   "Member write access updated successfully",
		"can_write": req.CanWrite,
	}, "")
}

// SearchUsers handles GET /api/users/search?q=...&project_uid=...
func (h *ProjectHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")

	var projectUID *uuid.UUID
	if projectUIDStr := c.Query("project_uid"); projectUIDStr != "" {
		parsed, err := uuid.Parse(projectUIDStr)
		if err == nil {
			projectUID = &parsed
		}
	}

	users, err := h.projectService.SearchUsers(c.Request.Context(), query, projectUID)
	if err != nil {
		logger.WithComponent("project-handler").
			WithFields(map[string]interface{}{
				"query": query,
				"error": err.Error(),
			}).
			Error("Failed to search users")
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to search users")
		return
	}

	utils.SuccessResponse(c, users, "")
}
