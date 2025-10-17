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

type projectRepository struct {
	db       *pgxpool.Pool
	userRepo UserRepository
}

func NewProjectRepository(db *pgxpool.Pool, userRepo UserRepository) ProjectRepository {
	return &projectRepository{
		db:       db,
		userRepo: userRepo,
	}
}

func (r *projectRepository) GetAll(ctx context.Context, userID int) ([]models.Project, error) {
	query := `
		SELECT p.id, p.project_uid, p.user_id, p.name, p.description, p.status, p.color, p.position, p.start_date, p.end_date,
			   p.is_private, p.dbml_content, p.dbml_layout_data, p.flowchart_content, p.created_at, p.created_by, p.updated_at, p.updated_by, p.is_active
		FROM project p
		INNER JOIN project_member pm ON p.id = pm.project_id
		INNER JOIN users u ON pm.user_id = u.id
		WHERE p.is_active = true AND u.id = $1
		GROUP BY p.id, p.project_uid, p.user_id, p.name, p.description, p.status, p.color, p.position, p.start_date, p.end_date,
			   p.is_private, p.dbml_content, p.dbml_layout_data, p.flowchart_content, p.created_at, p.created_by, p.updated_at, p.updated_by, p.is_active
		ORDER BY COALESCE(p.position, 999999), p.created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		err := rows.Scan(
			&p.ID, &p.ProjectUID, &p.UserID, &p.Name, &p.Description, &p.Status, &p.Color, &p.Position,
			&p.StartDate, &p.EndDate, &p.IsPrivate, &p.DbmlContent, &p.DbmlLayoutData, &p.FlowchartContent, &p.CreatedAt, &p.CreatedBy,
			&p.UpdatedAt, &p.UpdatedBy, &p.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (r *projectRepository) GetByUID(ctx context.Context, uid uuid.UUID) (*models.Project, error) {
	query := `
		SELECT id, project_uid, user_id, name, description, status, color, position, start_date, end_date,
			   is_private, dbml_content, dbml_layout_data, flowchart_content, created_at, created_by, updated_at, updated_by, is_active
		FROM project
		WHERE project_uid = $1 AND is_active = true`

	var p models.Project
	err := r.db.QueryRow(ctx, query, uid).Scan(
		&p.ID, &p.ProjectUID, &p.UserID, &p.Name, &p.Description, &p.Status, &p.Color, &p.Position,
		&p.StartDate, &p.EndDate, &p.IsPrivate, &p.DbmlContent, &p.DbmlLayoutData, &p.FlowchartContent, &p.CreatedAt, &p.CreatedBy,
		&p.UpdatedAt, &p.UpdatedBy, &p.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found")
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return &p, nil
}

func (r *projectRepository) GetByID(ctx context.Context, id int) (*models.Project, error) {
	query := `
		SELECT id, project_uid, user_id, name, description, status, color, position, start_date, end_date,
			   is_private, dbml_content, dbml_layout_data, flowchart_content, created_at, created_by, updated_at, updated_by, is_active
		FROM project
		WHERE id = $1 AND is_active = true`

	var p models.Project
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.ProjectUID, &p.UserID, &p.Name, &p.Description, &p.Status, &p.Color, &p.Position,
		&p.StartDate, &p.EndDate, &p.IsPrivate, &p.DbmlContent, &p.DbmlLayoutData, &p.FlowchartContent, &p.CreatedAt, &p.CreatedBy,
		&p.UpdatedAt, &p.UpdatedBy, &p.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found")
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return &p, nil
}

func (r *projectRepository) GetWithLists(ctx context.Context, uid uuid.UUID) (*models.ProjectWithListsResponse, error) {
	// First get the project
	project, err := r.GetByUID(ctx, uid)
	if err != nil {
		return nil, err
	}

	// Get lists with tasks for this project
	query := `
		SELECT 
			l.id, l.list_uid, l.project_id, l.name, l.color, l.position,
			l.created_at, l.created_by, l.updated_at, l.updated_by, l.is_active,
			t.id, t.task_uid, t.list_id, t.title, t.description, t.priority, 
			t.status, t.color, t.position, t.is_completed, t.due_date, t.completed_at,
			t.created_at, t.created_by, t.updated_at, t.updated_by, t.is_active
		FROM list l
		LEFT JOIN task t ON l.id = t.list_id AND t.is_active = true
		WHERE l.project_id = $1 AND l.is_active = true
		ORDER BY l.position ASC, COALESCE(t.position, 999999) ASC, t.created_at ASC
	`

	rows, err := r.db.Query(ctx, query, project.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to query lists with tasks: %w", err)
	}
	defer rows.Close()

	listsMap := make(map[uuid.UUID]*models.ListWithTasksResponse)
	var listOrder []uuid.UUID

	for rows.Next() {
		var l models.List
		var t models.Task
		var taskID, taskListID *int
		var taskUID *uuid.UUID
		var taskTitle, taskStatus, taskColor, taskCreatedBy, taskUpdatedBy *string
		var taskPosition *int
		var taskIsCompleted, taskIsActive *bool
		var taskDueDate, taskCompletedAt, taskCreatedAt, taskUpdatedAt *time.Time

		err := rows.Scan(
			&l.ID, &l.ListUID, &l.ProjectID, &l.Name, &l.Color, &l.Position,
			&l.CreatedAt, &l.CreatedBy, &l.UpdatedAt, &l.UpdatedBy, &l.IsActive,
			&taskID, &taskUID, &taskListID, &taskTitle, &t.Description, &t.Priority,
			&taskStatus, &taskColor, &taskPosition, &taskIsCompleted, &taskDueDate, &taskCompletedAt,
			&taskCreatedAt, &taskCreatedBy, &taskUpdatedAt, &taskUpdatedBy, &taskIsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan list with task: %w", err)
		}

		// Initialize list if not already present
		if _, exists := listsMap[l.ListUID]; !exists {
			listsMap[l.ListUID] = &models.ListWithTasksResponse{
				ListResponse: models.ListResponse{
					ListUID:   l.ListUID,
					Name:      l.Name,
					Color:     l.Color,
					Position:  l.Position,
					CreatedAt: l.CreatedAt,
					UpdatedAt: l.UpdatedAt,
				},
				Tasks: []models.TaskResponse{},
			}
			listOrder = append(listOrder, l.ListUID)
		}

		// Add task if present
		if taskID != nil && taskUID != nil {
			task := models.TaskResponse{
				TaskUID:     *taskUID,
				Title:       safeStringDeref(taskTitle),
				Description: t.Description,
				Priority:    t.Priority,
				Status:      safeStringDeref(taskStatus),
				Color:       safeStringDeref(taskColor),
				Position:    taskPosition,
				IsCompleted: safeBoolDeref(taskIsCompleted),
				DueDate:     taskDueDate,
				CompletedAt: taskCompletedAt,
				CreatedAt:   safeTimeDeref(taskCreatedAt),
				UpdatedAt:   taskUpdatedAt,
			}
			listsMap[l.ListUID].Tasks = append(listsMap[l.ListUID].Tasks, task)
		}
	}

	// Convert map to slice in order
	var finalLists []models.ListWithTasksResponse
	for _, listUID := range listOrder {
		finalLists = append(finalLists, *listsMap[listUID])
	}

	response := &models.ProjectWithListsResponse{
		ProjectResponse: models.ProjectResponse{
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
		},
		Lists: finalLists,
	}

	return response, nil
}

func (r *projectRepository) Create(ctx context.Context, project *models.Project) error {
	query := `
		INSERT INTO project (project_uid, user_id, name, description, status, color, position, start_date, end_date, is_private, dbml_content, dbml_layout_data, flowchart_content, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query,
		project.ProjectUID, project.UserID, project.Name, project.Description, project.Status,
		project.Color, project.Position, project.StartDate, project.EndDate, project.IsPrivate, project.DbmlContent, project.DbmlLayoutData, project.FlowchartContent, project.CreatedBy,
	).Scan(&project.ID, &project.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	return nil
}

func (r *projectRepository) Update(ctx context.Context, uid uuid.UUID, project *models.Project) error {
	query := `
		UPDATE project 
		SET name = $2, description = $3, status = $4, color = $5, position = $6, start_date = $7, end_date = $8,
			is_private = $9, dbml_content = $10, dbml_layout_data = $11, flowchart_content = $12, updated_at = $13, updated_by = $14
		WHERE project_uid = $1 AND is_active = true`

	now := time.Now()
	result, err := r.db.Exec(ctx, query,
		uid, project.Name, project.Description, project.Status, project.Color, project.Position,
		project.StartDate, project.EndDate, project.IsPrivate, project.DbmlContent, project.DbmlLayoutData, project.FlowchartContent, now, project.UpdatedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("project not found")
	}

	return nil
}

func (r *projectRepository) PartialUpdate(ctx context.Context, uid uuid.UUID, updates models.ProjectUpdateRequest) error {
	setParts := []string{}
	args := []interface{}{uid}
	argCount := 2

	if updates.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argCount))
		args = append(args, *updates.Name)
		argCount++
	}
	if updates.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argCount))
		args = append(args, *updates.Description)
		argCount++
	}
	if updates.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argCount))
		args = append(args, *updates.Status)
		argCount++
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
	if updates.StartDate != nil {
		setParts = append(setParts, fmt.Sprintf("start_date = $%d", argCount))
		args = append(args, *updates.StartDate)
		argCount++
	}
	if updates.EndDate != nil {
		setParts = append(setParts, fmt.Sprintf("end_date = $%d", argCount))
		args = append(args, *updates.EndDate)
		argCount++
	}
	if updates.IsPrivate != nil {
		setParts = append(setParts, fmt.Sprintf("is_private = $%d", argCount))
		args = append(args, *updates.IsPrivate)
		argCount++
	}
	if updates.DbmlContent != nil {
		fmt.Printf("DEBUG: DbmlContent is not nil, value: '%s', length: %d\n", *updates.DbmlContent, len(*updates.DbmlContent))
		setParts = append(setParts, fmt.Sprintf("dbml_content = $%d", argCount))
		args = append(args, *updates.DbmlContent)
		argCount++
	} else {
		fmt.Printf("DEBUG: DbmlContent is nil\n")
	}
	if updates.DbmlLayoutData != nil {
		setParts = append(setParts, fmt.Sprintf("dbml_layout_data = $%d", argCount))
		args = append(args, *updates.DbmlLayoutData)
		argCount++
	}
	if updates.FlowchartContent != nil {
		setParts = append(setParts, fmt.Sprintf("flowchart_content = $%d", argCount))
		args = append(args, *updates.FlowchartContent)
		argCount++
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	now := time.Now()
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argCount))
	args = append(args, now)

	query := fmt.Sprintf(`
		UPDATE project 
		SET %s
		WHERE project_uid = $1 AND is_active = true`,
		strings.Join(setParts, ", "))

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("project not found")
	}

	return nil
}

func (r *projectRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	query := `UPDATE project SET is_active = false WHERE project_uid = $1 AND is_active = true`

	result, err := r.db.Exec(ctx, query, uid)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("project not found")
	}

	return nil
}

func (r *projectRepository) AddMember(ctx context.Context, projectID int, userUID uuid.UUID, role string) error {
	// Get user by UUID to get the integer ID
	user, err := r.userRepo.GetByUID(ctx, userUID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	query := `
		INSERT INTO project_member (project_id, user_id, role)
		VALUES ($1, $2, $3)`

	_, err = r.db.Exec(ctx, query, projectID, user.ID, role)
	if err != nil {
		return fmt.Errorf("failed to add project member: %w", err)
	}

	return nil
}

func (r *projectRepository) RemoveMember(ctx context.Context, projectID int, userUID uuid.UUID) error {
	// Get user by UUID to get the integer ID
	user, err := r.userRepo.GetByUID(ctx, userUID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	query := `
		DELETE FROM project_member 
		WHERE project_id = $1 AND user_id = $2`

	_, err = r.db.Exec(ctx, query, projectID, user.ID)
	if err != nil {
		return fmt.Errorf("failed to remove project member: %w", err)
	}

	return nil
}

// Helper functions
func safeStringDeref(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func safeBoolDeref(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}

func safeTimeDeref(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Time{}
}

func (r *projectRepository) GetMembers(ctx context.Context, projectID int) ([]models.ProjectMember, error) {
	query := `
		SELECT id, project_id, user_id, role, joined_at
		FROM project_member
		WHERE project_id = $1
		ORDER BY joined_at ASC`

	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query project members: %w", err)
	}
	defer rows.Close()

	var members []models.ProjectMember
	for rows.Next() {
		var m models.ProjectMember
		err := rows.Scan(&m.ID, &m.ProjectID, &m.UserID, &m.Role, &m.JoinedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project member: %w", err)
		}
		members = append(members, m)
	}

	return members, nil
}

func (r *projectRepository) IsMember(ctx context.Context, projectID int, userUID uuid.UUID) (bool, error) {
	// Get user by UUID to get the integer ID
	user, err := r.userRepo.GetByUID(ctx, userUID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return false, nil // User doesn't exist, so not a member
	}

	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM project_member 
			WHERE project_id = $1 AND user_id = $2
		)`

	var exists bool
	err = r.db.QueryRow(ctx, query, projectID, user.ID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check membership: %w", err)
	}

	return exists, nil
}
