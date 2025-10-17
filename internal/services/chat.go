package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/internal/repositories"
)

type ChatService struct {
	chatRepo    repositories.ChatRepositoryInterface
	projectRepo repositories.ProjectRepository
	userRepo    repositories.UserRepository
}

func NewChatService(chatRepo repositories.ChatRepositoryInterface, projectRepo repositories.ProjectRepository, userRepo repositories.UserRepository) *ChatService {
	return &ChatService{
		chatRepo:    chatRepo,
		projectRepo: projectRepo,
		userRepo:    userRepo,
	}
}

func (s *ChatService) GetConversationsByProjectUID(projectUID uuid.UUID, userID int) ([]models.ChatConversationResponse, error) {
	// First get the project to verify ownership
	project, err := s.projectRepo.GetByUID(context.TODO(), projectUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	if project == nil {
		return nil, fmt.Errorf("project not found")
	}

	// Get conversations
	conversations, err := s.chatRepo.GetConversationsByProjectID(project.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}

	// Convert to response DTOs
	var responses []models.ChatConversationResponse
	for _, conv := range conversations {
		responses = append(responses, models.ChatConversationResponse{
			ConversationUID: conv.ConversationUID,
			ProjectUID:      projectUID,
			Name:            conv.Name,
			CreatedAt:       conv.CreatedAt,
			UpdatedAt:       conv.UpdatedAt,
		})
	}

	return responses, nil
}

func (s *ChatService) GetConversationWithMessages(conversationUID uuid.UUID, userID int) (*models.ChatConversationWithMessagesResponse, error) {
	// Get conversation
	conversation, err := s.chatRepo.GetConversationByUID(conversationUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	if conversation == nil {
		return nil, fmt.Errorf("conversation not found")
	}

	// Get the project by ID to get the project UID
	project, err := s.projectRepo.GetByID(context.TODO(), conversation.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify project access: %w", err)
	}
	if project == nil {
		return nil, fmt.Errorf("project not found")
	}

	// Get messages
	messages, err := s.chatRepo.GetMessagesByConversationID(conversation.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	// Convert to response DTOs
	var messageResponses []models.ChatMessageResponse
	for _, msg := range messages {
		messageResponses = append(messageResponses, models.ChatMessageResponse{
			MessageUID:      msg.MessageUID,
			ConversationUID: conversationUID,
			MessageType:     msg.MessageType,
			Content:         msg.Content,
			CreatedAt:       msg.CreatedAt,
		})
	}

	return &models.ChatConversationWithMessagesResponse{
		ChatConversationResponse: models.ChatConversationResponse{
			ConversationUID: conversation.ConversationUID,
			ProjectUID:      project.ProjectUID,
			Name:            conversation.Name,
			CreatedAt:       conversation.CreatedAt,
			UpdatedAt:       conversation.UpdatedAt,
		},
		Messages: messageResponses,
	}, nil
}

func (s *ChatService) CreateConversation(req models.ChatConversationRequest, projectUID uuid.UUID, userUID uuid.UUID) (*models.ChatConversationResponse, error) {
	// Get project to verify ownership and get ID
	project, err := s.projectRepo.GetByUID(context.TODO(), projectUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	if project == nil {
		return nil, fmt.Errorf("project not found")
	}

	// Get user to get integer ID
	user, err := s.userRepo.GetByUID(context.TODO(), userUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Create conversation
	conversation := &models.ChatConversation{
		ConversationUID: uuid.New(),
		ProjectID:       project.ID,
		Name:            req.Name,
		CreatedBy:       &user.ID,
		IsActive:        true,
	}

	err = s.chatRepo.CreateConversation(conversation)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	return &models.ChatConversationResponse{
		ConversationUID: conversation.ConversationUID,
		ProjectUID:      projectUID,
		Name:            conversation.Name,
		CreatedAt:       conversation.CreatedAt,
		UpdatedAt:       conversation.UpdatedAt,
	}, nil
}

func (s *ChatService) CreateMessage(req models.ChatMessageRequest, userUID uuid.UUID) (*models.ChatMessageResponse, error) {
	// Get conversation to verify it exists and get ID
	conversation, err := s.chatRepo.GetConversationByUID(req.ConversationUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	if conversation == nil {
		return nil, fmt.Errorf("conversation not found")
	}

	// Get user to get integer ID
	user, err := s.userRepo.GetByUID(context.TODO(), userUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Create message (no special processing - handled on frontend)
	message := &models.ChatMessage{
		MessageUID:     uuid.New(),
		ConversationID: conversation.ID,
		MessageType:    req.MessageType,
		Content:        req.Content,
		CreatedBy:      &user.ID,
	}

	err = s.chatRepo.CreateMessage(message)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	return &models.ChatMessageResponse{
		MessageUID:      message.MessageUID,
		ConversationUID: req.ConversationUID,
		MessageType:     message.MessageType,
		Content:         message.Content,
		CreatedAt:       message.CreatedAt,
	}, nil
}

func (s *ChatService) DeleteConversation(conversationUID uuid.UUID, userID int) error {
	// Verify conversation exists and user has access
	conversation, err := s.chatRepo.GetConversationByUID(conversationUID)
	if err != nil {
		return fmt.Errorf("failed to get conversation: %w", err)
	}
	if conversation == nil {
		return fmt.Errorf("conversation not found")
	}

	// Verify user owns the project
	project, err := s.projectRepo.GetByID(context.TODO(), conversation.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to verify project access: %w", err)
	}
	if project == nil {
		return fmt.Errorf("project not found")
	}

	return s.chatRepo.DeleteConversation(conversationUID)
}
