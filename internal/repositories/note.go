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

type noteRepository struct {
	db *pgxpool.Pool
}

func NewNoteRepository(db *pgxpool.Pool) NoteRepository {
	return &noteRepository{db: db}
}

func (r *noteRepository) GetByProjectUID(ctx context.Context, projectUID uuid.UUID) ([]models.Note, error) {
	query := `
		SELECT n.id, n.note_uid, n.project_id, n.folder_id, n.title, n.content_json, n.position, n.created_at, n.updated_at, n.created_by, n.updated_by, n.is_active
		FROM note n
		INNER JOIN project p ON n.project_id = p.id
		WHERE p.project_uid = $1 AND p.is_active = true AND n.is_active = true
		ORDER BY n.position ASC, n.created_at ASC
	`

	rows, err := r.db.Query(ctx, query, projectUID)
	if err != nil {
		logger.WithComponent("note-repository").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Failed to get notes by project UID")
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		err := rows.Scan(
			&note.ID,
			&note.NoteUID,
			&note.ProjectID,
			&note.FolderID,
			&note.Title,
			&note.ContentJSON,
			&note.Position,
			&note.CreatedAt,
			&note.UpdatedAt,
			&note.CreatedBy,
			&note.UpdatedBy,
			&note.IsActive,
		)
		if err != nil {
			logger.WithComponent("note-repository").
				WithFields(map[string]interface{}{"error": err.Error()}).
				Error("Failed to scan note row")
			return nil, err
		}
		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		logger.WithComponent("note-repository").
			WithFields(map[string]interface{}{"error": err.Error()}).
			Error("Error iterating over note rows")
		return nil, err
	}

	logger.WithComponent("note-repository").
		WithFields(map[string]interface{}{
			"project_uid": projectUID.String(),
			"count":       len(notes),
		}).
		Info("Successfully retrieved notes")

	return notes, nil
}

func (r *noteRepository) GetByUID(ctx context.Context, uid uuid.UUID) (*models.Note, error) {
	query := `
		SELECT id, note_uid, project_id, folder_id, title, content_json, position, created_at, updated_at, created_by, updated_by, is_active
		FROM note 
		WHERE note_uid = $1 AND is_active = true
	`

	var note models.Note
	err := r.db.QueryRow(ctx, query, uid).Scan(
		&note.ID,
		&note.NoteUID,
		&note.ProjectID,
		&note.FolderID,
		&note.Title,
		&note.ContentJSON,
		&note.Position,
		&note.CreatedAt,
		&note.UpdatedAt,
		&note.CreatedBy,
		&note.UpdatedBy,
		&note.IsActive,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			logger.WithComponent("note-repository").
				WithFields(map[string]interface{}{"note_uid": uid.String()}).
				Info("Note not found")
			return nil, pgx.ErrNoRows
		}
		logger.WithComponent("note-repository").
			WithFields(map[string]interface{}{
				"note_uid": uid.String(),
				"error":    err.Error(),
			}).
			Error("Failed to get note by UID")
		return nil, err
	}

	logger.WithComponent("note-repository").
		WithFields(map[string]interface{}{"note_uid": uid.String()}).
		Info("Successfully retrieved note")

	return &note, nil
}

func (r *noteRepository) Create(ctx context.Context, note *models.Note) error {
	query := `
		INSERT INTO note (note_uid, project_id, folder_id, title, content_json, position, created_at, created_by, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	note.NoteUID = uuid.New()
	note.CreatedAt = time.Now().UTC()
	note.IsActive = true

	err := r.db.QueryRow(ctx, query,
		note.NoteUID,
		note.ProjectID,
		note.FolderID,
		note.Title,
		note.ContentJSON,
		note.Position,
		note.CreatedAt,
		note.CreatedBy,
		note.IsActive,
	).Scan(&note.ID)

	if err != nil {
		logger.WithComponent("note-repository").
			WithFields(map[string]interface{}{
				"project_id": note.ProjectID,
				"title":      note.Title,
				"error":      err.Error(),
			}).
			Error("Failed to create note")
		return err
	}

	logger.WithComponent("note-repository").
		WithFields(map[string]interface{}{
			"note_uid":   note.NoteUID.String(),
			"project_id": note.ProjectID,
			"title":      note.Title,
		}).
		Info("Successfully created note")

	return nil
}

func (r *noteRepository) Update(ctx context.Context, uid uuid.UUID, note *models.Note) error {
	query := `
		UPDATE note 
		SET title = $1, content_json = $2, position = $3, updated_at = $4, updated_by = $5
		WHERE note_uid = $6 AND is_active = true
	`

	updatedAt := time.Now().UTC()
	result, err := r.db.Exec(ctx, query,
		note.Title,
		note.ContentJSON,
		note.Position,
		updatedAt,
		note.UpdatedBy,
		uid,
	)

	if err != nil {
		logger.WithComponent("note-repository").
			WithFields(map[string]interface{}{
				"note_uid": uid.String(),
				"error":    err.Error(),
			}).
			Error("Failed to update note")
		return err
	}

	if result.RowsAffected() == 0 {
		logger.WithComponent("note-repository").
			WithFields(map[string]interface{}{"note_uid": uid.String()}).
			Warn("No note found to update")
		return fmt.Errorf("note not found: %s", uid.String())
	}

	logger.WithComponent("note-repository").
		WithFields(map[string]interface{}{"note_uid": uid.String()}).
		Info("Successfully updated note")

	return nil
}

func (r *noteRepository) PartialUpdate(ctx context.Context, uid uuid.UUID, updates models.NoteUpdateRequest) error {
	// Build dynamic query based on provided fields
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if updates.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *updates.Title)
		argIndex++
	}

	if updates.ContentJSONString != nil {
		setParts = append(setParts, fmt.Sprintf("content_json = $%d", argIndex))
		args = append(args, *updates.ContentJSONString)
		argIndex++
	}

	if updates.FolderID != nil {
		setParts = append(setParts, fmt.Sprintf("folder_id = $%d", argIndex))
		args = append(args, *updates.FolderID)
		argIndex++
	}

	if updates.Position != nil {
		setParts = append(setParts, fmt.Sprintf("position = $%d", argIndex))
		args = append(args, *updates.Position)
		argIndex++
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Add updated_at
	updatedAt := time.Now().UTC()
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, updatedAt)
	argIndex++

	// Add WHERE clause
	args = append(args, uid)

	// Build the final query
	setClause := ""
	for i, part := range setParts {
		if i > 0 {
			setClause += ", "
		}
		setClause += part
	}

	query := fmt.Sprintf("UPDATE note SET %s WHERE note_uid = $%d AND is_active = true", setClause, argIndex)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		logger.WithComponent("note-repository").
			WithFields(map[string]interface{}{
				"note_uid": uid.String(),
				"error":    err.Error(),
			}).
			Error("Failed to partially update note")
		return err
	}

	if result.RowsAffected() == 0 {
		logger.WithComponent("note-repository").
			WithFields(map[string]interface{}{"note_uid": uid.String()}).
			Warn("No note found to update")
		return fmt.Errorf("note not found: %s", uid.String())
	}

	logger.WithComponent("note-repository").
		WithFields(map[string]interface{}{"note_uid": uid.String()}).
		Info("Successfully partially updated note")

	return nil
}

func (r *noteRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	query := `
		UPDATE note 
		SET is_active = false, updated_at = $1
		WHERE note_uid = $2 AND is_active = true
	`

	updatedAt := time.Now().UTC()
	result, err := r.db.Exec(ctx, query, updatedAt, uid)
	if err != nil {
		logger.WithComponent("note-repository").
			WithFields(map[string]interface{}{
				"note_uid": uid.String(),
				"error":    err.Error(),
			}).
			Error("Failed to delete note")
		return err
	}

	if result.RowsAffected() == 0 {
		logger.WithComponent("note-repository").
			WithFields(map[string]interface{}{"note_uid": uid.String()}).
			Warn("No note found to delete")
		return fmt.Errorf("note not found: %s", uid.String())
	}

	logger.WithComponent("note-repository").
		WithFields(map[string]interface{}{"note_uid": uid.String()}).
		Info("Successfully deleted note")

	return nil
}

func (r *noteRepository) GetMaxPositionByProject(ctx context.Context, projectID int) (int, error) {
	query := `
		SELECT COALESCE(MAX(position), -1) + 1 
		FROM note 
		WHERE project_id = $1 AND is_active = true
	`

	var maxPosition int
	err := r.db.QueryRow(ctx, query, projectID).Scan(&maxPosition)
	if err != nil {
		logger.WithComponent("note-repository").
			WithFields(map[string]interface{}{
				"project_id": projectID,
				"error":      err.Error(),
			}).
			Error("Failed to get max position")
		return 0, err
	}

	return maxPosition, nil
}

// UpdateFolderID updates the folder_id of a note, handling both setting and unsetting (NULL)
func (r *noteRepository) UpdateFolderID(ctx context.Context, uid uuid.UUID, folderID *int) error {
	updatedAt := time.Now().UTC()

	query := `
		UPDATE note 
		SET folder_id = $1, updated_at = $2 
		WHERE note_uid = $3 AND is_active = true
	`

	result, err := r.db.Exec(ctx, query, folderID, updatedAt, uid)
	if err != nil {
		logger.WithComponent("note-repository").
			WithFields(map[string]interface{}{
				"note_uid":  uid.String(),
				"folder_id": folderID,
				"error":     err.Error(),
			}).
			Error("Failed to update note folder_id")
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		logger.WithComponent("note-repository").
			WithFields(map[string]interface{}{
				"note_uid": uid.String(),
			}).
			Error("Note not found for folder_id update")
		return fmt.Errorf("note not found")
	}

	logger.WithComponent("note-repository").
		WithFields(map[string]interface{}{
			"note_uid":  uid.String(),
			"folder_id": folderID,
		}).
		Info("Successfully updated note folder_id")

	return nil
}
