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
