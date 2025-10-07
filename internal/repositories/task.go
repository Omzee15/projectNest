package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"lucid-lists-backend/internal/models"
)

type taskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) GetByListID(ctx context.Context, listID int) ([]models.Task, error) {
	query := `
		SELECT id, task_uid, list_id, title, description, priority, status, color, position, is_completed,
			   due_date, completed_at, created_at, created_by, updated_at, updated_by, is_active
		FROM task
		WHERE list_id = $1 AND is_active = true
		ORDER BY COALESCE(position, 999999), created_at`

	rows, err := r.db.Query(ctx, query, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		err := rows.Scan(
			&t.ID, &t.TaskUID, &t.ListID, &t.Title, &t.Description, &t.Priority, &t.Status,
			&t.Color, &t.Position, &t.IsCompleted, &t.DueDate, &t.CompletedAt,
			&t.CreatedAt, &t.CreatedBy, &t.UpdatedAt, &t.UpdatedBy, &t.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *taskRepository) GetByUID(ctx context.Context, uid uuid.UUID) (*models.Task, error) {
	query := `
		SELECT id, task_uid, list_id, title, description, priority, status, color, position, is_completed,
			   due_date, completed_at, created_at, created_by, updated_at, updated_by, is_active
		FROM task
		WHERE task_uid = $1 AND is_active = true`

	var t models.Task
	err := r.db.QueryRow(ctx, query, uid).Scan(
		&t.ID, &t.TaskUID, &t.ListID, &t.Title, &t.Description, &t.Priority, &t.Status, &t.Color, &t.Position, &t.IsCompleted,
		&t.DueDate, &t.CompletedAt, &t.CreatedAt, &t.CreatedBy, &t.UpdatedAt, &t.UpdatedBy, &t.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &t, nil
}

func (r *taskRepository) Create(ctx context.Context, task *models.Task) error {
	query := `
		INSERT INTO task (task_uid, list_id, title, description, priority, status, color, position, is_completed, due_date, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query,
		task.TaskUID, task.ListID, task.Title, task.Description, task.Priority, task.Status, task.Color, task.Position, task.IsCompleted, task.DueDate, task.CreatedBy,
	).Scan(&task.ID, &task.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

func (r *taskRepository) Update(ctx context.Context, uid uuid.UUID, task *models.Task) error {
	now := time.Now()

	// Handle completion logic based on is_completed field
	var completedAt *time.Time
	if task.IsCompleted && task.CompletedAt == nil {
		completedAt = &now
	} else if !task.IsCompleted {
		completedAt = nil
	} else {
		completedAt = task.CompletedAt
	}

	query := `
		UPDATE task 
		SET title = $2, description = $3, priority = $4, status = $5, color = $6, position = $7, is_completed = $8,
			due_date = $9, completed_at = $10, updated_at = $11, updated_by = $12
		WHERE task_uid = $1 AND is_active = true`

	result, err := r.db.Exec(ctx, query,
		uid, task.Title, task.Description, task.Priority, task.Status, task.Color, task.Position, task.IsCompleted,
		task.DueDate, completedAt, now, task.UpdatedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (r *taskRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	query := `UPDATE task SET is_active = false WHERE task_uid = $1 AND is_active = true`

	result, err := r.db.Exec(ctx, query, uid)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (r *taskRepository) MoveToList(ctx context.Context, uid uuid.UUID, newListID int) error {
	query := `
		UPDATE task 
		SET list_id = $2, updated_at = $3
		WHERE task_uid = $1 AND is_active = true`

	now := time.Now()
	result, err := r.db.Exec(ctx, query, uid, newListID, now)

	if err != nil {
		return fmt.Errorf("failed to move task: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (r *taskRepository) GetByProjectID(ctx context.Context, projectID int) ([]models.Task, error) {
	query := `
		SELECT t.id, t.task_uid, t.list_id, t.title, t.description, t.priority, t.status, t.color, t.position, t.is_completed,
			   t.due_date, t.completed_at, t.created_at, t.created_by, t.updated_at, t.updated_by, t.is_active
		FROM task t
		INNER JOIN list l ON t.list_id = l.id
		WHERE l.project_id = $1 AND t.is_active = true AND l.is_active = true
		ORDER BY l.position, COALESCE(t.position, 999999), t.created_at`

	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks by project: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		err := rows.Scan(
			&t.ID, &t.TaskUID, &t.ListID, &t.Title, &t.Description, &t.Priority, &t.Status, &t.Color, &t.Position, &t.IsCompleted,
			&t.DueDate, &t.CompletedAt, &t.CreatedAt, &t.CreatedBy, &t.UpdatedAt, &t.UpdatedBy, &t.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *taskRepository) PartialUpdate(ctx context.Context, uid uuid.UUID, updates models.TaskUpdateRequest) error {
	setParts := []string{}
	args := []interface{}{uid}
	argCount := 2

	if updates.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argCount))
		args = append(args, *updates.Title)
		argCount++
	}
	if updates.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argCount))
		args = append(args, *updates.Description)
		argCount++
	}
	if updates.Priority != nil {
		setParts = append(setParts, fmt.Sprintf("priority = $%d", argCount))
		args = append(args, *updates.Priority)
		argCount++
	}
	if updates.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argCount))
		args = append(args, *updates.Status)
		argCount++

		// Handle completion logic
		now := time.Now()
		if *updates.Status == "completed" {
			setParts = append(setParts, fmt.Sprintf("completed_at = $%d", argCount))
			args = append(args, now)
			argCount++
		} else {
			setParts = append(setParts, fmt.Sprintf("completed_at = NULL"))
		}
	}
	if updates.Color != nil {
		setParts = append(setParts, fmt.Sprintf("color = $%d", argCount))
		args = append(args, *updates.Color)
		argCount++
	}
	if updates.Position != nil {
		setParts = append(setParts, fmt.Sprintf("position = $%d", argCount))
		args = append(args, *updates.Position)
		argCount++
	}
	if updates.DueDate != nil {
		setParts = append(setParts, fmt.Sprintf("due_date = $%d", argCount))
		args = append(args, *updates.DueDate)
		argCount++
	}
	if updates.IsCompleted != nil {
		setParts = append(setParts, fmt.Sprintf("is_completed = $%d", argCount))
		args = append(args, *updates.IsCompleted)
		argCount++

		// Handle completion logic
		now := time.Now()
		if *updates.IsCompleted {
			setParts = append(setParts, fmt.Sprintf("completed_at = $%d", argCount))
			args = append(args, now)
			argCount++
		} else {
			setParts = append(setParts, "completed_at = NULL")
		}
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	now := time.Now()
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argCount))
	args = append(args, now)

	query := fmt.Sprintf(`
		UPDATE task 
		SET %s
		WHERE task_uid = $1 AND is_active = true`,
		strings.Join(setParts, ", "))

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (r *taskRepository) GetMaxPositionByList(ctx context.Context, listID int) (int, error) {
	query := `SELECT COALESCE(MAX(position), 0) FROM task WHERE list_id = $1 AND is_active = true`

	var maxPosition int
	err := r.db.QueryRow(ctx, query, listID).Scan(&maxPosition)
	if err != nil {
		return 0, fmt.Errorf("failed to get max position: %w", err)
	}

	return maxPosition, nil
}

// CountByProjectUID counts all active tasks for a project by project UID
func (r *taskRepository) CountByProjectUID(ctx context.Context, projectUID uuid.UUID) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM task t
		JOIN list l ON t.list_id = l.id
		JOIN project p ON l.project_id = p.id
		WHERE p.project_uid = $1 AND t.is_active = true AND l.is_active = true AND p.is_active = true`

	var count int64
	err := r.db.QueryRow(ctx, query, projectUID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count tasks by project: %w", err)
	}

	return count, nil
}

// CountCompletedByProjectUID counts completed tasks for a project by project UID
func (r *taskRepository) CountCompletedByProjectUID(ctx context.Context, projectUID uuid.UUID) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM task t
		JOIN list l ON t.list_id = l.id
		JOIN project p ON l.project_id = p.id
		WHERE p.project_uid = $1 AND t.is_completed = true AND t.is_active = true AND l.is_active = true AND p.is_active = true`

	var count int64
	err := r.db.QueryRow(ctx, query, projectUID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count completed tasks by project: %w", err)
	}

	return count, nil
}
