package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"lucid-lists-backend/internal/models"
)

type TaskCategoryRepository struct {
	db *gorm.DB
}

func NewTaskCategoryRepository(db *gorm.DB) *TaskCategoryRepository {
	return &TaskCategoryRepository{db: db}
}

// Create creates a new task category
func (r *TaskCategoryRepository) Create(ctx context.Context, category *models.TaskCategory) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// GetByUID retrieves a category by its UID
func (r *TaskCategoryRepository) GetByUID(ctx context.Context, categoryUID uuid.UUID) (*models.TaskCategory, error) {
	var category models.TaskCategory
	err := r.db.WithContext(ctx).Where("category_uid = ? AND is_active = ?", categoryUID, true).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetByProjectID retrieves all categories for a project
func (r *TaskCategoryRepository) GetByProjectID(ctx context.Context, projectID int) ([]models.TaskCategory, error) {
	var categories []models.TaskCategory
	err := r.db.WithContext(ctx).
		Where("project_id = ? AND is_active = ?", projectID, true).
		Order("name ASC").
		Find(&categories).Error
	return categories, err
}

// GetByProjectUID retrieves all categories for a project by project UID
func (r *TaskCategoryRepository) GetByProjectUID(ctx context.Context, projectUID uuid.UUID) ([]models.TaskCategory, error) {
	var categories []models.TaskCategory
	err := r.db.WithContext(ctx).
		Joins("JOIN project ON project.id = task_category.project_id").
		Where("project.project_uid = ? AND task_category.is_active = ?", projectUID, true).
		Order("task_category.name ASC").
		Find(&categories).Error
	return categories, err
}

// Update updates a category
func (r *TaskCategoryRepository) Update(ctx context.Context, category *models.TaskCategory) error {
	return r.db.WithContext(ctx).Save(category).Error
}

// Delete soft deletes a category
func (r *TaskCategoryRepository) Delete(ctx context.Context, categoryUID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.TaskCategory{}).
		Where("category_uid = ?", categoryUID).
		Update("is_active", false).Error
}

// AssignToTask assigns a category to a task
func (r *TaskCategoryRepository) AssignToTask(ctx context.Context, taskID, categoryID int, assignedBy *int) error {
	categoryMap := &models.TaskCategoryMap{
		TaskID:     taskID,
		CategoryID: categoryID,
		AssignedBy: assignedBy,
	}
	return r.db.WithContext(ctx).Create(categoryMap).Error
}

// RemoveFromTask removes a category from a task
func (r *TaskCategoryRepository) RemoveFromTask(ctx context.Context, taskID, categoryID int) error {
	return r.db.WithContext(ctx).
		Where("task_id = ? AND category_id = ?", taskID, categoryID).
		Delete(&models.TaskCategoryMap{}).Error
}

// GetCategoriesByTaskID retrieves all categories for a task
func (r *TaskCategoryRepository) GetCategoriesByTaskID(ctx context.Context, taskID int) ([]models.TaskCategory, error) {
	var categories []models.TaskCategory
	err := r.db.WithContext(ctx).
		Joins("JOIN task_category_map ON task_category_map.category_id = task_category.id").
		Where("task_category_map.task_id = ? AND task_category.is_active = ?", taskID, true).
		Find(&categories).Error
	return categories, err
}

// GetCategoriesByTaskUID retrieves all categories for a task by task UID
func (r *TaskCategoryRepository) GetCategoriesByTaskUID(ctx context.Context, taskUID uuid.UUID) ([]models.TaskCategory, error) {
	var categories []models.TaskCategory
	err := r.db.WithContext(ctx).
		Joins("JOIN task_category_map ON task_category_map.category_id = task_category.id").
		Joins("JOIN task ON task.id = task_category_map.task_id").
		Where("task.task_uid = ? AND task_category.is_active = ?", taskUID, true).
		Find(&categories).Error
	return categories, err
}

// ClearTaskCategories removes all categories from a task
func (r *TaskCategoryRepository) ClearTaskCategories(ctx context.Context, taskID int) error {
	return r.db.WithContext(ctx).
		Where("task_id = ?", taskID).
		Delete(&models.TaskCategoryMap{}).Error
}

// GetTasksByCategory retrieves all tasks for a category
func (r *TaskCategoryRepository) GetTasksByCategory(ctx context.Context, categoryID int) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.WithContext(ctx).
		Joins("JOIN task_category_map ON task_category_map.task_id = task.id").
		Where("task_category_map.category_id = ? AND task.is_active = ?", categoryID, true).
		Find(&tasks).Error
	return tasks, err
}
