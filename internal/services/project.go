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
	userRepo    repositories.UserRepository
}

func NewProjectService(projectRepo repositories.ProjectRepository, listRepo repositories.ListRepository, taskRepo repositories.TaskRepository, userRepo repositories.UserRepository) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		listRepo:    listRepo,
		taskRepo:    taskRepo,
		userRepo:    userRepo,
	}
}

func (s *ProjectService) GetAllProjects(ctx context.Context, userID int) ([]models.ProjectResponse, error) {
	projects, err := s.projectRepo.GetAll(ctx, userID)
	if err != nil {
		return nil, utils.NewInternalError("Failed to retrieve projects")
	}

	var response []models.ProjectResponse
	for _, project := range projects {
		response = append(response, models.ProjectResponse{
			ProjectUID:       project.ProjectUID,
			Name:             project.Name,
			Description:      project.Description,
			Status:           project.Status,
			Color:            project.Color,
			Position:         project.Position,
			StartDate:        project.StartDate,
			EndDate:          project.EndDate,
			IsPrivate:        project.IsPrivate,
			DbmlContent:      project.DbmlContent,
			FlowchartContent: project.FlowchartContent,
			CreatedAt:        project.CreatedAt,
			UpdatedAt:        project.UpdatedAt,
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

func (s *ProjectService) CreateProject(ctx context.Context, req *models.ProjectRequest, userID int, userUID uuid.UUID) (*models.ProjectResponse, error) {
	// Set default color if not provided
	color := req.Color
	if color == "" {
		color = "#FFFFFF"
	}

	// Create project model
	project := &models.Project{
		ProjectUID:       uuid.New(),
		UserID:           userID,
		Name:             req.Name,
		Description:      req.Description,
		Status:           req.Status,
		Color:            color,
		Position:         req.Position,
		StartDate:        req.StartDate,
		EndDate:          req.EndDate,
		IsPrivate:        req.IsPrivate != nil && *req.IsPrivate,
		DbmlContent:      req.DbmlContent,
		FlowchartContent: req.FlowchartContent,
		IsActive:         true,
		CreatedBy:        &userID,
	}

	if req.Status == "" {
		project.Status = "active"
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, utils.NewInternalError("Failed to create project")
	}

	// Add the creator as an owner in the project_member table
	if err := s.projectRepo.AddMember(ctx, project.ID, userUID, "owner"); err != nil {
		// Log the error but don't fail the project creation
		// In a production environment, you might want to handle this more carefully
		return nil, utils.NewInternalError("Failed to add project owner")
	}

	return &models.ProjectResponse{
		ProjectUID:       project.ProjectUID,
		Name:             project.Name,
		Description:      project.Description,
		Status:           project.Status,
		Color:            project.Color,
		Position:         project.Position,
		StartDate:        project.StartDate,
		EndDate:          project.EndDate,
		IsPrivate:        project.IsPrivate,
		DbmlContent:      project.DbmlContent,
		FlowchartContent: project.FlowchartContent,
		CreatedAt:        project.CreatedAt,
		UpdatedAt:        project.UpdatedAt,
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
		Name:             req.Name,
		Description:      req.Description,
		Status:           req.Status,
		Color:            color,
		Position:         req.Position,
		StartDate:        req.StartDate,
		EndDate:          req.EndDate,
		IsPrivate:        req.IsPrivate != nil && *req.IsPrivate,
		DbmlContent:      req.DbmlContent,
		FlowchartContent: req.FlowchartContent,
		UpdatedBy:        nil, // No user authentication yet
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
		IsPrivate:   updatedProject.IsPrivate,
		DbmlContent: updatedProject.DbmlContent,
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
		ProjectUID:       updatedProject.ProjectUID,
		Name:             updatedProject.Name,
		Description:      updatedProject.Description,
		Status:           updatedProject.Status,
		Color:            updatedProject.Color,
		Position:         updatedProject.Position,
		StartDate:        updatedProject.StartDate,
		EndDate:          updatedProject.EndDate,
		IsPrivate:        updatedProject.IsPrivate,
		DbmlContent:      updatedProject.DbmlContent,
		FlowchartContent: updatedProject.FlowchartContent,
		CreatedAt:        updatedProject.CreatedAt,
		UpdatedAt:        updatedProject.UpdatedAt,
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
func (s *ProjectService) GetAllProjectsWithProgress(ctx context.Context, userID int) ([]models.ProjectWithProgressResponse, error) {
	projects, err := s.projectRepo.GetAll(ctx, userID)
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

// AddMemberByEmail adds a user to a project by their email address
func (s *ProjectService) AddMemberByEmail(ctx context.Context, projectUID uuid.UUID, email string, role string) error {
	// Get the project to get its ID
	project, err := s.projectRepo.GetByUID(ctx, projectUID)
	if err != nil {
		if err.Error() == "project not found" {
			return utils.NewNotFoundError("Project not found")
		}
		return utils.NewInternalError("Failed to get project")
	}

	// Find the user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return utils.NewNotFoundError("user not found")
	}

	// Check if user is already a member
	isMember, err := s.projectRepo.IsMember(ctx, project.ID, user.UserUID)
	if err != nil {
		return utils.NewInternalError("Failed to check membership")
	}

	if isMember {
		return utils.NewBadRequestError("user already a member")
	}

	// Add the user as a member
	err = s.projectRepo.AddMember(ctx, project.ID, user.UserUID, role)
	if err != nil {
		return utils.NewInternalError("Failed to add member")
	}

	return nil
}

// GetProjectMembers returns all members of a project
func (s *ProjectService) GetProjectMembers(ctx context.Context, projectUID uuid.UUID) ([]models.ProjectMemberResponse, error) {
	// Get the project to get its ID
	project, err := s.projectRepo.GetByUID(ctx, projectUID)
	if err != nil {
		if err.Error() == "project not found" {
			return nil, utils.NewNotFoundError("Project not found")
		}
		return nil, utils.NewInternalError("Failed to get project")
	}

	// Get all members
	members, err := s.projectRepo.GetMembers(ctx, project.ID)
	if err != nil {
		return nil, utils.NewInternalError("Failed to get project members")
	}

	var response []models.ProjectMemberResponse
	for _, member := range members {
		// Get user details
		user, err := s.userRepo.GetByID(ctx, member.UserID)
		if err != nil {
			continue // Skip if user not found
		}

		response = append(response, models.ProjectMemberResponse{
			UserUID:  user.UserUID,
			Email:    user.Email,
			Name:     user.Name,
			Role:     member.Role,
			JoinedAt: member.JoinedAt,
		})
	}

	return response, nil
}
