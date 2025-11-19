package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"lucid-lists-backend/internal/models"
)

type taskAssigneeRepository struct {
	db *pgxpool.Pool
}

func NewTaskAssigneeRepository(db *pgxpool.Pool) TaskAssigneeRepository {
	return &taskAssigneeRepository{db: db}
}

// GetByTaskID returns all assignees for a given task by its internal ID
func (r *taskAssigneeRepository) GetByTaskID(ctx context.Context, taskID int) ([]models.TaskAssigneeResponse, error) {
	query := `
		SELECT u.user_uid, u.name, u.email, ta.assigned_at
		FROM task_assignee ta
		JOIN users u ON ta.user_id = u.id
		WHERE ta.task_id = $1
		ORDER BY ta.assigned_at`

	rows, err := r.db.Query(ctx, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to query task assignees: %w", err)
	}
	defer rows.Close()

	var assignees []models.TaskAssigneeResponse
	for rows.Next() {
		var a models.TaskAssigneeResponse
		if err := rows.Scan(&a.UserUID, &a.Name, &a.Email, &a.AssignedAt); err != nil {
			return nil, fmt.Errorf("failed to scan assignee: %w", err)
		}
		assignees = append(assignees, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating assignee rows: %w", err)
	}

	return assignees, nil
}

// GetByTaskUID returns all assignees for a given task by its UUID
func (r *taskAssigneeRepository) GetByTaskUID(ctx context.Context, taskUID uuid.UUID) ([]models.TaskAssigneeResponse, error) {
	query := `
		SELECT u.user_uid, u.name, u.email, ta.assigned_at
		FROM task_assignee ta
		JOIN task t ON ta.task_id = t.id
		JOIN users u ON ta.user_id = u.id
		WHERE t.task_uid = $1
		ORDER BY ta.assigned_at`

	rows, err := r.db.Query(ctx, query, taskUID)
	if err != nil {
		return nil, fmt.Errorf("failed to query task assignees: %w", err)
	}
	defer rows.Close()

	var assignees []models.TaskAssigneeResponse
	for rows.Next() {
		var a models.TaskAssigneeResponse
		if err := rows.Scan(&a.UserUID, &a.Name, &a.Email, &a.AssignedAt); err != nil {
			return nil, fmt.Errorf("failed to scan assignee: %w", err)
		}
		assignees = append(assignees, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating assignee rows: %w", err)
	}

	return assignees, nil
}

// AssignUser assigns a user to a task
func (r *taskAssigneeRepository) AssignUser(ctx context.Context, taskID int, userID int) error {
	query := `
		INSERT INTO task_assignee (task_id, user_id, assigned_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (task_id, user_id) DO NOTHING`

	_, err := r.db.Exec(ctx, query, taskID, userID)
	if err != nil {
		return fmt.Errorf("failed to assign user to task: %w", err)
	}

	return nil
}

// UnassignUser removes a user assignment from a task
func (r *taskAssigneeRepository) UnassignUser(ctx context.Context, taskID int, userID int) error {
	query := `DELETE FROM task_assignee WHERE task_id = $1 AND user_id = $2`

	_, err := r.db.Exec(ctx, query, taskID, userID)
	if err != nil {
		return fmt.Errorf("failed to unassign user from task: %w", err)
	}

	return nil
}

// BulkAssign replaces all assignees for a task with a new set of users
func (r *taskAssigneeRepository) BulkAssign(ctx context.Context, taskID int, userIDs []int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Remove all existing assignees
	deleteQuery := `DELETE FROM task_assignee WHERE task_id = $1`
	if _, err := tx.Exec(ctx, deleteQuery, taskID); err != nil {
		return fmt.Errorf("failed to remove existing assignees: %w", err)
	}

	// Add new assignees
	if len(userIDs) > 0 {
		insertQuery := `
			INSERT INTO task_assignee (task_id, user_id, assigned_at)
			VALUES ($1, $2, NOW())`

		for _, userID := range userIDs {
			if _, err := tx.Exec(ctx, insertQuery, taskID, userID); err != nil {
				return fmt.Errorf("failed to assign user %d: %w", userID, err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// RemoveAllAssignees removes all user assignments from a task
func (r *taskAssigneeRepository) RemoveAllAssignees(ctx context.Context, taskID int) error {
	query := `DELETE FROM task_assignee WHERE task_id = $1`

	_, err := r.db.Exec(ctx, query, taskID)
	if err != nil {
		return fmt.Errorf("failed to remove all assignees: %w", err)
	}

	return nil
}
