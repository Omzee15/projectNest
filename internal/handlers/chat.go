package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/internal/services"
	"lucid-lists-backend/internal/utils"
)

type ChatHandler struct {
	chatService *services.ChatService
}

func NewChatHandler(chatService *services.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

func (h *ChatHandler) GetConversations(c *gin.Context) {
	// Get project UID from URL
	projectUIDStr := c.Param("projectUid")
	projectUID, err := uuid.Parse(projectUIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project ID")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	conversations, err := h.chatService.GetConversationsByProjectUID(projectUID, userID.(int))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get conversations")
		return
	}

	utils.SuccessResponse(c, conversations, "Conversations retrieved successfully")
}

func (h *ChatHandler) GetConversationWithMessages(c *gin.Context) {
	// Get conversation UID from URL
	conversationUIDStr := c.Param("conversationUid")
	conversationUID, err := uuid.Parse(conversationUIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid conversation ID")
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	conversation, err := h.chatService.GetConversationWithMessages(conversationUID, userID.(int))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get conversation")
		return
	}

	utils.SuccessResponse(c, conversation, "Conversation retrieved successfully")
}

func (h *ChatHandler) CreateConversation(c *gin.Context) {
	// Get project UID from URL
	projectUIDStr := c.Param("projectUid")
	projectUID, err := uuid.Parse(projectUIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project ID")
		return
	}

	// Get user UID from context
	userUID, exists := c.Get("user_uid")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Parse request body
	var req models.ChatConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	conversation, err := h.chatService.CreateConversation(req, projectUID, userUID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create conversation")
		return
	}

	utils.CreatedResponse(c, conversation, "Conversation created successfully")
}

func (h *ChatHandler) CreateMessage(c *gin.Context) {
	// Get user UID from context
	userUID, exists := c.Get("user_uid")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Parse request body
	var req models.ChatMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	message, err := h.chatService.CreateMessage(req, userUID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create message")
		return
	}

	utils.CreatedResponse(c, message, "Message created successfully")
}

func (h *ChatHandler) DeleteConversation(c *gin.Context) {
	// Get conversation UID from URL
	conversationUIDStr := c.Param("conversationUid")
	conversationUID, err := uuid.Parse(conversationUIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid conversation ID")
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	err = h.chatService.DeleteConversation(conversationUID, userID.(int))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete conversation")
		return
	}

	utils.SuccessResponse(c, nil, "Conversation deleted successfully")
}
