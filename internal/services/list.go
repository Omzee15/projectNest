package services

import (
	"context"

	"github.com/google/uuid"

	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/internal/repositories"
	"lucid-lists-backend/internal/utils"
)

type ListService struct {
	listRepo    repositories.ListRepository
	taskRepo    repositories.TaskRepository
	projectRepo repositories.ProjectRepository
}

func NewListService(listRepo repositories.ListRepository, taskRepo repositories.TaskRepository, projectRepo repositories.ProjectRepository) *ListService {
	return &ListService{
		listRepo:    listRepo,
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
	}
}

func (s *ListService) CreateList(ctx context.Context, req *models.ListRequest) (*models.ListResponse, error) {
	// Get the project to verify it exists and get its internal ID
	project, err := s.projectRepo.GetByUID(ctx, req.ProjectUID)
	if err != nil {
		if err.Error() == "project not found" {
			return nil, utils.NewNotFoundError("Project not found")
		}
		return nil, utils.NewInternalError("Failed to get project: " + err.Error())
	}

	// Get max position for the project if position is not specified
	if req.Position == 0 {
		maxPosition, err := s.listRepo.GetMaxPositionByProject(ctx, project.ID)
		if err != nil {
			return nil, utils.NewInternalError("Failed to get max position: " + err.Error())
		}
		req.Position = maxPosition + 1
	}

	// Set default color if not provided
	color := req.Color
	if color == "" {
		color = "#FFFFFF"
	}

	// Create list model
	list := &models.List{
		ListUID:   uuid.New(),
		ProjectID: project.ID,
		Name:      req.Name,
		Color:     color,
		Position:  req.Position,
		IsActive:  true,
		CreatedBy: nil, // No user authentication yet
	}

	if err := s.listRepo.Create(ctx, list); err != nil {
		return nil, utils.NewInternalError("Failed to create list: " + err.Error())
	}

	return &models.ListResponse{
		ListUID:   list.ListUID,
		Name:      list.Name,
		Color:     list.Color,
		Position:  list.Position,
		CreatedAt: list.CreatedAt,
		UpdatedAt: list.UpdatedAt,
	}, nil
}

func (s *ListService) UpdateList(ctx context.Context, uid uuid.UUID, req *models.ListRequest) (*models.ListResponse, error) {
	// Check if list exists
	_, err := s.listRepo.GetByUID(ctx, uid)
	if err != nil {
		if err.Error() == "list not found" {
			return nil, utils.NewNotFoundError("List not found")
		}
		return nil, utils.NewInternalError("Failed to get list")
	}

	// Update list fields
	list := &models.List{
		Name:  req.Name,
		Color: req.Color,
	}

	if err := s.listRepo.Update(ctx, uid, list); err != nil {
		if err.Error() == "list not found" {
			return nil, utils.NewNotFoundError("List not found")
		}
		return nil, utils.NewInternalError("Failed to update list")
	}

	// Get updated list
	updatedList, err := s.listRepo.GetByUID(ctx, uid)
	if err != nil {
		return nil, utils.NewInternalError("Failed to get updated list")
	}

	return &models.ListResponse{
		ListUID:   updatedList.ListUID,
		Name:      updatedList.Name,
		Position:  updatedList.Position,
		CreatedAt: updatedList.CreatedAt,
		UpdatedAt: updatedList.UpdatedAt,
	}, nil
}

func (s *ListService) DeleteList(ctx context.Context, uid uuid.UUID) error {
	if err := s.listRepo.Delete(ctx, uid); err != nil {
		if err.Error() == "list not found" {
			return utils.NewNotFoundError("List not found")
		}
		return utils.NewInternalError("Failed to delete list")
	}

	return nil
}

func (s *ListService) UpdatePosition(ctx context.Context, uid uuid.UUID, position int) (*models.ListResponse, error) {
	if err := s.listRepo.UpdatePosition(ctx, uid, position); err != nil {
		if err.Error() == "list not found" {
			return nil, utils.NewNotFoundError("List not found")
		}
		return nil, utils.NewInternalError("Failed to update list position")
	}

	// Get updated list
	updatedList, err := s.listRepo.GetByUID(ctx, uid)
	if err != nil {
		return nil, utils.NewInternalError("Failed to get updated list")
	}

	return &models.ListResponse{
		ListUID:   updatedList.ListUID,
		Name:      updatedList.Name,
		Color:     updatedList.Color,
		Position:  updatedList.Position,
		CreatedAt: updatedList.CreatedAt,
		UpdatedAt: updatedList.UpdatedAt,
	}, nil
}

func (s *ListService) PartialUpdateList(ctx context.Context, uid uuid.UUID, updates *models.ListUpdateRequest) (*models.ListResponse, error) {
	// Check if list exists
	_, err := s.listRepo.GetByUID(ctx, uid)
	if err != nil {
		if err.Error() == "list not found" {
			return nil, utils.NewNotFoundError("List not found")
		}
		return nil, utils.NewInternalError("Failed to get list")
	}

	// Perform partial update
	if err := s.listRepo.PartialUpdate(ctx, uid, *updates); err != nil {
		if err.Error() == "list not found" {
			return nil, utils.NewNotFoundError("List not found")
		}
		if err.Error() == "no fields to update" {
			return nil, utils.NewBadRequestError("No fields to update")
		}
		return nil, utils.NewInternalError("Failed to update list")
	}

	// Get updated list
	updatedList, err := s.listRepo.GetByUID(ctx, uid)
	if err != nil {
		return nil, utils.NewInternalError("Failed to get updated list")
	}

	return &models.ListResponse{
		ListUID:   updatedList.ListUID,
		Name:      updatedList.Name,
		Color:     updatedList.Color,
		Position:  updatedList.Position,
		CreatedAt: updatedList.CreatedAt,
		UpdatedAt: updatedList.UpdatedAt,
	}, nil
}
