package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/internal/repositories"
	"lucid-lists-backend/pkg/logger"
)

type NoteService struct {
	noteRepo    repositories.NoteRepository
	projectRepo repositories.ProjectRepository
}

func NewNoteService(noteRepo repositories.NoteRepository, projectRepo repositories.ProjectRepository) *NoteService {
	return &NoteService{
		noteRepo:    noteRepo,
		projectRepo: projectRepo,
	}
}

// GetNotesByProject retrieves all notes for a project
func (s *NoteService) GetNotesByProject(ctx context.Context, projectUID uuid.UUID) (*models.NotesResponse, error) {
	logger.WithComponent("note-service").
		WithFields(map[string]interface{}{"project_uid": projectUID.String()}).
		Info("Getting notes for project")

	// First check if project exists
	_, err := s.projectRepo.GetByUID(ctx, projectUID)
	if err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Project not found")
		return nil, fmt.Errorf("project not found")
	}

	// Get notes
	notes, err := s.noteRepo.GetByProjectUID(ctx, projectUID)
	if err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Failed to get notes")
		return nil, err
	}

	// Convert to response format
	noteResponses := make([]models.NoteResponse, len(notes))
	for i, note := range notes {
		var contentJSON models.NoteContent
		if err := json.Unmarshal([]byte(note.ContentJSON), &contentJSON); err != nil {
			logger.WithComponent("note-service").
				WithFields(map[string]interface{}{
					"note_uid": note.NoteUID.String(),
					"error":    err.Error(),
				}).
				Error("Failed to unmarshal note content JSON")
			// Use empty content on error
			contentJSON = models.NoteContent{Blocks: []models.NoteBlock{}}
		}

		noteResponses[i] = models.NoteResponse{
			NoteUID:     note.NoteUID,
			ProjectID:   note.ProjectID,
			Title:       note.Title,
			ContentJSON: contentJSON,
			FolderID:    note.FolderID,
			Position:    note.Position,
			CreatedAt:   note.CreatedAt,
			UpdatedAt:   note.UpdatedAt,
		}
	}

	response := &models.NotesResponse{
		Notes: noteResponses,
		Total: len(noteResponses),
	}

	logger.WithComponent("note-service").
		WithFields(map[string]interface{}{
			"project_uid": projectUID.String(),
			"count":       len(noteResponses),
		}).
		Info("Successfully retrieved notes")

	return response, nil
}

// GetNote retrieves a single note by UID
func (s *NoteService) GetNote(ctx context.Context, noteUID uuid.UUID) (*models.NoteResponse, error) {
	logger.WithComponent("note-service").
		WithFields(map[string]interface{}{"note_uid": noteUID.String()}).
		Info("Getting note")

	note, err := s.noteRepo.GetByUID(ctx, noteUID)
	if err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"note_uid": noteUID.String(),
				"error":    err.Error(),
			}).
			Error("Failed to get note")
		return nil, err
	}

	var contentJSON models.NoteContent
	if err := json.Unmarshal([]byte(note.ContentJSON), &contentJSON); err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"note_uid": note.NoteUID.String(),
				"error":    err.Error(),
			}).
			Error("Failed to unmarshal note content JSON")
		// Use empty content on error
		contentJSON = models.NoteContent{Blocks: []models.NoteBlock{}}
	}

	response := &models.NoteResponse{
		NoteUID:     note.NoteUID,
		ProjectID:   note.ProjectID,
		Title:       note.Title,
		ContentJSON: contentJSON,
		FolderID:    note.FolderID,
		Position:    note.Position,
		CreatedAt:   note.CreatedAt,
		UpdatedAt:   note.UpdatedAt,
	}

	logger.WithComponent("note-service").
		WithFields(map[string]interface{}{"note_uid": noteUID.String()}).
		Info("Successfully retrieved note")

	return response, nil
}

// CreateNote creates a new note in a project
func (s *NoteService) CreateNote(ctx context.Context, projectUID uuid.UUID, request *models.NoteRequest) (*models.NoteResponse, error) {
	logger.WithComponent("note-service").
		WithFields(map[string]interface{}{
			"project_uid": projectUID.String(),
			"title":       request.Title,
		}).
		Info("Creating note")

	// First check if project exists
	project, err := s.projectRepo.GetByUID(ctx, projectUID)
	if err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Project not found")
		return nil, fmt.Errorf("project not found")
	}

	// Set position if not provided
	position := 0
	if request.Position != nil {
		position = *request.Position
	} else {
		// Get next available position
		maxPos, err := s.noteRepo.GetMaxPositionByProject(ctx, project.ID)
		if err != nil {
			logger.WithComponent("note-service").
				WithFields(map[string]interface{}{
					"project_id": project.ID,
					"error":      err.Error(),
				}).
				Error("Failed to get max position")
			return nil, err
		}
		position = maxPos
	}

	// Serialize content to JSON
	contentJSON, err := json.Marshal(request.ContentJSON)
	if err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"project_uid": projectUID.String(),
				"error":       err.Error(),
			}).
			Error("Failed to marshal note content JSON")
		return nil, fmt.Errorf("invalid note content format")
	}

	// Create note
	note := &models.Note{
		ProjectID:   project.ID,
		Title:       request.Title,
		ContentJSON: string(contentJSON),
		FolderID:    request.FolderID,
		Position:    &position,
	}

	err = s.noteRepo.Create(ctx, note)
	if err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"project_id": project.ID,
				"title":      request.Title,
				"error":      err.Error(),
			}).
			Error("Failed to create note")
		return nil, err
	}

	// Convert back to response format
	var responseContentJSON models.NoteContent
	if err := json.Unmarshal([]byte(note.ContentJSON), &responseContentJSON); err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"note_uid": note.NoteUID.String(),
				"error":    err.Error(),
			}).
			Error("Failed to unmarshal created note content JSON")
		responseContentJSON = models.NoteContent{Blocks: []models.NoteBlock{}}
	}

	response := &models.NoteResponse{
		NoteUID:     note.NoteUID,
		ProjectID:   note.ProjectID,
		Title:       note.Title,
		ContentJSON: responseContentJSON,
		FolderID:    note.FolderID,
		Position:    note.Position,
		CreatedAt:   note.CreatedAt,
		UpdatedAt:   note.UpdatedAt,
	}

	logger.WithComponent("note-service").
		WithFields(map[string]interface{}{
			"note_uid": note.NoteUID.String(),
			"title":    note.Title,
		}).
		Info("Successfully created note")

	return response, nil
}

// UpdateNote updates a note completely
func (s *NoteService) UpdateNote(ctx context.Context, noteUID uuid.UUID, request *models.NoteRequest) (*models.NoteResponse, error) {
	logger.WithComponent("note-service").
		WithFields(map[string]interface{}{
			"note_uid": noteUID.String(),
			"title":    request.Title,
		}).
		Info("Updating note")

	// Check if note exists
	existingNote, err := s.noteRepo.GetByUID(ctx, noteUID)
	if err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"note_uid": noteUID.String(),
				"error":    err.Error(),
			}).
			Error("Note not found")
		return nil, err
	}

	// Set position if not provided
	position := existingNote.Position
	if request.Position != nil {
		position = request.Position
	}

	// Serialize content to JSON
	contentJSON, err := json.Marshal(request.ContentJSON)
	if err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"note_uid": noteUID.String(),
				"error":    err.Error(),
			}).
			Error("Failed to marshal note content JSON")
		return nil, fmt.Errorf("invalid note content format")
	}

	// Update note
	updatedNote := &models.Note{
		Title:       request.Title,
		ContentJSON: string(contentJSON),
		Position:    position,
	}

	err = s.noteRepo.Update(ctx, noteUID, updatedNote)
	if err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"note_uid": noteUID.String(),
				"error":    err.Error(),
			}).
			Error("Failed to update note")
		return nil, err
	}

	// Get updated note
	note, err := s.noteRepo.GetByUID(ctx, noteUID)
	if err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"note_uid": noteUID.String(),
				"error":    err.Error(),
			}).
			Error("Failed to get updated note")
		return nil, err
	}

	// Convert back to response format
	var responseContentJSON models.NoteContent
	if err := json.Unmarshal([]byte(note.ContentJSON), &responseContentJSON); err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"note_uid": note.NoteUID.String(),
				"error":    err.Error(),
			}).
			Error("Failed to unmarshal updated note content JSON")
		responseContentJSON = models.NoteContent{Blocks: []models.NoteBlock{}}
	}

	response := &models.NoteResponse{
		NoteUID:     note.NoteUID,
		ProjectID:   note.ProjectID,
		Title:       note.Title,
		ContentJSON: responseContentJSON,
		FolderID:    note.FolderID,
		Position:    note.Position,
		CreatedAt:   note.CreatedAt,
		UpdatedAt:   note.UpdatedAt,
	}

	logger.WithComponent("note-service").
		WithFields(map[string]interface{}{"note_uid": noteUID.String()}).
		Info("Successfully updated note")

	return response, nil
}

// PartialUpdateNote updates specific fields of a note
func (s *NoteService) PartialUpdateNote(ctx context.Context, noteUID uuid.UUID, request *models.NoteUpdateRequest) (*models.NoteResponse, error) {
	logger.WithComponent("note-service").
		WithFields(map[string]interface{}{"note_uid": noteUID.String()}).
		Info("Partially updating note")

	// Check if note exists
	_, err := s.noteRepo.GetByUID(ctx, noteUID)
	if err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"note_uid": noteUID.String(),
				"error":    err.Error(),
			}).
			Error("Note not found")
		return nil, err
	}

	// Prepare the update request for repository - convert ContentJSON if provided
	repoUpdateRequest := models.NoteUpdateRequest{
		Title:    request.Title,
		Position: request.Position,
	}

	// Handle ContentJSON marshaling if provided
	if request.ContentJSON != nil {
		contentJSON, err := json.Marshal(*request.ContentJSON)
		if err != nil {
			logger.WithComponent("note-service").
				WithFields(map[string]interface{}{
					"note_uid": noteUID.String(),
					"error":    err.Error(),
				}).
				Error("Failed to marshal note content JSON for partial update")
			return nil, fmt.Errorf("invalid note content format")
		}
		contentJSONStr := string(contentJSON)
		repoUpdateRequest.ContentJSONString = &contentJSONStr
	}

	// Update note
	err = s.noteRepo.PartialUpdate(ctx, noteUID, repoUpdateRequest)
	if err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"note_uid": noteUID.String(),
				"error":    err.Error(),
			}).
			Error("Failed to partially update note")
		return nil, err
	}

	// Get updated note
	note, err := s.noteRepo.GetByUID(ctx, noteUID)
	if err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"note_uid": noteUID.String(),
				"error":    err.Error(),
			}).
			Error("Failed to get updated note")
		return nil, err
	}

	// Convert back to response format
	var responseContentJSON models.NoteContent
	if err := json.Unmarshal([]byte(note.ContentJSON), &responseContentJSON); err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"note_uid": note.NoteUID.String(),
				"error":    err.Error(),
			}).
			Error("Failed to unmarshal partially updated note content JSON")
		responseContentJSON = models.NoteContent{Blocks: []models.NoteBlock{}}
	}

	response := &models.NoteResponse{
		NoteUID:     note.NoteUID,
		ProjectID:   note.ProjectID,
		Title:       note.Title,
		ContentJSON: responseContentJSON,
		FolderID:    note.FolderID,
		Position:    note.Position,
		CreatedAt:   note.CreatedAt,
		UpdatedAt:   note.UpdatedAt,
	}

	logger.WithComponent("note-service").
		WithFields(map[string]interface{}{"note_uid": noteUID.String()}).
		Info("Successfully partially updated note")

	return response, nil
}

// DeleteNote soft deletes a note
func (s *NoteService) DeleteNote(ctx context.Context, noteUID uuid.UUID) error {
	logger.WithComponent("note-service").
		WithFields(map[string]interface{}{"note_uid": noteUID.String()}).
		Info("Deleting note")

	// Check if note exists
	_, err := s.noteRepo.GetByUID(ctx, noteUID)
	if err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"note_uid": noteUID.String(),
				"error":    err.Error(),
			}).
			Error("Note not found")
		return err
	}

	err = s.noteRepo.Delete(ctx, noteUID)
	if err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"note_uid": noteUID.String(),
				"error":    err.Error(),
			}).
			Error("Failed to delete note")
		return err
	}

	logger.WithComponent("note-service").
		WithFields(map[string]interface{}{"note_uid": noteUID.String()}).
		Info("Successfully deleted note")

	return nil
}

// MoveNoteToFolder moves a note to a specific folder or removes it from folder (if folderID is nil)
func (s *NoteService) MoveNoteToFolder(ctx context.Context, noteUID uuid.UUID, folderID *int) (*models.NoteResponse, error) {
	logger.WithComponent("note-service").
		WithFields(map[string]interface{}{
			"note_uid":  noteUID.String(),
			"folder_id": folderID,
		}).
		Info("Moving note to folder")

	// Check if note exists
	_, err := s.noteRepo.GetByUID(ctx, noteUID)
	if err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"note_uid": noteUID.String(),
				"error":    err.Error(),
			}).
			Error("Note not found")
		return nil, fmt.Errorf("note not found")
	}

	// Use the specialized method for updating folder_id
	err = s.noteRepo.UpdateFolderID(ctx, noteUID, folderID)
	if err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"note_uid":  noteUID.String(),
				"folder_id": folderID,
				"error":     err.Error(),
			}).
			Error("Failed to move note to folder")
		return nil, err
	}

	// Fetch the updated note
	updatedNote, err := s.noteRepo.GetByUID(ctx, noteUID)
	if err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"note_uid": noteUID.String(),
				"error":    err.Error(),
			}).
			Error("Failed to fetch updated note")
		return nil, err
	}

	// Convert to response format
	var responseContentJSON models.NoteContent
	if err := json.Unmarshal([]byte(updatedNote.ContentJSON), &responseContentJSON); err != nil {
		logger.WithComponent("note-service").
			WithFields(map[string]interface{}{
				"note_uid": noteUID.String(),
				"error":    err.Error(),
			}).
			Error("Failed to unmarshal note content JSON")
		responseContentJSON = models.NoteContent{Blocks: []models.NoteBlock{}}
	}

	response := &models.NoteResponse{
		NoteUID:     updatedNote.NoteUID,
		ProjectID:   updatedNote.ProjectID,
		Title:       updatedNote.Title,
		ContentJSON: responseContentJSON,
		FolderID:    updatedNote.FolderID,
		Position:    updatedNote.Position,
		CreatedAt:   updatedNote.CreatedAt,
		UpdatedAt:   updatedNote.UpdatedAt,
	}

	logger.WithComponent("note-service").
		WithFields(map[string]interface{}{
			"note_uid":  noteUID.String(),
			"folder_id": folderID,
		}).
		Info("Successfully moved note to folder")

	return response, nil
}
