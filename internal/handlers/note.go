package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/internal/services"
	"lucid-lists-backend/internal/utils"
	"lucid-lists-backend/pkg/logger"
)

type NoteHandler struct {
	noteService *services.NoteService
}

func NewNoteHandler(noteService *services.NoteService) *NoteHandler {
	return &NoteHandler{
		noteService: noteService,
	}
}

// GetNotesByProject handles GET /api/projects/:uid/notes
func (h *NoteHandler) GetNotesByProject(c *gin.Context) {
	uidParam := c.Param("uid")

	logger.WithComponent("note-handler").
		WithFields(map[string]interface{}{"project_uid": uidParam}).
		Info("Getting notes for project")

	projectUID, err := uuid.Parse(uidParam)
	if err != nil {
		logger.WithComponent("note-handler").
			WithFields(map[string]interface{}{"invalid_uid": uidParam}).
			Warn("Invalid project UID format")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project UID format")
		return
	}

	notes, err := h.noteService.GetNotesByProject(c.Request.Context(), projectUID)
	if err != nil {
		logger.WithComponent("note-handler").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Failed to get notes")

		if err.Error() == "project not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Project not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve notes")
		}
		return
	}

	logger.WithComponent("note-handler").
		WithFields(map[string]interface{}{
			"project_uid": projectUID.String(),
			"count":       notes.Total,
		}).
		Info("Successfully retrieved notes")

	utils.SuccessResponse(c, notes, "")
}

// GetNote handles GET /api/notes/:uid
func (h *NoteHandler) GetNote(c *gin.Context) {
	uidParam := c.Param("uid")

	logger.WithComponent("note-handler").
		WithFields(map[string]interface{}{"note_uid": uidParam}).
		Info("Getting note")

	noteUID, err := uuid.Parse(uidParam)
	if err != nil {
		logger.WithComponent("note-handler").
			WithFields(map[string]interface{}{"invalid_uid": uidParam}).
			Warn("Invalid note UID format")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid note UID format")
		return
	}

	note, err := h.noteService.GetNote(c.Request.Context(), noteUID)
	if err != nil {
		logger.WithComponent("note-handler").
			WithFields(map[string]interface{}{
				"note_uid": noteUID.String(),
				"error":    err.Error(),
			}).
			Error("Failed to get note")
		utils.ErrorResponse(c, http.StatusNotFound, "Note not found")
		return
	}

	logger.WithComponent("note-handler").
		WithFields(map[string]interface{}{"note_uid": noteUID.String()}).
		Info("Successfully retrieved note")

	utils.SuccessResponse(c, note, "")
}

// CreateNote handles POST /api/projects/:uid/notes
func (h *NoteHandler) CreateNote(c *gin.Context) {
	uidParam := c.Param("uid")
	var req models.NoteRequest

	logger.WithComponent("note-handler").
		WithFields(map[string]interface{}{"project_uid": uidParam}).
		Info("Creating note for project")

	projectUID, err := uuid.Parse(uidParam)
	if err != nil {
		logger.WithComponent("note-handler").
			WithFields(map[string]interface{}{"invalid_uid": uidParam}).
			Warn("Invalid project UID format")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project UID format")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithComponent("note-handler").
			WithFields(map[string]interface{}{"error": err.Error()}).
			Warn("Invalid request body for create note")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	note, err := h.noteService.CreateNote(c.Request.Context(), projectUID, &req)
	if err != nil {
		logger.WithComponent("note-handler").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"title":       req.Title,
				"error":       err.Error(),
			}).
			Error("Failed to create note")

		if err.Error() == "project not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Project not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create note")
		}
		return
	}

	logger.WithComponent("note-handler").
		WithFields(map[string]interface{}{
			"note_uid": note.NoteUID.String(),
			"title":    note.Title,
		}).
		Info("Successfully created note")

	utils.SuccessResponse(c, note, "Note created successfully")
}

// UpdateNote handles PUT /api/notes/:uid
func (h *NoteHandler) UpdateNote(c *gin.Context) {
	uidParam := c.Param("uid")
	var req models.NoteRequest

	logger.WithComponent("note-handler").
		WithFields(map[string]interface{}{"note_uid": uidParam}).
		Info("Updating note")

	noteUID, err := uuid.Parse(uidParam)
	if err != nil {
		logger.WithComponent("note-handler").
			WithFields(map[string]interface{}{"invalid_uid": uidParam}).
			Warn("Invalid note UID format")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid note UID format")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithComponent("note-handler").
			WithFields(map[string]interface{}{"error": err.Error()}).
			Warn("Invalid request body for update note")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	note, err := h.noteService.UpdateNote(c.Request.Context(), noteUID, &req)
	if err != nil {
		logger.WithComponent("note-handler").
			WithFields(map[string]interface{}{
				"note_uid": noteUID.String(),
				"error":    err.Error(),
			}).
			Error("Failed to update note")
		utils.ErrorResponse(c, http.StatusNotFound, "Note not found")
		return
	}

	logger.WithComponent("note-handler").
		WithFields(map[string]interface{}{"note_uid": noteUID.String()}).
		Info("Successfully updated note")

	utils.SuccessResponse(c, note, "Note updated successfully")
}

// PartialUpdateNote handles PATCH /api/notes/:uid
func (h *NoteHandler) PartialUpdateNote(c *gin.Context) {
	uidParam := c.Param("uid")
	var req models.NoteUpdateRequest

	logger.WithComponent("note-handler").
		WithFields(map[string]interface{}{"note_uid": uidParam}).
		Info("Partially updating note")

	noteUID, err := uuid.Parse(uidParam)
	if err != nil {
		logger.WithComponent("note-handler").
			WithFields(map[string]interface{}{"invalid_uid": uidParam}).
			Warn("Invalid note UID format")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid note UID format")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithComponent("note-handler").
			WithFields(map[string]interface{}{"error": err.Error()}).
			Warn("Invalid request body for partial update note")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	note, err := h.noteService.PartialUpdateNote(c.Request.Context(), noteUID, &req)
	if err != nil {
		logger.WithComponent("note-handler").
			WithFields(map[string]interface{}{
				"note_uid": noteUID.String(),
				"error":    err.Error(),
			}).
			Error("Failed to partially update note")
		utils.ErrorResponse(c, http.StatusNotFound, "Note not found")
		return
	}

	logger.WithComponent("note-handler").
		WithFields(map[string]interface{}{"note_uid": noteUID.String()}).
		Info("Successfully partially updated note")

	utils.SuccessResponse(c, note, "Note updated successfully")
}

// DeleteNote handles DELETE /api/notes/:uid
func (h *NoteHandler) DeleteNote(c *gin.Context) {
	uidParam := c.Param("uid")

	logger.WithComponent("note-handler").
		WithFields(map[string]interface{}{"note_uid": uidParam}).
		Info("Deleting note")

	noteUID, err := uuid.Parse(uidParam)
	if err != nil {
		logger.WithComponent("note-handler").
			WithFields(map[string]interface{}{"invalid_uid": uidParam}).
			Warn("Invalid note UID format")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid note UID format")
		return
	}

	err = h.noteService.DeleteNote(c.Request.Context(), noteUID)
	if err != nil {
		logger.WithComponent("note-handler").
			WithFields(map[string]interface{}{
				"note_uid": noteUID.String(),
				"error":    err.Error(),
			}).
			Error("Failed to delete note")
		utils.ErrorResponse(c, http.StatusNotFound, "Note not found")
		return
	}

	logger.WithComponent("note-handler").
		WithFields(map[string]interface{}{"note_uid": noteUID.String()}).
		Info("Successfully deleted note")

	utils.SuccessResponse(c, gin.H{"message": "Note deleted successfully"}, "Note deleted successfully")
}

// MoveNoteToFolder handles POST /api/notes/:uid/move-to-folder
func (h *NoteHandler) MoveNoteToFolder(c *gin.Context) {
	noteUIDStr := c.Param("uid")
	noteUID, err := uuid.Parse(noteUIDStr)
	if err != nil {
		logger.WithComponent("note-handler").
			WithFields(map[string]interface{}{
				"note_uid": noteUIDStr,
				"error":    err.Error(),
			}).
			Error("Invalid note UID format")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid note UID format")
		return
	}

	var request struct {
		FolderID *int `json:"folder_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.WithComponent("note-handler").
			WithFields(map[string]interface{}{
				"note_uid": noteUID.String(),
				"error":    err.Error(),
			}).
			Error("Invalid request body for move note")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	logger.WithComponent("note-handler").
		WithFields(map[string]interface{}{
			"note_uid":  noteUID.String(),
			"folder_id": request.FolderID,
		}).
		Info("Moving note to folder")

	note, err := h.noteService.MoveNoteToFolder(c.Request.Context(), noteUID, request.FolderID)
	if err != nil {
		logger.WithComponent("note-handler").
			WithFields(map[string]interface{}{
				"note_uid":  noteUID.String(),
				"folder_id": request.FolderID,
				"error":     err.Error(),
			}).
			Error("Failed to move note to folder")

		if err.Error() == "note not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Note not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to move note")
		}
		return
	}

	logger.WithComponent("note-handler").
		WithFields(map[string]interface{}{
			"note_uid":  noteUID.String(),
			"folder_id": request.FolderID,
		}).
		Info("Successfully moved note to folder")

	utils.SuccessResponse(c, note, "Note moved successfully")
}
