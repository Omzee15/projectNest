package models

import (
	"time"

	"github.com/google/uuid"
)

// Database models - these map directly to the database schema
// These use integer IDs internally for efficiency

type User struct {
	ID        int       `db:"id"`
	UserUID   uuid.UUID `db:"user_uid"`
	Email     string    `db:"email"`
	Password  string    `db:"password_hash"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	IsActive  bool      `db:"is_active"`
}

type Workspace struct {
	ID           int        `db:"id"`
	WorkspaceUID uuid.UUID  `db:"workspace_uid"`
	Name         string     `db:"name"`
	Description  *string    `db:"description"`
	CreatedAt    time.Time  `db:"created_at"`
	CreatedBy    *uuid.UUID `db:"created_by"`
	UpdatedAt    *time.Time `db:"updated_at"`
	UpdatedBy    *uuid.UUID `db:"updated_by"`
	IsActive     bool       `db:"is_active"`
}

type Project struct {
	ID               int        `db:"id"`
	ProjectUID       uuid.UUID  `db:"project_uid"`
	UserID           int        `db:"user_id"`
	Name             string     `db:"name"`
	Description      *string    `db:"description"`
	Status           string     `db:"status"`
	Color            string     `db:"color"`
	Position         *int       `db:"position"`
	StartDate        *time.Time `db:"start_date"`
	EndDate          *time.Time `db:"end_date"`
	IsPrivate        bool       `db:"is_private"`
	DbmlContent      *string    `db:"dbml_content"`
	DbmlLayoutData   *string    `db:"dbml_layout_data"`
	FlowchartContent *string    `db:"flowchart_content"`
	CreatedAt        time.Time  `db:"created_at"`
	CreatedBy        *uuid.UUID `db:"created_by"`
	UpdatedAt        *time.Time `db:"updated_at"`
	UpdatedBy        *uuid.UUID `db:"updated_by"`
	IsActive         bool       `db:"is_active"`
}

type ProjectMember struct {
	ID        int       `db:"id"`
	ProjectID int       `db:"project_id"`
	UserID    uuid.UUID `db:"user_id"`
	Role      string    `db:"role"` // 'owner' or 'member'
	JoinedAt  time.Time `db:"joined_at"`
}

type List struct {
	ID        int        `db:"id"`
	ListUID   uuid.UUID  `db:"list_uid"`
	ProjectID int        `db:"project_id"`
	Name      string     `db:"name"`
	Color     string     `db:"color"`
	Position  int        `db:"position"`
	CreatedAt time.Time  `db:"created_at"`
	CreatedBy *uuid.UUID `db:"created_by"`
	UpdatedAt *time.Time `db:"updated_at"`
	UpdatedBy *uuid.UUID `db:"updated_by"`
	IsActive  bool       `db:"is_active"`
}

type Task struct {
	ID          int        `db:"id"`
	TaskUID     uuid.UUID  `db:"task_uid"`
	ListID      int        `db:"list_id"`
	Title       string     `db:"title"`
	Description *string    `db:"description"`
	Priority    *string    `db:"priority"`
	Status      string     `db:"status"`
	Color       string     `db:"color"`
	Position    *int       `db:"position"`
	IsCompleted bool       `db:"is_completed"`
	DueDate     *time.Time `db:"due_date"`
	CompletedAt *time.Time `db:"completed_at"`
	CreatedAt   time.Time  `db:"created_at"`
	CreatedBy   *uuid.UUID `db:"created_by"`
	UpdatedAt   *time.Time `db:"updated_at"`
	UpdatedBy   *uuid.UUID `db:"updated_by"`
	IsActive    bool       `db:"is_active"`
}

// Phase 3: Brainstorming & Planning Layer Models

type BrainstormCanvas struct {
	ID        int        `db:"id"`
	CanvasUID uuid.UUID  `db:"canvas_uid"`
	ProjectID int        `db:"project_id"`
	StateJSON string     `db:"state_json"` // JSONB stored as string
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	CreatedBy *uuid.UUID `db:"created_by"`
	UpdatedBy *uuid.UUID `db:"updated_by"`
}

type NoteFolder struct {
	ID             int        `db:"id"`
	FolderUID      uuid.UUID  `db:"folder_uid"`
	ProjectID      int        `db:"project_id"`
	ParentFolderID *int       `db:"parent_folder_id"`
	Name           string     `db:"name"`
	Position       *int       `db:"position"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at"`
	CreatedBy      *uuid.UUID `db:"created_by"`
	UpdatedBy      *uuid.UUID `db:"updated_by"`
	IsActive       bool       `db:"is_active"`
}

type Note struct {
	ID          int        `db:"id"`
	NoteUID     uuid.UUID  `db:"note_uid"`
	ProjectID   int        `db:"project_id"`
	FolderID    *int       `db:"folder_id"`
	Title       string     `db:"title"`
	ContentJSON string     `db:"content_json"` // JSONB stored as string for rich content
	Position    *int       `db:"position"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	CreatedBy   *uuid.UUID `db:"created_by"`
	UpdatedBy   *uuid.UUID `db:"updated_by"`
	IsActive    bool       `db:"is_active"`
}

// Chat Conversation models for DevSprint-AI
type ChatConversation struct {
	ID              int        `db:"id"`
	ConversationUID uuid.UUID  `db:"conversation_uid"`
	ProjectID       int        `db:"project_id"`
	Name            string     `db:"name"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at"`
	CreatedBy       *uuid.UUID `db:"created_by"`
	UpdatedBy       *uuid.UUID `db:"updated_by"`
	IsActive        bool       `db:"is_active"`
}

type ChatMessage struct {
	ID             int        `db:"id"`
	MessageUID     uuid.UUID  `db:"message_uid"`
	ConversationID int        `db:"conversation_id"`
	MessageType    string     `db:"message_type"` // 'user' or 'ai'
	Content        string     `db:"content"`
	CreatedAt      time.Time  `db:"created_at"`
	CreatedBy      *uuid.UUID `db:"created_by"`
}
