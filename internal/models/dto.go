package models

import (
	"time"

	"github.com/google/uuid"
)

// DTOs (Data Transfer Objects) for API requests and responses
// These are exposed to the frontend and use UUID fields

type ProjectRequest struct {
	Name             string     `json:"name" validate:"required,min=1,max=255"`
	Description      *string    `json:"description"`
	Status           string     `json:"status" validate:"oneof=active inactive completed"`
	Color            string     `json:"color" validate:"omitempty,len=7,startswith=#"`
	Position         *int       `json:"position" validate:"omitempty,min=0"`
	StartDate        *time.Time `json:"start_date"`
	EndDate          *time.Time `json:"end_date"`
	IsPrivate        *bool      `json:"is_private"`
	DbmlContent      *string    `json:"dbml_content"`
	DbmlLayoutData   *string    `json:"dbml_layout_data"`
	FlowchartContent *string    `json:"flowchart_content"`
}

type ProjectResponse struct {
	ProjectUID       uuid.UUID  `json:"project_uid"`
	Name             string     `json:"name"`
	Description      *string    `json:"description"`
	Status           string     `json:"status"`
	Color            string     `json:"color"`
	Position         *int       `json:"position"`
	StartDate        *time.Time `json:"start_date"`
	EndDate          *time.Time `json:"end_date"`
	IsPrivate        bool       `json:"is_private"`
	DbmlContent      *string    `json:"dbml_content"`
	DbmlLayoutData   *string    `json:"dbml_layout_data"`
	FlowchartContent *string    `json:"flowchart_content"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
}

type ProjectWithListsResponse struct {
	ProjectResponse
	Lists []ListWithTasksResponse `json:"lists"`
}

type ListRequest struct {
	ProjectUID uuid.UUID `json:"project_uid" validate:"required"`
	Name       string    `json:"name" validate:"required,min=1,max=255"`
	Color      string    `json:"color" validate:"omitempty,len=7,startswith=#"`
	Position   int       `json:"position" validate:"min=0"`
}

type ListResponse struct {
	ListUID   uuid.UUID  `json:"list_uid"`
	Name      string     `json:"name"`
	Color     string     `json:"color"`
	Position  int        `json:"position"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type ListWithTasksResponse struct {
	ListResponse
	Tasks []TaskResponse `json:"tasks"`
}

type TaskRequest struct {
	ListUID     uuid.UUID  `json:"list_uid" validate:"required"`
	Title       string     `json:"title" validate:"required,min=1,max=255"`
	Description *string    `json:"description"`
	Priority    *string    `json:"priority" validate:"omitempty,oneof=low medium high"`
	Status      string     `json:"status" validate:"oneof=todo in_progress completed"`
	Color       string     `json:"color" validate:"omitempty,len=7,startswith=#"`
	Position    *int       `json:"position" validate:"omitempty,min=0"`
	IsCompleted *bool      `json:"is_completed"`
	DueDate     *time.Time `json:"due_date"`
}

type TaskResponse struct {
	TaskUID     uuid.UUID  `json:"task_uid"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Priority    *string    `json:"priority"`
	Status      string     `json:"status"`
	Color       string     `json:"color"`
	Position    *int       `json:"position"`
	IsCompleted bool       `json:"is_completed"`
	DueDate     *time.Time `json:"due_date"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type MoveTaskRequest struct {
	ListUID  uuid.UUID `json:"list_uid" validate:"required"`
	Position *int      `json:"position"`
}

type UpdatePositionRequest struct {
	Position int `json:"position" validate:"min=0"`
}

// Update request models for editing existing entities
type ProjectUpdateRequest struct {
	Name             *string    `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description      *string    `json:"description,omitempty"`
	Status           *string    `json:"status,omitempty" validate:"omitempty,oneof=active inactive completed"`
	Color            *string    `json:"color,omitempty" validate:"omitempty,len=7,startswith=#"`
	Position         *int       `json:"position,omitempty" validate:"omitempty,min=0"`
	StartDate        *time.Time `json:"start_date,omitempty"`
	EndDate          *time.Time `json:"end_date,omitempty"`
	IsPrivate        *bool      `json:"is_private,omitempty"`
	DbmlContent      *string    `json:"dbml_content,omitempty"`
	DbmlLayoutData   *string    `json:"dbml_layout_data,omitempty"`
	FlowchartContent *string    `json:"flowchart_content,omitempty"`
}

type ListUpdateRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Color    *string `json:"color,omitempty" validate:"omitempty,len=7,startswith=#"`
	Position *int    `json:"position,omitempty" validate:"omitempty,min=0"`
}

type TaskUpdateRequest struct {
	Title       *string    `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string    `json:"description,omitempty"`
	Priority    *string    `json:"priority,omitempty" validate:"omitempty,oneof=low medium high"`
	Status      *string    `json:"status,omitempty" validate:"omitempty,oneof=todo in_progress completed"`
	Color       *string    `json:"color,omitempty" validate:"omitempty,len=7,startswith=#"`
	Position    *int       `json:"position,omitempty" validate:"omitempty,min=0"`
	IsCompleted *bool      `json:"is_completed,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

// Progress tracking DTOs for Phase 2
type ProjectProgressResponse struct {
	TotalTasks     int     `json:"total_tasks"`
	CompletedTasks int     `json:"completed_tasks"`
	TodoTasks      int     `json:"todo_tasks"`
	Progress       float64 `json:"progress"`
}

type ProjectWithProgressResponse struct {
	ProjectResponse
	TaskStats ProjectProgressResponse `json:"task_stats"`
}

// Phase 3: Brainstorming & Planning Layer DTOs

type CanvasRequest struct {
	StateJSON string `json:"state_json" validate:"required"`
}

type CanvasResponse struct {
	CanvasUID uuid.UUID  `json:"canvas_uid"`
	ProjectID int        `json:"project_id"`
	StateJSON string     `json:"state_json"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// NoteContent represents the structured JSON content of a note
type NoteContent struct {
	Blocks []NoteBlock `json:"blocks"`
}

// NoteBlock represents different types of content blocks in a note
type NoteBlock struct {
	ID       string              `json:"id"`
	Type     string              `json:"type"` // "text", "checklist", "heading"
	Content  string              `json:"content,omitempty"`
	Metadata *NoteBlockMetadata  `json:"metadata,omitempty"`
	Items    []NoteChecklistItem `json:"items,omitempty"`    // For checklist blocks
	Children []NoteBlock         `json:"children,omitempty"` // For nested blocks
}

// NoteBlockMetadata holds additional properties for blocks
type NoteBlockMetadata struct {
	Level     *int  `json:"level,omitempty"` // For heading level (1-6)
	Bold      *bool `json:"bold,omitempty"`
	Italic    *bool `json:"italic,omitempty"`
	Underline *bool `json:"underline,omitempty"`
}

// NoteChecklistItem represents an item in a checklist
type NoteChecklistItem struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

type NoteRequest struct {
	Title       string      `json:"title" validate:"required,min=1,max=255"`
	ContentJSON NoteContent `json:"content"`
	FolderID    *int        `json:"folder_id,omitempty"`
	Position    *int        `json:"position,omitempty" validate:"omitempty,min=0"`
}

type NoteUpdateRequest struct {
	Title             *string      `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	ContentJSON       *NoteContent `json:"content,omitempty"` // For API layer
	ContentJSONString *string      `json:"-"`                 // For repository layer (internal use)
	FolderID          *int         `json:"folder_id,omitempty"`
	Position          *int         `json:"position,omitempty" validate:"omitempty,min=0"`
}

type NoteResponse struct {
	NoteUID     uuid.UUID   `json:"note_uid"`
	ProjectID   int         `json:"project_id"`
	Title       string      `json:"title"`
	ContentJSON NoteContent `json:"content"`
	FolderID    *int        `json:"folder_id"`
	Position    *int        `json:"position"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   *time.Time  `json:"updated_at"`
}

// Bulk notes response for getting all notes in a project
type NotesResponse struct {
	Notes []NoteResponse `json:"notes"`
	Total int            `json:"total"`
}

// Folder DTOs
type NoteFolderRequest struct {
	Name           string `json:"name" validate:"required,min=1,max=255"`
	ParentFolderID *int   `json:"parent_folder_id,omitempty"`
	Position       *int   `json:"position,omitempty" validate:"omitempty,min=0"`
}

type NoteFolderUpdateRequest struct {
	Name           *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	ParentFolderID *int    `json:"parent_folder_id,omitempty"`
	Position       *int    `json:"position,omitempty" validate:"omitempty,min=0"`
}

type NoteFolderResponse struct {
	ID             int        `json:"id"`
	FolderUID      uuid.UUID  `json:"folder_uid"`
	ProjectID      int        `json:"project_id"`
	ParentFolderID *int       `json:"parent_folder_id"`
	Name           string     `json:"name"`
	Position       *int       `json:"position"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

type NoteFoldersResponse struct {
	Folders []NoteFolderResponse `json:"folders"`
	Total   int                  `json:"total"`
}

// Authentication DTOs

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type UserResponse struct {
	UserUID   uuid.UUID `json:"user_uid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Chat Conversation DTOs for DevSprint-AI
type ChatConversationRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type ChatConversationResponse struct {
	ConversationUID uuid.UUID  `json:"conversation_uid"`
	ProjectUID      uuid.UUID  `json:"project_uid"`
	Name            string     `json:"name"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

type ChatMessageRequest struct {
	ConversationUID uuid.UUID `json:"conversation_uid" validate:"required"`
	Content         string    `json:"content" validate:"required,min=1"`
	MessageType     string    `json:"message_type" validate:"required,oneof=user ai"`
}

type ChatMessageResponse struct {
	MessageUID      uuid.UUID `json:"message_uid"`
	ConversationUID uuid.UUID `json:"conversation_uid"`
	MessageType     string    `json:"message_type"`
	Content         string    `json:"content"`
	CreatedAt       time.Time `json:"created_at"`
}

type ChatConversationWithMessagesResponse struct {
	ChatConversationResponse
	Messages []ChatMessageResponse `json:"messages"`
}

// Project Member DTOs
type ProjectMemberResponse struct {
	UserUID  uuid.UUID `json:"user_uid"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

// User Settings DTOs
type UserSettingsRequest struct {
	Theme                *string `json:"theme" validate:"omitempty,min=1,max=100"`
	Language             *string `json:"language" validate:"omitempty,len=2"`
	Timezone             *string `json:"timezone" validate:"omitempty,min=1,max=100"`
	NotificationsEnabled *bool   `json:"notifications_enabled"`
	EmailNotifications   *bool   `json:"email_notifications"`
	SoundEnabled         *bool   `json:"sound_enabled"`
	CompactMode          *bool   `json:"compact_mode"`
	AutoSave             *bool   `json:"auto_save"`
	AutoSaveInterval     *int    `json:"auto_save_interval" validate:"omitempty,min=10,max=600"`
}

type UserSettingsResponse struct {
	SettingsUID          uuid.UUID `json:"settings_uid"`
	Theme                string    `json:"theme"`
	Language             string    `json:"language"`
	Timezone             string    `json:"timezone"`
	NotificationsEnabled bool      `json:"notifications_enabled"`
	EmailNotifications   bool      `json:"email_notifications"`
	SoundEnabled         bool      `json:"sound_enabled"`
	CompactMode          bool      `json:"compact_mode"`
	AutoSave             bool      `json:"auto_save"`
	AutoSaveInterval     int       `json:"auto_save_interval"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
