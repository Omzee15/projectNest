package models

import (
	"time"

	"github.com/google/uuid"
)

// DTOs (Data Transfer Objects) for API requests and responses
// These are exposed to the frontend and use UUID fields

type ProjectRequest struct {
	Name        string     `json:"name" validate:"required,min=1,max=255"`
	Description *string    `json:"description"`
	Status      string     `json:"status" validate:"oneof=active inactive completed"`
	Color       string     `json:"color" validate:"omitempty,len=7,startswith=#"`
	Position    *int       `json:"position" validate:"omitempty,min=0"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

type ProjectResponse struct {
	ProjectUID  uuid.UUID  `json:"project_uid"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Status      string     `json:"status"`
	Color       string     `json:"color"`
	Position    *int       `json:"position"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
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
	Status      string     `json:"status" validate:"oneof=todo completed"`
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
	Name        *string    `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string    `json:"description,omitempty"`
	Status      *string    `json:"status,omitempty" validate:"omitempty,oneof=active inactive completed"`
	Color       *string    `json:"color,omitempty" validate:"omitempty,len=7,startswith=#"`
	Position    *int       `json:"position,omitempty" validate:"omitempty,min=0"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
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
	Status      *string    `json:"status,omitempty" validate:"omitempty,oneof=todo completed"`
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
