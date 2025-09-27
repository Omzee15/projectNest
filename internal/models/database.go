package models

import (
	"time"

	"github.com/google/uuid"
)

// Database models - these map directly to the database schema
// These use integer IDs internally for efficiency

type Workspace struct {
	ID          int        `db:"id"`
	WorkspaceUID uuid.UUID `db:"workspace_uid"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	CreatedBy   *uuid.UUID `db:"created_by"`
	UpdatedAt   *time.Time `db:"updated_at"`
	UpdatedBy   *uuid.UUID `db:"updated_by"`
	IsActive    bool       `db:"is_active"`
}

type Project struct {
	ID          int        `db:"id"`
	ProjectUID  uuid.UUID  `db:"project_uid"`
	WorkspaceID int        `db:"workspace_id"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`
	Status      string     `db:"status"`
	Color       string     `db:"color"`
	Position    *int       `db:"position"`
	StartDate   *time.Time `db:"start_date"`
	EndDate     *time.Time `db:"end_date"`
	CreatedAt   time.Time  `db:"created_at"`
	CreatedBy   *uuid.UUID `db:"created_by"`
	UpdatedAt   *time.Time `db:"updated_at"`
	UpdatedBy   *uuid.UUID `db:"updated_by"`
	IsActive    bool       `db:"is_active"`
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
