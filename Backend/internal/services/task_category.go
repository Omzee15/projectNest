package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/internal/repositories"
)

type TaskCategoryService struct {
	categoryRepo *repositories.TaskCategoryRepository
	projectRepo  repositories.ProjectRepository
	taskRepo     repositories.TaskRepository
}

func NewTaskCategoryService(
	categoryRepo *repositories.TaskCategoryRepository,
	projectRepo repositories.ProjectRepository,
	taskRepo repositories.TaskRepository,
) *TaskCategoryService {
	return &TaskCategoryService{
		categoryRepo: categoryRepo,
		projectRepo:  projectRepo,
		taskRepo:     taskRepo,
	}
}

// CreateCategory creates a new task category
func (s *TaskCategoryService) CreateCategory(ctx context.Context, req *models.TaskCategoryRequest, userID int) (*models.TaskCategoryResponse, error) {
	// Get project by UID
	project, err := s.projectRepo.GetByUID(ctx, req.ProjectUID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// Set default color if not provided
	color := req.Color
	if color == "" {
		color = "#808080"
	}

	// Create category
	category := &models.TaskCategory{
		ProjectID:   project.ID,
		Name:        req.Name,
		Color:       color,
		Description: req.Description,
		CreatedBy:   &userID,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	// Build response
	response := &models.TaskCategoryResponse{
		CategoryUID: category.CategoryUID,
		ProjectUID:  project.ProjectUID,
		Name:        category.Name,
		Color:       category.Color,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}

	return response, nil
}

// GetCategoriesByProjectUID retrieves all categories for a project
func (s *TaskCategoryService) GetCategoriesByProjectUID(ctx context.Context, projectUID uuid.UUID) ([]models.TaskCategoryResponse, error) {
	// Get project
	project, err := s.projectRepo.GetByUID(ctx, projectUID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// Get categories
	categories, err := s.categoryRepo.GetByProjectID(ctx, project.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	// Build responses
	responses := make([]models.TaskCategoryResponse, 0, len(categories))
	for _, category := range categories {
		responses = append(responses, models.TaskCategoryResponse{
			CategoryUID: category.CategoryUID,
			ProjectUID:  project.ProjectUID,
			Name:        category.Name,
			Color:       category.Color,
			Description: category.Description,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		})
	}

	return responses, nil
}

// UpdateCategory updates a category
func (s *TaskCategoryService) UpdateCategory(ctx context.Context, categoryUID uuid.UUID, req *models.TaskCategoryUpdateRequest, userID int) (*models.TaskCategoryResponse, error) {
	// Get category
	category, err := s.categoryRepo.GetByUID(ctx, categoryUID)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	// Update fields
	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Color != nil {
		category.Color = *req.Color
	}
	if req.Description != nil {
		category.Description = req.Description
	}
	category.UpdatedBy = &userID

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	// Get project info
	project, err := s.projectRepo.GetByID(ctx, category.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Build response
	response := &models.TaskCategoryResponse{
		CategoryUID: category.CategoryUID,
		ProjectUID:  project.ProjectUID,
		Name:        category.Name,
		Color:       category.Color,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}

	return response, nil
}

// DeleteCategory deletes a category
func (s *TaskCategoryService) DeleteCategory(ctx context.Context, categoryUID uuid.UUID) error {
	// Get category
	_, err := s.categoryRepo.GetByUID(ctx, categoryUID)
	if err != nil {
		return fmt.Errorf("category not found: %w", err)
	}

	// Delete category
	if err := s.categoryRepo.Delete(ctx, categoryUID); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

// AssignCategoriesToTask assigns multiple categories to a task
func (s *TaskCategoryService) AssignCategoriesToTask(ctx context.Context, taskUID uuid.UUID, categoryUIDs []uuid.UUID, userID int) error {
	// Get task
	task, err := s.taskRepo.GetByUID(ctx, taskUID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// Clear existing categories
	if err := s.categoryRepo.ClearTaskCategories(ctx, task.ID); err != nil {
		return fmt.Errorf("failed to clear existing categories: %w", err)
	}

	// Assign new categories
	for _, categoryUID := range categoryUIDs {
		category, err := s.categoryRepo.GetByUID(ctx, categoryUID)
		if err != nil {
			continue // Skip invalid categories
		}

		if err := s.categoryRepo.AssignToTask(ctx, task.ID, category.ID, &userID); err != nil {
			return fmt.Errorf("failed to assign category: %w", err)
		}
	}

	return nil
}

// GetCategoriesByTaskUID retrieves all categories for a task
func (s *TaskCategoryService) GetCategoriesByTaskUID(ctx context.Context, taskUID uuid.UUID) ([]models.TaskCategoryResponse, error) {
	// Get task
	task, err := s.taskRepo.GetByUID(ctx, taskUID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// Get categories
	categories, err := s.categoryRepo.GetCategoriesByTaskID(ctx, task.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	// Build responses - we need to get the project UID from the task's list
	responses := make([]models.TaskCategoryResponse, 0, len(categories))
	for _, category := range categories {
		// Get project for this category
		project, err := s.projectRepo.GetByID(ctx, category.ProjectID)
		if err != nil {
			continue // Skip if project not found
		}

		responses = append(responses, models.TaskCategoryResponse{
			CategoryUID: category.CategoryUID,
			ProjectUID:  project.ProjectUID,
			Name:        category.Name,
			Color:       category.Color,
			Description: category.Description,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		})
	}

	return responses, nil
}

// RemoveCategoryFromTask removes a category from a task
func (s *TaskCategoryService) RemoveCategoryFromTask(ctx context.Context, taskUID, categoryUID uuid.UUID) error {
	// Get task
	task, err := s.taskRepo.GetByUID(ctx, taskUID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// Get category
	category, err := s.categoryRepo.GetByUID(ctx, categoryUID)
	if err != nil {
		return fmt.Errorf("category not found: %w", err)
	}

	// Remove category from task
	if err := s.categoryRepo.RemoveFromTask(ctx, task.ID, category.ID); err != nil {
		return fmt.Errorf("failed to remove category from task: %w", err)
	}

	return nil
}
