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

type listRepository struct {
	db *pgxpool.Pool
}

func NewListRepository(db *pgxpool.Pool) ListRepository {
	return &listRepository{db: db}
}

func (r *listRepository) GetByProjectID(ctx context.Context, projectID int) ([]models.List, error) {
	query := `
		SELECT id, list_uid, project_id, name, color, position,
			   created_at, created_by, updated_at, updated_by, is_active
		FROM list
		WHERE project_id = $1 AND is_active = true
		ORDER BY position`

	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query lists: %w", err)
	}
	defer rows.Close()

	var lists []models.List
	for rows.Next() {
		var l models.List
		err := rows.Scan(
			&l.ID, &l.ListUID, &l.ProjectID, &l.Name, &l.Color, &l.Position,
			&l.CreatedAt, &l.CreatedBy, &l.UpdatedAt, &l.UpdatedBy, &l.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan list: %w", err)
		}
		lists = append(lists, l)
	}

	return lists, nil
}

func (r *listRepository) GetByUID(ctx context.Context, uid uuid.UUID) (*models.List, error) {
	query := `
		SELECT id, list_uid, project_id, name, color, position,
			   created_at, created_by, updated_at, updated_by, is_active
		FROM list
		WHERE list_uid = $1 AND is_active = true`

	var l models.List
	err := r.db.QueryRow(ctx, query, uid).Scan(
		&l.ID, &l.ListUID, &l.ProjectID, &l.Name, &l.Color, &l.Position,
		&l.CreatedAt, &l.CreatedBy, &l.UpdatedAt, &l.UpdatedBy, &l.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("list not found")
		}
		return nil, fmt.Errorf("failed to get list: %w", err)
	}

	return &l, nil
}

func (r *listRepository) Create(ctx context.Context, list *models.List) error {
	query := `
		INSERT INTO list (list_uid, project_id, name, color, position, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query,
		list.ListUID, list.ProjectID, list.Name, list.Color, list.Position, list.CreatedBy,
	).Scan(&list.ID, &list.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create list: %w", err)
	}

	return nil
}

func (r *listRepository) Update(ctx context.Context, uid uuid.UUID, list *models.List) error {
	query := `
		UPDATE list 
		SET name = $2, color = $3, updated_at = $4, updated_by = $5
		WHERE list_uid = $1 AND is_active = true`

	now := time.Now()
	result, err := r.db.Exec(ctx, query,
		uid, list.Name, list.Color, now, list.UpdatedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to update list: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("list not found")
	}

	return nil
}

func (r *listRepository) PartialUpdate(ctx context.Context, uid uuid.UUID, updates models.ListUpdateRequest) error {
	setParts := []string{}
	args := []interface{}{uid}
	argCount := 2

	if updates.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argCount))
		args = append(args, *updates.Name)
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

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	now := time.Now()
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argCount))
	args = append(args, now)

	query := fmt.Sprintf(`
		UPDATE list 
		SET %s
		WHERE list_uid = $1 AND is_active = true`,
		strings.Join(setParts, ", "))

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update list: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("list not found")
	}

	return nil
}

func (r *listRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	query := `UPDATE list SET is_active = false WHERE list_uid = $1 AND is_active = true`

	result, err := r.db.Exec(ctx, query, uid)
	if err != nil {
		return fmt.Errorf("failed to delete list: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("list not found")
	}

	return nil
}

func (r *listRepository) UpdatePosition(ctx context.Context, uid uuid.UUID, position int) error {
	query := `
		UPDATE list 
		SET position = $2, updated_at = $3
		WHERE list_uid = $1 AND is_active = true`

	now := time.Now()
	result, err := r.db.Exec(ctx, query, uid, position, now)

	if err != nil {
		return fmt.Errorf("failed to update list position: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("list not found")
	}

	return nil
}

func (r *listRepository) GetMaxPositionByProject(ctx context.Context, projectID int) (int, error) {
	query := `
		SELECT COALESCE(MAX(position), 0)
		FROM list
		WHERE project_id = $1 AND is_active = true`

	var maxPosition int
	err := r.db.QueryRow(ctx, query, projectID).Scan(&maxPosition)
	if err != nil {
		return 0, fmt.Errorf("failed to get max position: %w", err)
	}

	return maxPosition, nil
}
