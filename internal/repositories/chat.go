package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"lucid-lists-backend/internal/models"
)

type ChatRepository struct {
	db *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) ChatRepositoryInterface {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) GetConversationsByProjectID(projectID int) ([]models.ChatConversation, error) {
	query := `
		SELECT id, conversation_uid, project_id, name, created_at, updated_at, created_by, updated_by, is_active
		FROM chat_conversations 
		WHERE project_id = $1 AND is_active = true 
		ORDER BY updated_at DESC NULLS LAST, created_at DESC
		LIMIT 10`

	rows, err := r.db.Query(context.Background(), query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []models.ChatConversation
	for rows.Next() {
		var conv models.ChatConversation
		err := rows.Scan(
			&conv.ID, &conv.ConversationUID, &conv.ProjectID, &conv.Name,
			&conv.CreatedAt, &conv.UpdatedAt, &conv.CreatedBy, &conv.UpdatedBy, &conv.IsActive,
		)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}

	return conversations, rows.Err()
}

func (r *ChatRepository) GetConversationByUID(conversationUID uuid.UUID) (*models.ChatConversation, error) {
	query := `
		SELECT id, conversation_uid, project_id, name, created_at, updated_at, created_by, updated_by, is_active
		FROM chat_conversations 
		WHERE conversation_uid = $1 AND is_active = true`

	var conversation models.ChatConversation
	err := r.db.QueryRow(context.Background(), query, conversationUID).Scan(
		&conversation.ID, &conversation.ConversationUID, &conversation.ProjectID, &conversation.Name,
		&conversation.CreatedAt, &conversation.UpdatedAt, &conversation.CreatedBy, &conversation.UpdatedBy, &conversation.IsActive,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return &conversation, nil
}

func (r *ChatRepository) CreateConversation(conversation *models.ChatConversation) error {
	// First check if we already have 10 conversations for this project
	countQuery := `
		SELECT COUNT(*) 
		FROM chat_conversations 
		WHERE project_id = $1 AND is_active = true`

	var count int
	err := r.db.QueryRow(context.Background(), countQuery, conversation.ProjectID).Scan(&count)
	if err != nil {
		return err
	}

	if count >= 10 {
		// Delete the oldest conversation to make room
		deleteQuery := `
			UPDATE chat_conversations 
			SET is_active = false, updated_at = now()
			WHERE id = (
				SELECT id FROM chat_conversations 
				WHERE project_id = $1 AND is_active = true 
				ORDER BY updated_at ASC NULLS FIRST, created_at ASC 
				LIMIT 1
			)`
		_, err = r.db.Exec(context.Background(), deleteQuery, conversation.ProjectID)
		if err != nil {
			return err
		}
	}

	query := `
		INSERT INTO chat_conversations (conversation_uid, project_id, name, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err = r.db.QueryRow(context.Background(), query,
		conversation.ConversationUID,
		conversation.ProjectID,
		conversation.Name,
		conversation.CreatedBy,
	).Scan(&conversation.ID, &conversation.CreatedAt)

	return err
}

func (r *ChatRepository) UpdateConversation(conversation *models.ChatConversation) error {
	query := `
		UPDATE chat_conversations 
		SET name = $1, updated_at = now(), updated_by = $2
		WHERE conversation_uid = $3 AND is_active = true`

	result, err := r.db.Exec(context.Background(), query, conversation.Name, conversation.UpdatedBy, conversation.ConversationUID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("conversation not found or inactive")
	}

	return nil
}

func (r *ChatRepository) DeleteConversation(conversationUID uuid.UUID) error {
	query := `
		UPDATE chat_conversations 
		SET is_active = false, updated_at = now()
		WHERE conversation_uid = $1`

	result, err := r.db.Exec(context.Background(), query, conversationUID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("conversation not found")
	}

	return nil
}

// Message methods
func (r *ChatRepository) GetMessagesByConversationID(conversationID int) ([]models.ChatMessage, error) {
	query := `
		SELECT id, message_uid, conversation_id, message_type, content, created_at, created_by
		FROM chat_messages 
		WHERE conversation_id = $1 
		ORDER BY created_at ASC`

	rows, err := r.db.Query(context.Background(), query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.ChatMessage
	for rows.Next() {
		var msg models.ChatMessage
		err := rows.Scan(
			&msg.ID, &msg.MessageUID, &msg.ConversationID, &msg.MessageType,
			&msg.Content, &msg.CreatedAt, &msg.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

func (r *ChatRepository) CreateMessage(message *models.ChatMessage) error {
	query := `
		INSERT INTO chat_messages (message_uid, conversation_id, message_type, content, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	err := r.db.QueryRow(context.Background(), query,
		message.MessageUID,
		message.ConversationID,
		message.MessageType,
		message.Content,
		message.CreatedBy,
	).Scan(&message.ID, &message.CreatedAt)

	return err
}
