package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/internal/repositories"
)

type TaskCommentService struct {
	commentRepo *repositories.TaskCommentRepository
	taskRepo    repositories.TaskRepository
	userRepo    repositories.UserRepository
}

func NewTaskCommentService(
	commentRepo *repositories.TaskCommentRepository,
	taskRepo    repositories.TaskRepository,
	userRepo    repositories.UserRepository,
) *TaskCommentService {
	return &TaskCommentService{
		commentRepo: commentRepo,
		taskRepo:    taskRepo,
		userRepo:    userRepo,
	}
}

// CreateComment creates a new comment on a task
func (s *TaskCommentService) CreateComment(ctx context.Context, req *models.TaskCommentRequest, userID int) (*models.TaskCommentResponse, error) {
	// Get task by UID
	task, err := s.taskRepo.GetByUID(ctx, req.TaskUID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// Create comment
	comment := &models.TaskComment{
		TaskID:  task.ID,
		UserID:  userID,
		Content: req.Content,
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Get user info
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Build response
	response := &models.TaskCommentResponse{
		CommentUID: comment.CommentUID,
		TaskUID:    task.TaskUID,
		UserUID:    user.UserUID,
		UserName:   user.Name,
		Content:    comment.Content,
		CreatedAt:  comment.CreatedAt,
		UpdatedAt:  comment.UpdatedAt,
	}

	return response, nil
}

// GetCommentsByTaskUID retrieves all comments for a task
func (s *TaskCommentService) GetCommentsByTaskUID(ctx context.Context, taskUID uuid.UUID) ([]models.TaskCommentResponse, error) {
	// Get task
	task, err := s.taskRepo.GetByUID(ctx, taskUID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// Get comments
	comments, err := s.commentRepo.GetByTaskID(ctx, task.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	// Build responses
	responses := make([]models.TaskCommentResponse, 0, len(comments))
	for _, comment := range comments {
		user, err := s.userRepo.GetByID(ctx, comment.UserID)
		if err != nil {
			continue // Skip if user not found
		}

		responses = append(responses, models.TaskCommentResponse{
			CommentUID: comment.CommentUID,
			TaskUID:    task.TaskUID,
			UserUID:    user.UserUID,
			UserName:   user.Name,
			Content:    comment.Content,
			CreatedAt:  comment.CreatedAt,
			UpdatedAt:  comment.UpdatedAt,
		})
	}

	return responses, nil
}

// UpdateComment updates a comment
func (s *TaskCommentService) UpdateComment(ctx context.Context, commentUID uuid.UUID, req *models.TaskCommentUpdateRequest, userID int) (*models.TaskCommentResponse, error) {
	// Get comment
	comment, err := s.commentRepo.GetByUID(ctx, commentUID)
	if err != nil {
		return nil, fmt.Errorf("comment not found: %w", err)
	}

	// Check ownership
	if comment.UserID != userID {
		return nil, fmt.Errorf("unauthorized: you can only edit your own comments")
	}

	// Update comment
	comment.Content = req.Content

	if err := s.commentRepo.Update(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	// Get task and user info
	task, err := s.taskRepo.GetByID(ctx, comment.TaskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Build response
	response := &models.TaskCommentResponse{
		CommentUID: comment.CommentUID,
		TaskUID:    task.TaskUID,
		UserUID:    user.UserUID,
		UserName:   user.Name,
		Content:    comment.Content,
		CreatedAt:  comment.CreatedAt,
		UpdatedAt:  comment.UpdatedAt,
	}

	return response, nil
}

// DeleteComment deletes a comment
func (s *TaskCommentService) DeleteComment(ctx context.Context, commentUID uuid.UUID, userID int) error {
	// Get comment
	comment, err := s.commentRepo.GetByUID(ctx, commentUID)
	if err != nil {
		return fmt.Errorf("comment not found: %w", err)
	}

	// Check ownership
	if comment.UserID != userID {
		return fmt.Errorf("unauthorized: you can only delete your own comments")
	}

	// Delete comment
	if err := s.commentRepo.Delete(ctx, commentUID); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}
