package repositories

import (
	"context"

	"github.com/google/uuid"

	"lucid-lists-backend/internal/models"
)

// ProjectRepository defines the interface for project data operations
type ProjectRepository interface {
	GetAll(ctx context.Context) ([]models.Project, error)
	GetByUID(ctx context.Context, uid uuid.UUID) (*models.Project, error)
	GetWithLists(ctx context.Context, uid uuid.UUID) (*models.ProjectWithListsResponse, error)
	Create(ctx context.Context, project *models.Project) error
	Update(ctx context.Context, uid uuid.UUID, project *models.Project) error
	PartialUpdate(ctx context.Context, uid uuid.UUID, updates models.ProjectUpdateRequest) error
	Delete(ctx context.Context, uid uuid.UUID) error
	GetMaxPositionByWorkspace(ctx context.Context, workspaceID int) (int, error)
}

// ListRepository defines the interface for list data operations  
type ListRepository interface {
	GetByProjectID(ctx context.Context, projectID int) ([]models.List, error)
	GetByUID(ctx context.Context, uid uuid.UUID) (*models.List, error)
	Create(ctx context.Context, list *models.List) error
	Update(ctx context.Context, uid uuid.UUID, list *models.List) error
	PartialUpdate(ctx context.Context, uid uuid.UUID, updates models.ListUpdateRequest) error
	Delete(ctx context.Context, uid uuid.UUID) error
	UpdatePosition(ctx context.Context, uid uuid.UUID, position int) error
	GetMaxPositionByProject(ctx context.Context, projectID int) (int, error)
}

// TaskRepository defines the interface for task data operations
type TaskRepository interface {
	GetByListID(ctx context.Context, listID int) ([]models.Task, error)
	GetByUID(ctx context.Context, uid uuid.UUID) (*models.Task, error)
	Create(ctx context.Context, task *models.Task) error
	Update(ctx context.Context, uid uuid.UUID, task *models.Task) error
	PartialUpdate(ctx context.Context, uid uuid.UUID, updates models.TaskUpdateRequest) error
	Delete(ctx context.Context, uid uuid.UUID) error
	MoveToList(ctx context.Context, uid uuid.UUID, newListID int) error
	GetByProjectID(ctx context.Context, projectID int) ([]models.Task, error)
	GetMaxPositionByList(ctx context.Context, listID int) (int, error)
}