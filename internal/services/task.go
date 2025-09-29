package services

import (
	"context"

	"github.com/google/uuid"

	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/internal/repositories"
	"lucid-lists-backend/internal/utils"
)

type TaskService struct {
	taskRepo repositories.TaskRepository
	listRepo repositories.ListRepository
}

func NewTaskService(taskRepo repositories.TaskRepository, listRepo repositories.ListRepository) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
		listRepo: listRepo,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, req *models.TaskRequest) (*models.TaskResponse, error) {
	// We need to resolve list_uid to list_id
	list, err := s.listRepo.GetByUID(ctx, req.ListUID)
	if err != nil {
		if err.Error() == "list not found" {
			return nil, utils.NewNotFoundError("List not found")
		}
		return nil, utils.NewInternalError("Failed to get list")
	}

	// Get next position if not specified
	position := req.Position
	if position == nil {
		maxPos, err := s.taskRepo.GetMaxPositionByList(ctx, list.ID)
		if err == nil {
			newPos := maxPos + 1
			position = &newPos
		}
	}

	// Set default color if not provided
	color := req.Color
	if color == "" {
		color = "#FFFFFF"
	}

	// Set default is_completed if not provided
	isCompleted := false
	if req.IsCompleted != nil {
		isCompleted = *req.IsCompleted
	}

	// Create task model
	task := &models.Task{
		TaskUID:     uuid.New(),
		ListID:      list.ID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
		Color:       color,
		Position:    position,
		IsCompleted: isCompleted,
		DueDate:     req.DueDate,
		IsActive:    true,
		CreatedBy:   nil, // No user authentication yet
	}

	if req.Status == "" {
		task.Status = "todo"
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, utils.NewInternalError("Failed to create task")
	}

	return &models.TaskResponse{
		TaskUID:     task.TaskUID,
		Title:       task.Title,
		Description: task.Description,
		Priority:    task.Priority,
		Status:      task.Status,
		Color:       task.Color,
		Position:    task.Position,
		IsCompleted: task.IsCompleted,
		DueDate:     task.DueDate,
		CompletedAt: task.CompletedAt,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, uid uuid.UUID, req *models.TaskRequest) (*models.TaskResponse, error) {
	// Check if task exists
	_, err := s.taskRepo.GetByUID(ctx, uid)
	if err != nil {
		if err.Error() == "task not found" {
			return nil, utils.NewNotFoundError("Task not found")
		}
		return nil, utils.NewInternalError("Failed to get task")
	}

	// Update task fields
	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
		Color:       req.Color,
		Position:    req.Position,
		DueDate:     req.DueDate,
	}

	if err := s.taskRepo.Update(ctx, uid, task); err != nil {
		if err.Error() == "task not found" {
			return nil, utils.NewNotFoundError("Task not found")
		}
		return nil, utils.NewInternalError("Failed to update task")
	}

	// Get updated task
	updatedTask, err := s.taskRepo.GetByUID(ctx, uid)
	if err != nil {
		return nil, utils.NewInternalError("Failed to get updated task")
	}

	return &models.TaskResponse{
		TaskUID:     updatedTask.TaskUID,
		Title:       updatedTask.Title,
		Description: updatedTask.Description,
		Priority:    updatedTask.Priority,
		Status:      updatedTask.Status,
		Color:       updatedTask.Color,
		Position:    updatedTask.Position,
		IsCompleted: updatedTask.IsCompleted,
		DueDate:     updatedTask.DueDate,
		CompletedAt: updatedTask.CompletedAt,
		CreatedAt:   updatedTask.CreatedAt,
		UpdatedAt:   updatedTask.UpdatedAt,
	}, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, uid uuid.UUID) error {
	if err := s.taskRepo.Delete(ctx, uid); err != nil {
		if err.Error() == "task not found" {
			return utils.NewNotFoundError("Task not found")
		}
		return utils.NewInternalError("Failed to delete task")
	}

	return nil
}

func (s *TaskService) MoveTask(ctx context.Context, uid uuid.UUID, req *models.MoveTaskRequest) (*models.TaskResponse, error) {
	// Get the target list to verify it exists and get its internal ID
	list, err := s.listRepo.GetByUID(ctx, req.ListUID)
	if err != nil {
		if err.Error() == "list not found" {
			return nil, utils.NewNotFoundError("List not found")
		}
		return nil, utils.NewInternalError("Failed to get list")
	}

	// Move the task
	if err := s.taskRepo.MoveToList(ctx, uid, list.ID); err != nil {
		if err.Error() == "task not found" {
			return nil, utils.NewNotFoundError("Task not found")
		}
		return nil, utils.NewInternalError("Failed to move task")
	}

	// Get updated task
	updatedTask, err := s.taskRepo.GetByUID(ctx, uid)
	if err != nil {
		return nil, utils.NewInternalError("Failed to get updated task")
	}

	return &models.TaskResponse{
		TaskUID:     updatedTask.TaskUID,
		Title:       updatedTask.Title,
		Description: updatedTask.Description,
		Priority:    updatedTask.Priority,
		Status:      updatedTask.Status,
		Color:       updatedTask.Color,
		Position:    updatedTask.Position,
		IsCompleted: updatedTask.IsCompleted,
		DueDate:     updatedTask.DueDate,
		CompletedAt: updatedTask.CompletedAt,
		CreatedAt:   updatedTask.CreatedAt,
		UpdatedAt:   updatedTask.UpdatedAt,
	}, nil
}

// PartialUpdateTask updates specific fields of a task
func (s *TaskService) PartialUpdateTask(ctx context.Context, uid uuid.UUID, updates *models.TaskUpdateRequest) (*models.TaskResponse, error) {
	// Use repository method for partial update
	if err := s.taskRepo.PartialUpdate(ctx, uid, *updates); err != nil {
		if err.Error() == "task not found" {
			return nil, utils.NewNotFoundError("Task not found")
		}
		return nil, utils.NewInternalError("Failed to update task")
	}

	// Get updated task
	updatedTask, err := s.taskRepo.GetByUID(ctx, uid)
	if err != nil {
		return nil, utils.NewInternalError("Failed to get updated task")
	}

	return &models.TaskResponse{
		TaskUID:     updatedTask.TaskUID,
		Title:       updatedTask.Title,
		Description: updatedTask.Description,
		Priority:    updatedTask.Priority,
		Status:      updatedTask.Status,
		Color:       updatedTask.Color,
		Position:    updatedTask.Position,
		IsCompleted: updatedTask.IsCompleted,
		DueDate:     updatedTask.DueDate,
		CompletedAt: updatedTask.CompletedAt,
		CreatedAt:   updatedTask.CreatedAt,
		UpdatedAt:   updatedTask.UpdatedAt,
	}, nil
}
