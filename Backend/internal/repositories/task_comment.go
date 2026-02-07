package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"lucid-lists-backend/internal/models"
)

type TaskCommentRepository struct {
	db *gorm.DB
}

func NewTaskCommentRepository(db *gorm.DB) *TaskCommentRepository {
	return &TaskCommentRepository{db: db}
}

// Create creates a new task comment
func (r *TaskCommentRepository) Create(ctx context.Context, comment *models.TaskComment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

// GetByUID retrieves a comment by its UID
func (r *TaskCommentRepository) GetByUID(ctx context.Context, commentUID uuid.UUID) (*models.TaskComment, error) {
	var comment models.TaskComment
	err := r.db.WithContext(ctx).Where("comment_uid = ? AND is_active = ?", commentUID, true).First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetByTaskID retrieves all comments for a task
func (r *TaskCommentRepository) GetByTaskID(ctx context.Context, taskID int) ([]models.TaskComment, error) {
	var comments []models.TaskComment
	err := r.db.WithContext(ctx).
		Where("task_id = ? AND is_active = ?", taskID, true).
		Order("created_at DESC").
		Find(&comments).Error
	return comments, err
}

// Update updates a comment
func (r *TaskCommentRepository) Update(ctx context.Context, comment *models.TaskComment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

// Delete soft deletes a comment
func (r *TaskCommentRepository) Delete(ctx context.Context, commentUID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.TaskComment{}).
		Where("comment_uid = ?", commentUID).
		Update("is_active", false).Error
}

// GetByTaskUID retrieves all comments for a task by task UID
func (r *TaskCommentRepository) GetByTaskUID(ctx context.Context, taskUID uuid.UUID) ([]models.TaskComment, error) {
	var comments []models.TaskComment
	err := r.db.WithContext(ctx).
		Joins("JOIN task ON task.id = task_comment.task_id").
		Where("task.task_uid = ? AND task_comment.is_active = ?", taskUID, true).
		Order("task_comment.created_at DESC").
		Find(&comments).Error
	return comments, err
}
