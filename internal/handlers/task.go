package handlers

import (
	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/internal/services"
	"lucid-lists-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type TaskHandler struct {
	taskService *services.TaskService
}

func NewTaskHandler(taskService *services.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// CreateTask handles POST /api/tasks
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req models.TaskRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.SendError(c, err)
		return
	}

	task, err := h.taskService.CreateTask(c.Request.Context(), &req)
	if err != nil {
		logrus.WithError(err).Error("Failed to create task")
		utils.SendError(c, err)
		return
	}

	utils.CreatedResponse(c, task, "Task created successfully")
}

// UpdateTask handles PUT /api/tasks/:uid
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	uidStr := c.Param("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		utils.SendValidationError(c, "Invalid task ID format")
		return
	}

	var req models.TaskRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.SendError(c, err)
		return
	}

	task, err := h.taskService.UpdateTask(c.Request.Context(), uid, &req)
	if err != nil {
		logrus.WithError(err).WithField("task_uid", uid).Error("Failed to update task")
		utils.SendError(c, err)
		return
	}

	utils.SuccessResponse(c, task, "Task updated successfully")
}

// DeleteTask handles DELETE /api/tasks/:uid
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	uidStr := c.Param("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		utils.SendValidationError(c, "Invalid task ID format")
		return
	}

	err = h.taskService.DeleteTask(c.Request.Context(), uid)
	if err != nil {
		logrus.WithError(err).WithField("task_uid", uid).Error("Failed to delete task")
		utils.SendError(c, err)
		return
	}

	utils.SuccessResponse(c, nil, "Task deleted successfully")
}

// MoveTask handles POST /api/tasks/:uid/move
func (h *TaskHandler) MoveTask(c *gin.Context) {
	uidStr := c.Param("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		utils.SendValidationError(c, "Invalid task ID format")
		return
	}

	var req models.MoveTaskRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.SendError(c, err)
		return
	}

	task, err := h.taskService.MoveTask(c.Request.Context(), uid, &req)
	if err != nil {
		logrus.WithError(err).WithField("task_uid", uid).Error("Failed to move task")
		utils.SendError(c, err)
		return
	}

	utils.SuccessResponse(c, task, "Task moved successfully")
}

// PartialUpdateTask handles PATCH /api/tasks/:uid
func (h *TaskHandler) PartialUpdateTask(c *gin.Context) {
	uidStr := c.Param("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		utils.SendValidationError(c, "Invalid task ID format")
		return
	}

	var req models.TaskUpdateRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.SendError(c, err)
		return
	}

	task, err := h.taskService.PartialUpdateTask(c.Request.Context(), uid, &req)
	if err != nil {
		logrus.WithError(err).WithField("task_uid", uid).Error("Failed to update task")
		utils.SendError(c, err)
		return
	}

	utils.SuccessResponse(c, task, "Task updated successfully")
}
