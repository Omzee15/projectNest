package handlers

import (
	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/internal/services"
	"lucid-lists-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ListHandler struct {
	listService *services.ListService
}

func NewListHandler(listService *services.ListService) *ListHandler {
	return &ListHandler{
		listService: listService,
	}
}

// CreateList handles POST /api/lists
func (h *ListHandler) CreateList(c *gin.Context) {
	var req models.ListRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		logrus.WithError(err).Error("Failed to bind and validate list request")
		utils.SendError(c, err)
		return
	}

	list, err := h.listService.CreateList(c.Request.Context(), &req)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"project_uid": req.ProjectUID,
			"name":        req.Name,
			"position":    req.Position,
		}).Error("Failed to create list")
		utils.SendError(c, err)
		return
	}

	utils.CreatedResponse(c, list, "List created successfully")
}

// UpdateList handles PUT /api/lists/:uid
func (h *ListHandler) UpdateList(c *gin.Context) {
	uidStr := c.Param("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		utils.SendValidationError(c, "Invalid list ID format")
		return
	}

	var req models.ListRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.SendError(c, err)
		return
	}

	list, err := h.listService.UpdateList(c.Request.Context(), uid, &req)
	if err != nil {
		logrus.WithError(err).WithField("list_uid", uid).Error("Failed to update list")
		utils.SendError(c, err)
		return
	}

	utils.SuccessResponse(c, list, "List updated successfully")
}

// DeleteList handles DELETE /api/lists/:uid
func (h *ListHandler) DeleteList(c *gin.Context) {
	uidStr := c.Param("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		utils.SendValidationError(c, "Invalid list ID format")
		return
	}

	err = h.listService.DeleteList(c.Request.Context(), uid)
	if err != nil {
		logrus.WithError(err).WithField("list_uid", uid).Error("Failed to delete list")
		utils.SendError(c, err)
		return
	}

	utils.SuccessResponse(c, nil, "List deleted successfully")
}

// UpdatePosition handles PUT /api/lists/:uid/position
func (h *ListHandler) UpdatePosition(c *gin.Context) {
	uidStr := c.Param("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		utils.SendValidationError(c, "Invalid list ID format")
		return
	}

	var req models.UpdatePositionRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.SendError(c, err)
		return
	}

	list, err := h.listService.UpdatePosition(c.Request.Context(), uid, req.Position)
	if err != nil {
		logrus.WithError(err).WithField("list_uid", uid).Error("Failed to update list position")
		utils.SendError(c, err)
		return
	}

	utils.SuccessResponse(c, list, "List position updated successfully")
}

// PartialUpdateList handles PATCH /api/lists/:uid
func (h *ListHandler) PartialUpdateList(c *gin.Context) {
	uidStr := c.Param("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		utils.SendValidationError(c, "Invalid list ID format")
		return
	}

	var req models.ListUpdateRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.SendError(c, err)
		return
	}

	list, err := h.listService.PartialUpdateList(c.Request.Context(), uid, &req)
	if err != nil {
		logrus.WithError(err).WithField("list_uid", uid).Error("Failed to update list")
		utils.SendError(c, err)
		return
	}

	utils.SuccessResponse(c, list, "List updated successfully")
}
