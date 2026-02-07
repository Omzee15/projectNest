package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/internal/services"
	"lucid-lists-backend/internal/utils"
)

type TaskCommentHandler struct {
	service *services.TaskCommentService
}

func NewTaskCommentHandler(service *services.TaskCommentService) *TaskCommentHandler {
	return &TaskCommentHandler{service: service}
}

// CreateComment creates a new comment on a task
// @Summary Create task comment
// @Tags Task Comments
// @Accept json
// @Produce json
// @Param comment body models.TaskCommentRequest true "Comment request"
// @Success 201 {object} models.TaskCommentResponse
// @Router /api/tasks/comments [post]
func (h *TaskCommentHandler) CreateComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.TaskCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	comment, err := h.service.CreateComment(c.Request.Context(), &req, userID.(int))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreatedResponse(c, comment, "Comment created successfully")
}

// GetCommentsByTaskUID retrieves all comments for a task
// @Summary Get task comments
// @Tags Task Comments
// @Produce json
// @Param uid path string true "Task UID"
// @Success 200 {array} models.TaskCommentResponse
// @Router /api/tasks/{uid}/comments [get]
func (h *TaskCommentHandler) GetCommentsByTaskUID(c *gin.Context) {
	taskUIDStr := c.Param("uid")
	taskUID, err := uuid.Parse(taskUIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid task UID")
		return
	}

	comments, err := h.service.GetCommentsByTaskUID(c.Request.Context(), taskUID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, comments, "Comments retrieved successfully")
}

// UpdateComment updates a comment
// @Summary Update task comment
// @Tags Task Comments
// @Accept json
// @Produce json
// @Param comment_uid path string true "Comment UID"
// @Param comment body models.TaskCommentUpdateRequest true "Comment update request"
// @Success 200 {object} models.TaskCommentResponse
// @Router /api/tasks/comments/{comment_uid} [put]
func (h *TaskCommentHandler) UpdateComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	commentUIDStr := c.Param("comment_uid")
	commentUID, err := uuid.Parse(commentUIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid comment UID")
		return
	}

	var req models.TaskCommentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	comment, err := h.service.UpdateComment(c.Request.Context(), commentUID, &req, userID.(int))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, comment, "Comment updated successfully")
}

// DeleteComment deletes a comment
// @Summary Delete task comment
// @Tags Task Comments
// @Produce json
// @Param comment_uid path string true "Comment UID"
// @Success 200 {object} map[string]string
// @Router /api/tasks/comments/{comment_uid} [delete]
func (h *TaskCommentHandler) DeleteComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	commentUIDStr := c.Param("comment_uid")
	commentUID, err := uuid.Parse(commentUIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid comment UID")
		return
	}

	if err := h.service.DeleteComment(c.Request.Context(), commentUID, userID.(int)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Comment deleted successfully")
}
