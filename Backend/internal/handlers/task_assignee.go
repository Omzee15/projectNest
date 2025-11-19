package handlers

import (
	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/internal/repositories"
	"lucid-lists-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type TaskAssigneeHandler struct {
	taskAssigneeRepo repositories.TaskAssigneeRepository
	taskRepo         repositories.TaskRepository
	userRepo         repositories.UserRepository
}

func NewTaskAssigneeHandler(
	taskAssigneeRepo repositories.TaskAssigneeRepository,
	taskRepo repositories.TaskRepository,
	userRepo repositories.UserRepository,
) *TaskAssigneeHandler {
	return &TaskAssigneeHandler{
		taskAssigneeRepo: taskAssigneeRepo,
		taskRepo:         taskRepo,
		userRepo:         userRepo,
	}
}

// GetTaskAssignees handles GET /api/tasks/:uid/assignees
func (h *TaskAssigneeHandler) GetTaskAssignees(c *gin.Context) {
	taskUIDStr := c.Param("uid")
	taskUID, err := uuid.Parse(taskUIDStr)
	if err != nil {
		utils.SendValidationError(c, "Invalid task ID format")
		return
	}

	assignees, err := h.taskAssigneeRepo.GetByTaskUID(c.Request.Context(), taskUID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get task assignees")
		utils.ErrorResponse(c, 500, "Failed to retrieve task assignees")
		return
	}

	utils.SuccessResponse(c, assignees, "Task assignees retrieved successfully")
}

// AssignUserToTask handles POST /api/tasks/:uid/assignees
func (h *TaskAssigneeHandler) AssignUserToTask(c *gin.Context) {
	taskUIDStr := c.Param("uid")
	taskUID, err := uuid.Parse(taskUIDStr)
	if err != nil {
		utils.SendValidationError(c, "Invalid task ID format")
		return
	}

	var req models.AssignTaskRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.SendError(c, err)
		return
	}

	// Get task by UID to get internal ID
	task, err := h.taskRepo.GetByUID(c.Request.Context(), taskUID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get task")
		utils.ErrorResponse(c, 404, "Task not found")
		return
	}

	// Get user by UID to get internal ID
	user, err := h.userRepo.GetByUID(c.Request.Context(), req.UserUID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user")
		utils.ErrorResponse(c, 404, "User not found")
		return
	}

	// Assign user to task
	if err := h.taskAssigneeRepo.AssignUser(c.Request.Context(), task.ID, user.ID); err != nil {
		logrus.WithError(err).Error("Failed to assign user to task")
		utils.ErrorResponse(c, 500, "Failed to assign user to task")
		return
	}

	utils.SuccessResponse(c, nil, "User assigned to task successfully")
}

// UnassignUserFromTask handles DELETE /api/tasks/:uid/assignees/:userUid
func (h *TaskAssigneeHandler) UnassignUserFromTask(c *gin.Context) {
	taskUIDStr := c.Param("uid")
	taskUID, err := uuid.Parse(taskUIDStr)
	if err != nil {
		utils.SendValidationError(c, "Invalid task ID format")
		return
	}

	userUIDStr := c.Param("userUid")
	userUID, err := uuid.Parse(userUIDStr)
	if err != nil {
		utils.SendValidationError(c, "Invalid user ID format")
		return
	}

	// Get task by UID to get internal ID
	task, err := h.taskRepo.GetByUID(c.Request.Context(), taskUID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get task")
		utils.ErrorResponse(c, 404, "Task not found")
		return
	}

	// Get user by UID to get internal ID
	user, err := h.userRepo.GetByUID(c.Request.Context(), userUID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user")
		utils.ErrorResponse(c, 404, "User not found")
		return
	}

	// Unassign user from task
	if err := h.taskAssigneeRepo.UnassignUser(c.Request.Context(), task.ID, user.ID); err != nil {
		logrus.WithError(err).Error("Failed to unassign user from task")
		utils.ErrorResponse(c, 500, "Failed to unassign user from task")
		return
	}

	utils.SuccessResponse(c, nil, "User unassigned from task successfully")
}

// BulkAssignUsersToTask handles POST /api/tasks/:uid/assignees/bulk
func (h *TaskAssigneeHandler) BulkAssignUsersToTask(c *gin.Context) {
	taskUIDStr := c.Param("uid")
	taskUID, err := uuid.Parse(taskUIDStr)
	if err != nil {
		utils.SendValidationError(c, "Invalid task ID format")
		return
	}

	var req models.BulkAssignTaskRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.SendError(c, err)
		return
	}

	// Get task by UID to get internal ID
	task, err := h.taskRepo.GetByUID(c.Request.Context(), taskUID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get task")
		utils.ErrorResponse(c, 404, "Task not found")
		return
	}

	// Convert user UIDs to internal IDs
	userIDs := make([]int, 0, len(req.UserUIDs))
	for _, userUID := range req.UserUIDs {
		user, err := h.userRepo.GetByUID(c.Request.Context(), userUID)
		if err != nil {
			logrus.WithError(err).WithField("user_uid", userUID).Warn("User not found, skipping")
			continue
		}
		userIDs = append(userIDs, user.ID)
	}

	// Bulk assign users
	if err := h.taskAssigneeRepo.BulkAssign(c.Request.Context(), task.ID, userIDs); err != nil {
		logrus.WithError(err).Error("Failed to bulk assign users to task")
		utils.ErrorResponse(c, 500, "Failed to assign users to task")
		return
	}

	utils.SuccessResponse(c, nil, "Users assigned to task successfully")
}
