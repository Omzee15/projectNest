package repositories

import (
	"context"

	"github.com/google/uuid"

	"lucid-lists-backend/internal/models"
)

// ProjectRepository defines the interface for project data operations
type ProjectRepository interface {
	GetAll(ctx context.Context, userID int) ([]models.Project, error)
	GetByID(ctx context.Context, id int) (*models.Project, error)
	GetByUID(ctx context.Context, uid uuid.UUID) (*models.Project, error)
	GetWithLists(ctx context.Context, uid uuid.UUID) (*models.ProjectWithListsResponse, error)
	Create(ctx context.Context, project *models.Project) error
	Update(ctx context.Context, uid uuid.UUID, project *models.Project) error
	PartialUpdate(ctx context.Context, uid uuid.UUID, updates models.ProjectUpdateRequest) error
	Delete(ctx context.Context, uid uuid.UUID) error
	AddMember(ctx context.Context, projectID int, userUID uuid.UUID, role string) error
	GetMembers(ctx context.Context, projectID int) ([]models.ProjectMember, error)
	IsMember(ctx context.Context, projectID int, userUID uuid.UUID) (bool, error)
}

// ProjectMemberRepository defines the interface for project member operations
type ProjectMemberRepository interface {
	AddMember(ctx context.Context, projectID int, userUID uuid.UUID, role string) error
	RemoveMember(ctx context.Context, projectID int, userUID uuid.UUID) error
	GetMembersByProjectID(ctx context.Context, projectID int) ([]models.ProjectMember, error)
	IsMember(ctx context.Context, projectID int, userUID uuid.UUID) (bool, error)
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
	// Progress tracking methods for Phase 2
	CountByProjectUID(ctx context.Context, projectUID uuid.UUID) (int64, error)
	CountCompletedByProjectUID(ctx context.Context, projectUID uuid.UUID) (int64, error)
}

// Phase 3: Brainstorming & Planning Layer Repository Interfaces

// CanvasRepository defines the interface for brainstorm canvas data operations
type CanvasRepository interface {
	GetByProjectUID(ctx context.Context, projectUID uuid.UUID) (*models.BrainstormCanvas, error)
	Create(ctx context.Context, canvas *models.BrainstormCanvas) error
	Update(ctx context.Context, projectUID uuid.UUID, stateJSON string) error
	Delete(ctx context.Context, projectUID uuid.UUID) error
}

// NoteRepository defines the interface for note data operations
type NoteRepository interface {
	GetByProjectUID(ctx context.Context, projectUID uuid.UUID) ([]models.Note, error)
	GetByUID(ctx context.Context, uid uuid.UUID) (*models.Note, error)
	Create(ctx context.Context, note *models.Note) error
	Update(ctx context.Context, uid uuid.UUID, note *models.Note) error
	PartialUpdate(ctx context.Context, uid uuid.UUID, updates models.NoteUpdateRequest) error
	Delete(ctx context.Context, uid uuid.UUID) error
	GetMaxPositionByProject(ctx context.Context, projectID int) (int, error)
	UpdateFolderID(ctx context.Context, uid uuid.UUID, folderID *int) error
}

// NoteFolderRepository defines the interface for note folder data operations
type NoteFolderRepositoryInterface interface {
	GetByProjectUID(ctx context.Context, projectUID uuid.UUID) ([]models.NoteFolder, error)
	Create(ctx context.Context, projectID int, folder models.NoteFolder) (*models.NoteFolder, error)
	Update(ctx context.Context, folderUID uuid.UUID, updates models.NoteFolder) (*models.NoteFolder, error)
	Delete(ctx context.Context, folderUID uuid.UUID) error
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByUID(ctx context.Context, uid uuid.UUID) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
}

// ChatRepositoryInterface defines the interface for chat conversation operations
type ChatRepositoryInterface interface {
	GetConversationsByProjectID(projectID int) ([]models.ChatConversation, error)
	GetConversationByUID(conversationUID uuid.UUID) (*models.ChatConversation, error)
	CreateConversation(conversation *models.ChatConversation) error
	UpdateConversation(conversation *models.ChatConversation) error
	DeleteConversation(conversationUID uuid.UUID) error
	GetMessagesByConversationID(conversationID int) ([]models.ChatMessage, error)
	CreateMessage(message *models.ChatMessage) error
}

// UserSettingsRepository defines the interface for user settings operations
type UserSettingsRepository interface {
	GetByUserID(ctx context.Context, userID int) (*models.UserSettings, error)
	Create(ctx context.Context, settings *models.UserSettings) error
	Update(ctx context.Context, userID int, settings *models.UserSettingsRequest) error
	CreateOrUpdate(ctx context.Context, userID int, settings *models.UserSettingsRequest) error
	DeleteByUserID(ctx context.Context, userID int) error
}
