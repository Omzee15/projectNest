package services

import (
	"context"

	"github.com/google/uuid"

	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/internal/repositories"
	"lucid-lists-backend/internal/utils"
)

type ProjectService struct {
	projectRepo repositories.ProjectRepository
	listRepo    repositories.ListRepository
	taskRepo    repositories.TaskRepository
}

func NewProjectService(projectRepo repositories.ProjectRepository, listRepo repositories.ListRepository, taskRepo repositories.TaskRepository) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		listRepo:    listRepo,
		taskRepo:    taskRepo,
	}
}

func (s *ProjectService) GetAllProjects(ctx context.Context) ([]models.ProjectResponse, error) {
	projects, err := s.projectRepo.GetAll(ctx)
	if err != nil {
		return nil, utils.NewInternalError("Failed to retrieve projects")
	}

	var response []models.ProjectResponse
	for _, project := range projects {
		response = append(response, models.ProjectResponse{
			ProjectUID:  project.ProjectUID,
			Name:        project.Name,
			Description: project.Description,
			Status:      project.Status,
			Color:       project.Color,
			Position:    project.Position,
			StartDate:   project.StartDate,
			EndDate:     project.EndDate,
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt,
		})
	}

	return response, nil
}

func (s *ProjectService) GetProjectWithLists(ctx context.Context, uid uuid.UUID) (*models.ProjectWithListsResponse, error) {
	projectWithLists, err := s.projectRepo.GetWithLists(ctx, uid)
	if err != nil {
		if err.Error() == "project not found" {
			return nil, utils.NewNotFoundError("Project not found")
		}
		// Temporarily show the actual error for debugging
		return nil, utils.NewInternalError("Failed to get project: " + err.Error())
	}

	return projectWithLists, nil
}

func (s *ProjectService) CreateProject(ctx context.Context, req *models.ProjectRequest) (*models.ProjectResponse, error) {
	// Get next position if not specified
	position := req.Position
	if position == nil {
		// For now, we'll assume workspace_id = 1 since we don't have user management yet
		maxPos, err := s.projectRepo.GetMaxPositionByWorkspace(ctx, 1)
		if err == nil {
			newPos := maxPos + 1
			position = &newPos
		}
	}

	// Set default color if not provided
	color := req.Color
	if color == "" {
		color = "#FFFFFF"
	}

	// Create project model
	project := &models.Project{
		ProjectUID:  uuid.New(),
		WorkspaceID: 1, // Default workspace
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		Color:       color,
		Position:    position,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		IsActive:    true,
		CreatedBy:   nil, // No user authentication yet
	}

	if req.Status == "" {
		project.Status = "active"
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, utils.NewInternalError("Failed to create project")
	}

	return &models.ProjectResponse{
		ProjectUID:  project.ProjectUID,
		Name:        project.Name,
		Description: project.Description,
		Status:      project.Status,
		Color:       project.Color,
		Position:    project.Position,
		StartDate:   project.StartDate,
		EndDate:     project.EndDate,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, uid uuid.UUID, req *models.ProjectRequest) (*models.ProjectResponse, error) {
	// Check if project exists
	_, err := s.projectRepo.GetByUID(ctx, uid)
	if err != nil {
		if err.Error() == "project not found" {
			return nil, utils.NewNotFoundError("Project not found")
		}
		return nil, utils.NewInternalError("Failed to get project")
	}

	// Set default color if not provided
	color := req.Color
	if color == "" {
		color = "#FFFFFF"
	}

	// Update project fields
	project := &models.Project{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		Color:       color,
		Position:    req.Position,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		UpdatedBy:   nil, // No user authentication yet
	}

	if req.Status == "" {
		project.Status = "active"
	}

	if err := s.projectRepo.Update(ctx, uid, project); err != nil {
		return nil, utils.NewInternalError("Failed to update project")
	}

	// Get updated project to return
	updatedProject, err := s.projectRepo.GetByUID(ctx, uid)
	if err != nil {
		return nil, utils.NewInternalError("Failed to get updated project")
	}

	return &models.ProjectResponse{
		ProjectUID:  updatedProject.ProjectUID,
		Name:        updatedProject.Name,
		Description: updatedProject.Description,
		Status:      updatedProject.Status,
		Color:       updatedProject.Color,
		Position:    updatedProject.Position,
		StartDate:   updatedProject.StartDate,
		EndDate:     updatedProject.EndDate,
		CreatedAt:   updatedProject.CreatedAt,
		UpdatedAt:   updatedProject.UpdatedAt,
	}, nil
}

func (s *ProjectService) PartialUpdateProject(ctx context.Context, uid uuid.UUID, updates *models.ProjectUpdateRequest) (*models.ProjectResponse, error) {
	// Check if project exists
	_, err := s.projectRepo.GetByUID(ctx, uid)
	if err != nil {
		if err.Error() == "project not found" {
			return nil, utils.NewNotFoundError("Project not found")
		}
		return nil, utils.NewInternalError("Failed to get project")
	}

	// Apply partial update
	if err := s.projectRepo.PartialUpdate(ctx, uid, *updates); err != nil {
		return nil, utils.NewInternalError("Failed to update project")
	}

	// Get updated project to return
	updatedProject, err := s.projectRepo.GetByUID(ctx, uid)
	if err != nil {
		return nil, utils.NewInternalError("Failed to get updated project")
	}

	return &models.ProjectResponse{
		ProjectUID:  updatedProject.ProjectUID,
		Name:        updatedProject.Name,
		Description: updatedProject.Description,
		Status:      updatedProject.Status,
		Color:       updatedProject.Color,
		Position:    updatedProject.Position,
		StartDate:   updatedProject.StartDate,
		EndDate:     updatedProject.EndDate,
		CreatedAt:   updatedProject.CreatedAt,
		UpdatedAt:   updatedProject.UpdatedAt,
	}, nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, uid uuid.UUID) error {
	// Check if project exists
	_, err := s.projectRepo.GetByUID(ctx, uid)
	if err != nil {
		if err.Error() == "project not found" {
			return utils.NewNotFoundError("Project not found")
		}
		return utils.NewInternalError("Failed to get project")
	}

	if err := s.projectRepo.Delete(ctx, uid); err != nil {
		return utils.NewInternalError("Failed to delete project")
	}

	return nil
}

// GetProjectProgress calculates and returns progress statistics for a project
func (s *ProjectService) GetProjectProgress(ctx context.Context, uid uuid.UUID) (*models.ProjectProgressResponse, error) {
	// Check if project exists
	_, err := s.projectRepo.GetByUID(ctx, uid)
	if err != nil {
		if err.Error() == "project not found" {
			return nil, utils.NewNotFoundError("Project not found")
		}
		return nil, utils.NewInternalError("Failed to get project")
	}

	// Count total tasks
	totalTasks, err := s.taskRepo.CountByProjectUID(ctx, uid)
	if err != nil {
		return nil, utils.NewInternalError("Failed to count total tasks")
	}

	// Count completed tasks (using is_completed field)
	completedTasks, err := s.taskRepo.CountCompletedByProjectUID(ctx, uid)
	if err != nil {
		return nil, utils.NewInternalError("Failed to count completed tasks")
	}

	// Calculate todo tasks (total - completed)
	todoTasks := totalTasks - completedTasks

	// Calculate progress percentage
	var progress float64
	if totalTasks > 0 {
		progress = float64(completedTasks) / float64(totalTasks)
	}

	return &models.ProjectProgressResponse{
		TotalTasks:     int(totalTasks),
		CompletedTasks: int(completedTasks),
		TodoTasks:      int(todoTasks),
		Progress:       progress,
	}, nil
}

// GetAllProjectsWithProgress returns all projects with their progress statistics
func (s *ProjectService) GetAllProjectsWithProgress(ctx context.Context) ([]models.ProjectWithProgressResponse, error) {
	projects, err := s.projectRepo.GetAll(ctx)
	if err != nil {
		return nil, utils.NewInternalError("Failed to retrieve projects")
	}

	var response []models.ProjectWithProgressResponse
	for _, project := range projects {
		// Get project response
		projectResponse := models.ProjectResponse{
			ProjectUID:  project.ProjectUID,
			Name:        project.Name,
			Description: project.Description,
			Status:      project.Status,
			Color:       project.Color,
			Position:    project.Position,
			StartDate:   project.StartDate,
			EndDate:     project.EndDate,
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt,
		}

		// Get progress stats
		progressStats, err := s.GetProjectProgress(ctx, project.ProjectUID)
		if err != nil {
			// If we can't get progress stats, return empty stats instead of failing
			progressStats = &models.ProjectProgressResponse{
				TotalTasks:     0,
				CompletedTasks: 0,
				TodoTasks:      0,
				Progress:       0,
			}
		}

		response = append(response, models.ProjectWithProgressResponse{
			ProjectResponse: projectResponse,
			TaskStats:       *progressStats,
		})
	}

	return response, nil
}
