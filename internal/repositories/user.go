package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/pkg/logger"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, user_uid, email, password_hash, name, created_at, updated_at, is_active
		FROM users
		WHERE email = $1 AND is_active = true`

	var u models.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.UserUID, &u.Email, &u.Password, &u.Name,
		&u.CreatedAt, &u.UpdatedAt, &u.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		logger.WithComponent("user-repository").
			WithFields(map[string]interface{}{
				"email": email,
				"error": err.Error(),
			}).
			Error("Failed to get user by email")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	logger.WithComponent("user-repository").
		WithFields(map[string]interface{}{
			"user_uid": u.UserUID.String(),
			"email":    email,
		}).
		Info("Successfully retrieved user by email")

	return &u, nil
}

func (r *userRepository) GetByUID(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, user_uid, email, password_hash, name, created_at, updated_at, is_active
		FROM users
		WHERE user_uid = $1 AND is_active = true`

	var u models.User
	err := r.db.QueryRow(ctx, query, uid).Scan(
		&u.ID, &u.UserUID, &u.Email, &u.Password, &u.Name,
		&u.CreatedAt, &u.UpdatedAt, &u.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		logger.WithComponent("user-repository").
			WithFields(map[string]interface{}{
				"user_uid": uid.String(),
				"error":    err.Error(),
			}).
			Error("Failed to get user by UID")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	logger.WithComponent("user-repository").
		WithFields(map[string]interface{}{
			"user_uid": uid.String(),
		}).
		Info("Successfully retrieved user by UID")

	return &u, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	// Generate UUID for new user
	user.UserUID = uuid.New()
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = user.CreatedAt
	user.IsActive = true

	query := `
		INSERT INTO users (user_uid, email, password_hash, name, created_at, updated_at, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	err := r.db.QueryRow(ctx, query,
		user.UserUID, user.Email, user.Password, user.Name,
		user.CreatedAt, user.UpdatedAt, user.IsActive,
	).Scan(&user.ID)

	if err != nil {
		logger.WithComponent("user-repository").
			WithFields(map[string]interface{}{
				"email": user.Email,
				"name":  user.Name,
				"error": err.Error(),
			}).
			Error("Failed to create user")
		return fmt.Errorf("failed to create user: %w", err)
	}

	logger.WithComponent("user-repository").
		WithFields(map[string]interface{}{
			"user_uid": user.UserUID.String(),
			"email":    user.Email,
			"name":     user.Name,
		}).
		Info("Successfully created user")

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, user_uid, email, password_hash, name, created_at, updated_at, is_active
		FROM users
		WHERE id = $1 AND is_active = true`

	var u models.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.UserUID, &u.Email, &u.Password, &u.Name,
		&u.CreatedAt, &u.UpdatedAt, &u.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		logger.WithComponent("user-repository").
			WithFields(map[string]interface{}{
				"id":    id,
				"error": err.Error(),
			}).
			Error("Failed to get user by ID")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	logger.WithComponent("user-repository").
		WithFields(map[string]interface{}{
			"user_uid": u.UserUID.String(),
			"id":       u.ID,
		}).
		Debug("Successfully retrieved user by ID")

	return &u, nil
}
