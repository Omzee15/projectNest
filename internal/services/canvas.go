package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/internal/repositories"
	"lucid-lists-backend/pkg/logger"
)

type CanvasService struct {
	canvasRepo  repositories.CanvasRepository
	projectRepo repositories.ProjectRepository
}

func NewCanvasService(canvasRepo repositories.CanvasRepository, projectRepo repositories.ProjectRepository) *CanvasService {
	return &CanvasService{
		canvasRepo:  canvasRepo,
		projectRepo: projectRepo,
	}
}

// GetCanvas retrieves the canvas state for a project, or creates a default one if it doesn't exist
func (s *CanvasService) GetCanvas(ctx context.Context, projectUID uuid.UUID) (*models.CanvasResponse, error) {
	logger.WithComponent("canvas-service").
		WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
		Info("Getting canvas for project")

	// First check if project exists
	project, err := s.projectRepo.GetByUID(ctx, projectUID)
	if err != nil {
		logger.WithComponent("canvas-service").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Project not found")
		return nil, fmt.Errorf("project not found")
	}

	// Try to get existing canvas
	canvas, err := s.canvasRepo.GetByProjectUID(ctx, projectUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Create default canvas if none exists
			defaultCanvas := &models.BrainstormCanvas{
				ProjectID: project.ID,
				StateJSON: `{"nodes":[],"edges":[],"viewport":{"x":0,"y":0,"zoom":1}}`,
			}

			err = s.canvasRepo.Create(ctx, defaultCanvas)
			if err != nil {
				logger.WithComponent("canvas-service").
					WithFields(map[string]interface{}{
						"project_uid": projectUID.String(),
						"error":       err.Error(),
					}).
					Error("Failed to create default canvas")
				return nil, err
			}

			canvas = defaultCanvas
			logger.WithComponent("canvas-service").
				WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
				Info("Created default canvas for project")
		} else {
			logger.WithComponent("canvas-service").
				WithFields(map[string]interface{}{
					"project_uid": projectUID.String(),
					"error":       err.Error(),
				}).
				Error("Failed to get canvas")
			return nil, err
		}
	}

	response := &models.CanvasResponse{
		CanvasUID: canvas.CanvasUID,
		ProjectID: canvas.ProjectID,
		StateJSON: canvas.StateJSON,
		CreatedAt: canvas.CreatedAt,
		UpdatedAt: canvas.UpdatedAt,
	}

	logger.WithComponent("canvas-service").
		WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
		Info("Successfully retrieved canvas")

	return response, nil
}

// UpdateCanvas updates the canvas state for a project
func (s *CanvasService) UpdateCanvas(ctx context.Context, projectUID uuid.UUID, request *models.CanvasRequest) error {
	logger.WithComponent("canvas-service").
		WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
		Info("Updating canvas for project")

	// Validate JSON
	var jsonData interface{}
	if err := json.Unmarshal([]byte(request.StateJSON), &jsonData); err != nil {
		logger.WithComponent("canvas-service").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Invalid JSON in canvas state")
		return fmt.Errorf("invalid JSON format")
	}

	// First check if project exists
	_, err := s.projectRepo.GetByUID(ctx, projectUID)
	if err != nil {
		logger.WithComponent("canvas-service").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Project not found")
		return fmt.Errorf("project not found")
	}

	// Check if canvas exists, if not create it
	_, err = s.canvasRepo.GetByProjectUID(ctx, projectUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Canvas doesn't exist, create it
			project, _ := s.projectRepo.GetByUID(ctx, projectUID)
			canvas := &models.BrainstormCanvas{
				ProjectID: project.ID,
				StateJSON: request.StateJSON,
			}

			err = s.canvasRepo.Create(ctx, canvas)
			if err != nil {
				logger.WithComponent("canvas-service").
					WithFields(map[string]interface{}{
						"project_uid": projectUID.String(),
						"error":       err.Error(),
					}).
					Error("Failed to create canvas")
				return err
			}

			logger.WithComponent("canvas-service").
				WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
				Info("Created new canvas with updated state")

			return nil
		} else {
			logger.WithComponent("canvas-service").
				WithFields(map[string]interface{}{
					"project_uid": projectUID.String(),
					"error":       err.Error(),
				}).
				Error("Failed to check canvas existence")
			return err
		}
	}

	// Update existing canvas
	err = s.canvasRepo.Update(ctx, projectUID, request.StateJSON)
	if err != nil {
		logger.WithComponent("canvas-service").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Failed to update canvas")
		return err
	}

	logger.WithComponent("canvas-service").
		WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
		Info("Successfully updated canvas")

	return nil
}

// DeleteCanvas deletes the canvas for a project
func (s *CanvasService) DeleteCanvas(ctx context.Context, projectUID uuid.UUID) error {
	logger.WithComponent("canvas-service").
		WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
		Info("Deleting canvas for project")

	// First check if project exists
	_, err := s.projectRepo.GetByUID(ctx, projectUID)
	if err != nil {
		logger.WithComponent("canvas-service").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Project not found")
		return fmt.Errorf("project not found")
	}

	err = s.canvasRepo.Delete(ctx, projectUID)
	if err != nil {
		logger.WithComponent("canvas-service").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Failed to delete canvas")
		return err
	}

	logger.WithComponent("canvas-service").
		WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
		Info("Successfully deleted canvas")

	return nil
}
