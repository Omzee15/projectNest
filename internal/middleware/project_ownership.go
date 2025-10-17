package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"lucid-lists-backend/internal/repositories"
	"lucid-lists-backend/internal/utils"
	"lucid-lists-backend/pkg/logger"
)

// ProjectOwnership validates that the authenticated user has access to the requested project
// Updated to check project membership instead of ownership
func ProjectOwnership(projectRepo repositories.ProjectRepository, db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract project UID from URL parameter
		projectUIDParam := c.Param("uid")
		if projectUIDParam == "" {
			logger.WithComponent("project-ownership-middleware").
				Warn("No project UID provided in request")
			utils.ErrorResponse(c, http.StatusBadRequest, "Project UID is required")
			c.Abort()
			return
		}

		// Validate project UID format
		projectUID, err := uuid.Parse(projectUIDParam)
		if err != nil {
			logger.WithComponent("project-ownership-middleware").
				WithFields(map[string]interface{}{
					"project_uid": projectUIDParam,
					"error":       err.Error(),
				}).
				Warn("Invalid project UID format")
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project UID format")
			c.Abort()
			return
		}

		// Check if project exists
		project, err := projectRepo.GetByUID(c.Request.Context(), projectUID)
		if err != nil {
			logger.WithComponent("project-ownership-middleware").
				WithFields(map[string]interface{}{
					"project_uid": projectUID.String(),
					"error":       err.Error(),
				}).
				Error("Failed to retrieve project for ownership validation")
			utils.ErrorResponse(c, http.StatusNotFound, "Project not found")
			c.Abort()
			return
		}

		// Get authenticated user from context
		userUIDInterface, exists := c.Get("user_uid")
		if !exists {
			logger.WithComponent("project-ownership-middleware").
				Error("No authenticated user found in context")
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authentication required")
			c.Abort()
			return
		}

		userUID, ok := userUIDInterface.(uuid.UUID)
		if !ok {
			logger.WithComponent("project-ownership-middleware").
				Error("Invalid user UID type in context")
			utils.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context")
			c.Abort()
			return
		}

		// Check if user is a member of the project
		isMember, err := checkProjectMembership(c.Request.Context(), db, project.ID, userUID)
		if err != nil {
			logger.WithComponent("project-ownership-middleware").
				WithFields(map[string]interface{}{
					"project_uid": projectUID.String(),
					"user_uid":    userUID.String(),
					"error":       err.Error(),
				}).
				Error("Failed to check project membership")
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to validate access")
			c.Abort()
			return
		}

		if !isMember {
			logger.WithComponent("project-ownership-middleware").
				WithFields(map[string]interface{}{
					"project_uid": projectUID.String(),
					"user_uid":    userUID.String(),
				}).
				Warn("User is not a member of this project")
			utils.ErrorResponse(c, http.StatusForbidden, "Access denied - you are not a member of this project")
			c.Abort()
			return
		}

		logger.WithComponent("project-ownership-middleware").
			WithFields(map[string]interface{}{
				"project_uid":  projectUID.String(),
				"project_name": project.Name,
				"user_uid":     userUID.String(),
			}).
			Info("Project membership validation passed")

		// Store project information in context for use by handlers
		c.Set("project", project)
		c.Set("project_uid", projectUID)

		c.Next()
	}
}

// checkProjectMembership checks if a user is a member of a project
func checkProjectMembership(ctx context.Context, db *pgxpool.Pool, projectID int, userUID uuid.UUID) (bool, error) {
	// First get the user's integer ID from their UUID
	userQuery := `SELECT id FROM users WHERE user_uid = $1 AND is_active = true`
	var userID int
	err := db.QueryRow(ctx, userQuery, userUID).Scan(&userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil // User doesn't exist, so not a member
		}
		return false, fmt.Errorf("failed to get user ID: %w", err)
	}

	// Now check membership using the integer user ID
	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM project_member 
			WHERE project_id = $1 AND user_id = $2
		)`

	var exists bool
	err = db.QueryRow(ctx, query, projectID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check membership: %w", err)
	}

	return exists, nil
}

// ResourceOwnership validates ownership for individual resources (notes, folders) by checking their parent project
// This is a more generic approach that can be used for any resource that belongs to a project
func ResourceOwnership(projectRepo repositories.ProjectRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract resource UID from URL parameter
		resourceUIDParam := c.Param("uid")
		if resourceUIDParam == "" {
			logger.WithComponent("resource-ownership-middleware").
				Warn("No resource UID provided in request")
			utils.ErrorResponse(c, http.StatusBadRequest, "Resource UID is required")
			c.Abort()
			return
		}

		// Validate resource UID format
		resourceUID, err := uuid.Parse(resourceUIDParam)
		if err != nil {
			logger.WithComponent("resource-ownership-middleware").
				WithFields(map[string]interface{}{
					"resource_uid": resourceUIDParam,
					"error":        err.Error(),
				}).
				Warn("Invalid resource UID format")
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid resource UID format")
			c.Abort()
			return
		}

		// TODO: When user authentication is implemented:
		// 1. Extract the resource from database
		// 2. Find its parent project
		// 3. Validate project ownership against authenticated user
		//
		// For now, we just validate the UID format and allow access
		logger.WithComponent("resource-ownership-middleware").
			WithFields(map[string]interface{}{
				"resource_uid": resourceUID.String(),
			}).
			Info("Resource ownership validation passed (simplified implementation)")

		// Store resource UID in context
		c.Set("resource_uid", resourceUID)

		c.Next()
	}
}
