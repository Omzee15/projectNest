package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"lucid-lists-backend/internal/models"
)

type NoteFolderRepository struct {
	db *pgxpool.Pool
}

func NewNoteFolderRepository(db *pgxpool.Pool) NoteFolderRepositoryInterface {
	return &NoteFolderRepository{
		db: db,
	}
}

func (r *NoteFolderRepository) GetByProjectUID(ctx context.Context, projectUID uuid.UUID) ([]models.NoteFolder, error) {
	query := `
		SELECT nf.id, nf.folder_uid, nf.project_id, nf.parent_folder_id, nf.name, nf.position, 
		       nf.created_at, nf.updated_at, nf.created_by, nf.updated_by, nf.is_active
		FROM note_folder nf
		JOIN project p ON p.id = nf.project_id
		WHERE p.project_uid = $1 AND nf.is_active = true
		ORDER BY nf.position ASC, nf.created_at ASC
	`

	rows, err := r.db.Query(ctx, query, projectUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get folders: %w", err)
	}
	defer rows.Close()

	var folders []models.NoteFolder
	for rows.Next() {
		var folder models.NoteFolder
		err := rows.Scan(
			&folder.ID, &folder.FolderUID, &folder.ProjectID, &folder.ParentFolderID,
			&folder.Name, &folder.Position, &folder.CreatedAt, &folder.UpdatedAt,
			&folder.CreatedBy, &folder.UpdatedBy, &folder.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan folder: %w", err)
		}
		folders = append(folders, folder)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return folders, nil
}

func (r *NoteFolderRepository) Create(ctx context.Context, projectID int, folder models.NoteFolder) (*models.NoteFolder, error) {
	// Get the next position if not provided
	if folder.Position == nil {
		position, err := r.getMaxPosition(ctx, projectID, folder.ParentFolderID)
		if err != nil {
			return nil, fmt.Errorf("failed to get max position: %w", err)
		}
		folder.Position = &position
	}

	query := `
		INSERT INTO note_folder (project_id, parent_folder_id, name, position)
		VALUES ($1, $2, $3, $4)
		RETURNING id, folder_uid, created_at, updated_at, created_by, updated_by, is_active
	`

	var created models.NoteFolder
	err := r.db.QueryRow(ctx, query, projectID, folder.ParentFolderID, folder.Name, folder.Position).Scan(
		&created.ID, &created.FolderUID, &created.CreatedAt, &created.UpdatedAt,
		&created.CreatedBy, &created.UpdatedBy, &created.IsActive,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create folder: %w", err)
	}

	// Set the fields we know
	created.ProjectID = projectID
	created.ParentFolderID = folder.ParentFolderID
	created.Name = folder.Name
	created.Position = folder.Position

	return &created, nil
}

func (r *NoteFolderRepository) Update(ctx context.Context, folderUID uuid.UUID, updates models.NoteFolder) (*models.NoteFolder, error) {
	query := `
		UPDATE note_folder 
		SET name = $2, parent_folder_id = $3, position = $4, updated_at = NOW()
		WHERE folder_uid = $1 AND is_active = true
		RETURNING id, folder_uid, project_id, parent_folder_id, name, position, 
		          created_at, updated_at, created_by, updated_by, is_active
	`

	var updated models.NoteFolder
	err := r.db.QueryRow(ctx, query, folderUID, updates.Name, updates.ParentFolderID, updates.Position).Scan(
		&updated.ID, &updated.FolderUID, &updated.ProjectID, &updated.ParentFolderID,
		&updated.Name, &updated.Position, &updated.CreatedAt, &updated.UpdatedAt,
		&updated.CreatedBy, &updated.UpdatedBy, &updated.IsActive,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update folder: %w", err)
	}

	return &updated, nil
}

func (r *NoteFolderRepository) Delete(ctx context.Context, folderUID uuid.UUID) error {
	query := `
		UPDATE note_folder 
		SET is_active = false, updated_at = NOW()
		WHERE folder_uid = $1 AND is_active = true
	`

	result, err := r.db.Exec(ctx, query, folderUID)
	if err != nil {
		return fmt.Errorf("failed to delete folder: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("folder not found")
	}

	return nil
}

func (r *NoteFolderRepository) getMaxPosition(ctx context.Context, projectID int, parentFolderID *int) (int, error) {
	query := `
		SELECT COALESCE(MAX(position), 0) + 1
		FROM note_folder
		WHERE project_id = $1 AND parent_folder_id IS NOT DISTINCT FROM $2 AND is_active = true
	`

	var maxPosition int
	err := r.db.QueryRow(ctx, query, projectID, parentFolderID).Scan(&maxPosition)
	if err != nil {
		return 0, fmt.Errorf("failed to get max position: %w", err)
	}

	return maxPosition, nil
}