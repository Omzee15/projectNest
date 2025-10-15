package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/internal/repositories"
	"lucid-lists-backend/pkg/logger"
)

type NoteFolderService struct {
	folderRepo  repositories.NoteFolderRepositoryInterface
	projectRepo repositories.ProjectRepository
}

func NewNoteFolderService(folderRepo repositories.NoteFolderRepositoryInterface, projectRepo repositories.ProjectRepository) *NoteFolderService {
	return &NoteFolderService{
		folderRepo:  folderRepo,
		projectRepo: projectRepo,
	}
}

// GetFoldersByProject retrieves all folders for a project
func (s *NoteFolderService) GetFoldersByProject(ctx context.Context, projectUID uuid.UUID) (*models.NoteFoldersResponse, error) {
	logger.WithComponent("note-folder-service").
		WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
		Info("Getting folders for project")

	// Get folders from repository
	folders, err := s.folderRepo.GetByProjectUID(ctx, projectUID)
	if err != nil {
		logger.WithComponent("note-folder-service").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Failed to get folders")
		return nil, fmt.Errorf("failed to get folders: %w", err)
	}

	// Convert to response DTOs
	var folderResponses []models.NoteFolderResponse
	for _, folder := range folders {
		folderResponse := models.NoteFolderResponse{
			ID:             folder.ID,
			FolderUID:      folder.FolderUID,
			ProjectID:      folder.ProjectID,
			ParentFolderID: folder.ParentFolderID,
			Name:           folder.Name,
			Position:       folder.Position,
			CreatedAt:      folder.CreatedAt,
			UpdatedAt:      folder.UpdatedAt,
		}
		folderResponses = append(folderResponses, folderResponse)
	}

	response := &models.NoteFoldersResponse{
		Folders: folderResponses,
		Total:   len(folderResponses),
	}

	logger.WithComponent("note-folder-service").
		WithFields(map[string]interface{}{
			"project_uid": projectUID.String(),
			"count":       len(folderResponses),
		}).
		Info("Successfully retrieved folders")

	return response, nil
}

// CreateFolder creates a new folder
func (s *NoteFolderService) CreateFolder(ctx context.Context, projectUID uuid.UUID, request models.NoteFolderRequest) (*models.NoteFolderResponse, error) {
	logger.WithComponent("note-folder-service").
		WithFields(map[string]interface{}{
			"project_uid": projectUID.String(),
			"name":        request.Name,
		}).
		Info("Creating folder")

	// Get project to validate it exists and get project ID
	project, err := s.projectRepo.GetByUID(ctx, projectUID)
	if err != nil {
		logger.WithComponent("note-folder-service").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Failed to get project")
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Create folder model
	folder := models.NoteFolder{
		Name:           request.Name,
		ParentFolderID: request.ParentFolderID,
		Position:       request.Position,
	}

	// Create folder in repository
	createdFolder, err := s.folderRepo.Create(ctx, project.ID, folder)
	if err != nil {
		logger.WithComponent("note-folder-service").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"name":        request.Name,
				"error":       err.Error(),
			}).
			Error("Failed to create folder")
		return nil, fmt.Errorf("failed to create folder: %w", err)
	}

	// Convert to response DTO
	response := &models.NoteFolderResponse{
		ID:             createdFolder.ID,
		FolderUID:      createdFolder.FolderUID,
		ProjectID:      createdFolder.ProjectID,
		ParentFolderID: createdFolder.ParentFolderID,
		Name:           createdFolder.Name,
		Position:       createdFolder.Position,
		CreatedAt:      createdFolder.CreatedAt,
		UpdatedAt:      createdFolder.UpdatedAt,
	}

	logger.WithComponent("note-folder-service").
		WithFields(map[string]interface{}{
			"folder_uid": createdFolder.FolderUID.String(),
			"name":       createdFolder.Name,
		}).
		Info("Successfully created folder")

	return response, nil
}

// UpdateFolder updates an existing folder
func (s *NoteFolderService) UpdateFolder(ctx context.Context, folderUID uuid.UUID, request models.NoteFolderUpdateRequest) (*models.NoteFolderResponse, error) {
	logger.WithComponent("note-folder-service").
		WithFields(map[string]interface{}{"folder_uid": folderUID.String()}).
		Info("Updating folder")

	// Create update model
	updates := models.NoteFolder{
		Name:           *request.Name, // Assuming name is required for update
		ParentFolderID: request.ParentFolderID,
		Position:       request.Position,
	}

	// Update folder in repository
	updatedFolder, err := s.folderRepo.Update(ctx, folderUID, updates)
	if err != nil {
		logger.WithComponent("note-folder-service").
			WithFields(map[string]interface{}{
				"folder_uid": folderUID.String(),
				"error":      err.Error(),
			}).
			Error("Failed to update folder")
		return nil, fmt.Errorf("failed to update folder: %w", err)
	}

	// Convert to response DTO
	response := &models.NoteFolderResponse{
		ID:             updatedFolder.ID,
		FolderUID:      updatedFolder.FolderUID,
		ProjectID:      updatedFolder.ProjectID,
		ParentFolderID: updatedFolder.ParentFolderID,
		Name:           updatedFolder.Name,
		Position:       updatedFolder.Position,
		CreatedAt:      updatedFolder.CreatedAt,
		UpdatedAt:      updatedFolder.UpdatedAt,
	}

	logger.WithComponent("note-folder-service").
		WithFields(map[string]interface{}{"folder_uid": folderUID.String()}).
		Info("Successfully updated folder")

	return response, nil
}

// DeleteFolder deletes a folder
func (s *NoteFolderService) DeleteFolder(ctx context.Context, folderUID uuid.UUID) error {
	logger.WithComponent("note-folder-service").
		WithFields(map[string]interface{}{"folder_uid": folderUID.String()}).
		Info("Deleting folder")

	// Delete folder from repository
	err := s.folderRepo.Delete(ctx, folderUID)
	if err != nil {
		logger.WithComponent("note-folder-service").
			WithFields(map[string]interface{}{
				"folder_uid": folderUID.String(),
				"error":      err.Error(),
			}).
			Error("Failed to delete folder")
		return fmt.Errorf("failed to delete folder: %w", err)
	}

	logger.WithComponent("note-folder-service").
		WithFields(map[string]interface{}{"folder_uid": folderUID.String()}).
		Info("Successfully deleted folder")

	return nil
}
