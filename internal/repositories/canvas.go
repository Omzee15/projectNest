package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/pkg/logger"
)

type canvasRepository struct {
	db *pgxpool.Pool
}

func NewCanvasRepository(db *pgxpool.Pool) CanvasRepository {
	return &canvasRepository{db: db}
}

func (r *canvasRepository) GetByProjectUID(ctx context.Context, projectUID uuid.UUID) (*models.BrainstormCanvas, error) {
	query := `
		SELECT bc.id, bc.canvas_uid, bc.project_id, bc.state_json, bc.created_at, bc.updated_at, bc.created_by, bc.updated_by
		FROM brainstorm_canvas bc
		INNER JOIN project p ON bc.project_id = p.id
		WHERE p.project_uid = $1 AND p.is_active = true
	`

	var canvas models.BrainstormCanvas
	err := r.db.QueryRow(ctx, query, projectUID).Scan(
		&canvas.ID,
		&canvas.CanvasUID,
		&canvas.ProjectID,
		&canvas.StateJSON,
		&canvas.CreatedAt,
		&canvas.UpdatedAt,
		&canvas.CreatedBy,
		&canvas.UpdatedBy,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			logger.WithComponent("canvas-repository").
				WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
				Info("No canvas found for project")
			return nil, pgx.ErrNoRows
		}
		logger.WithComponent("canvas-repository").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Failed to get canvas by project UID")
		return nil, err
	}

	logger.WithComponent("canvas-repository").
		WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
		Info("Successfully retrieved canvas")

	return &canvas, nil
}

func (r *canvasRepository) Create(ctx context.Context, canvas *models.BrainstormCanvas) error {
	query := `
		INSERT INTO brainstorm_canvas (canvas_uid, project_id, state_json, created_at, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	canvas.CanvasUID = uuid.New()
	canvas.CreatedAt = time.Now().UTC()

	err := r.db.QueryRow(ctx, query,
		canvas.CanvasUID,
		canvas.ProjectID,
		canvas.StateJSON,
		canvas.CreatedAt,
		canvas.CreatedBy,
	).Scan(&canvas.ID)

	if err != nil {
		logger.WithComponent("canvas-repository").
			WithFields(map[string]interface{}{
				"project_id": canvas.ProjectID,
				"error":      err.Error(),
			}).
			Error("Failed to create canvas")
		return err
	}

	logger.WithComponent("canvas-repository").
		WithFields(map[string]interface{}{
			"canvas_uid": canvas.CanvasUID.String(),
			"project_id": canvas.ProjectID,
		}).
		Info("Successfully created canvas")

	return nil
}

func (r *canvasRepository) Update(ctx context.Context, projectUID uuid.UUID, stateJSON string) error {
	query := `
		UPDATE brainstorm_canvas 
		SET state_json = $1, updated_at = $2
		FROM project p
		WHERE brainstorm_canvas.project_id = p.id 
		AND p.project_uid = $3 
		AND p.is_active = true
	`

	updatedAt := time.Now().UTC()
	result, err := r.db.Exec(ctx, query, stateJSON, updatedAt, projectUID)
	if err != nil {
		logger.WithComponent("canvas-repository").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Failed to update canvas")
		return err
	}

	if result.RowsAffected() == 0 {
		logger.WithComponent("canvas-repository").
			WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
			Warn("No canvas found to update")
		return fmt.Errorf("canvas not found for project %s", projectUID.String())
	}

	logger.WithComponent("canvas-repository").
		WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
		Info("Successfully updated canvas")

	return nil
}

func (r *canvasRepository) Delete(ctx context.Context, projectUID uuid.UUID) error {
	query := `
		DELETE FROM brainstorm_canvas 
		USING project p
		WHERE brainstorm_canvas.project_id = p.id 
		AND p.project_uid = $1 
		AND p.is_active = true
	`

	result, err := r.db.Exec(ctx, query, projectUID)
	if err != nil {
		logger.WithComponent("canvas-repository").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Failed to delete canvas")
		return err
	}

	if result.RowsAffected() == 0 {
		logger.WithComponent("canvas-repository").
			WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
			Warn("No canvas found to delete")
		return fmt.Errorf("canvas not found for project %s", projectUID.String())
	}

	logger.WithComponent("canvas-repository").
		WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
		Info("Successfully deleted canvas")

	return nil
}