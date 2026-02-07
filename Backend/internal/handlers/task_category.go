package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/internal/services"
	"lucid-lists-backend/internal/utils"
)

type TaskCategoryHandler struct {
	service *services.TaskCategoryService
}

func NewTaskCategoryHandler(service *services.TaskCategoryService) *TaskCategoryHandler {
	return &TaskCategoryHandler{service: service}
}

// CreateCategory creates a new task category
// @Summary Create task category
// @Tags Task Categories
// @Accept json
// @Produce json
// @Param category body models.TaskCategoryRequest true "Category request"
// @Success 201 {object} models.TaskCategoryResponse
// @Router /api/categories [post]
func (h *TaskCategoryHandler) CreateCategory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.TaskCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	category, err := h.service.CreateCategory(c.Request.Context(), &req, userID.(int))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.CreatedResponse(c, category, "Category created successfully")
}

// GetCategoriesByProjectUID retrieves all categories for a project
// @Summary Get project categories
// @Tags Task Categories
// @Produce json
// @Param uid path string true "Project UID"
// @Success 200 {array} models.TaskCategoryResponse
// @Router /api/projects/{uid}/categories [get]
func (h *TaskCategoryHandler) GetCategoriesByProjectUID(c *gin.Context) {
	projectUIDStr := c.Param("uid")
	projectUID, err := uuid.Parse(projectUIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project UID")
		return
	}

	categories, err := h.service.GetCategoriesByProjectUID(c.Request.Context(), projectUID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, categories, "Categories retrieved successfully")
}

// UpdateCategory updates a category
// @Summary Update task category
// @Tags Task Categories
// @Accept json
// @Produce json
// @Param category_uid path string true "Category UID"
// @Param category body models.TaskCategoryUpdateRequest true "Category update request"
// @Success 200 {object} models.TaskCategoryResponse
// @Router /api/categories/{category_uid} [put]
func (h *TaskCategoryHandler) UpdateCategory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	categoryUIDStr := c.Param("category_uid")
	categoryUID, err := uuid.Parse(categoryUIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category UID")
		return
	}

	var req models.TaskCategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	category, err := h.service.UpdateCategory(c.Request.Context(), categoryUID, &req, userID.(int))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, category, "Category updated successfully")
}

// DeleteCategory deletes a category
// @Summary Delete task category
// @Tags Task Categories
// @Produce json
// @Param category_uid path string true "Category UID"
// @Success 200 {object} map[string]string
// @Router /api/categories/{category_uid} [delete]
func (h *TaskCategoryHandler) DeleteCategory(c *gin.Context) {
	categoryUIDStr := c.Param("category_uid")
	categoryUID, err := uuid.Parse(categoryUIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category UID")
		return
	}

	if err := h.service.DeleteCategory(c.Request.Context(), categoryUID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Category deleted successfully")
}

// AssignCategoriesToTask assigns categories to a task
// @Summary Assign categories to task
// @Tags Task Categories
// @Accept json
// @Produce json
// @Param task_uid path string true "Task UID"
// @Param categories body []string true "Category UIDs"
// @Success 200 {object} map[string]string
// @Router /api/tasks/{task_uid}/categories [post]
func (h *TaskCategoryHandler) AssignCategoriesToTask(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	taskUIDStr := c.Param("uid")
	taskUID, err := uuid.Parse(taskUIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid task UID")
		return
	}

	var categoryUIDs []uuid.UUID
	if err := c.ShouldBindJSON(&categoryUIDs); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if err := h.service.AssignCategoriesToTask(c.Request.Context(), taskUID, categoryUIDs, userID.(int)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Categories assigned successfully")
}

// GetCategoriesByTaskUID retrieves all categories for a task
// @Summary Get task categories
// @Tags Task Categories
// @Produce json
// @Param uid path string true "Task UID"
// @Success 200 {array} models.TaskCategoryResponse
// @Router /api/tasks/{uid}/categories [get]
func (h *TaskCategoryHandler) GetCategoriesByTaskUID(c *gin.Context) {
	taskUIDStr := c.Param("uid")
	taskUID, err := uuid.Parse(taskUIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid task UID")
		return
	}

	categories, err := h.service.GetCategoriesByTaskUID(c.Request.Context(), taskUID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, categories, "Categories retrieved successfully")
}

// RemoveCategoryFromTask removes a category from a task
// @Summary Remove category from task
// @Tags Task Categories
// @Produce json
// @Param uid path string true "Task UID"
// @Param category_uid path string true "Category UID"
// @Success 200 {object} map[string]string
// @Router /api/tasks/{uid}/categories/{category_uid} [delete]
func (h *TaskCategoryHandler) RemoveCategoryFromTask(c *gin.Context) {
	taskUIDStr := c.Param("uid")
	taskUID, err := uuid.Parse(taskUIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid task UID")
		return
	}

	categoryUIDStr := c.Param("category_uid")
	categoryUID, err := uuid.Parse(categoryUIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category UID")
		return
	}

	if err := h.service.RemoveCategoryFromTask(c.Request.Context(), taskUID, categoryUID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Category removed successfully")
}
