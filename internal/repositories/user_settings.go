package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"lucid-lists-backend/internal/models"
)

type userSettingsRepository struct {
	db *pgxpool.Pool
}

func NewUserSettingsRepository(db *pgxpool.Pool) UserSettingsRepository {
	return &userSettingsRepository{db: db}
}

func (r *userSettingsRepository) GetByUserID(ctx context.Context, userID int) (*models.UserSettings, error) {
	query := `
		SELECT id, settings_uid, user_id, theme, language, timezone, 
		       notifications_enabled, email_notifications, sound_enabled, 
		       compact_mode, auto_save, auto_save_interval, created_at, updated_at
		FROM user_settings 
		WHERE user_id = $1`

	settings := &models.UserSettings{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&settings.ID,
		&settings.SettingsUID,
		&settings.UserID,
		&settings.Theme,
		&settings.Language,
		&settings.Timezone,
		&settings.NotificationsEnabled,
		&settings.EmailNotifications,
		&settings.SoundEnabled,
		&settings.CompactMode,
		&settings.AutoSave,
		&settings.AutoSaveInterval,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Settings not found, return nil without error
		}
		return nil, err
	}

	return settings, nil
}

func (r *userSettingsRepository) Create(ctx context.Context, settings *models.UserSettings) error {
	query := `
		INSERT INTO user_settings (
			settings_uid, user_id, theme, language, timezone,
			notifications_enabled, email_notifications, sound_enabled,
			compact_mode, auto_save, auto_save_interval, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id`

	err := r.db.QueryRow(ctx, query,
		settings.SettingsUID,
		settings.UserID,
		settings.Theme,
		settings.Language,
		settings.Timezone,
		settings.NotificationsEnabled,
		settings.EmailNotifications,
		settings.SoundEnabled,
		settings.CompactMode,
		settings.AutoSave,
		settings.AutoSaveInterval,
		settings.CreatedAt,
		settings.UpdatedAt,
	).Scan(&settings.ID)

	return err
}

func (r *userSettingsRepository) Update(ctx context.Context, userID int, request *models.UserSettingsRequest) error {
	query := `
		UPDATE user_settings SET 
			theme = COALESCE($2, theme),
			language = COALESCE($3, language),
			timezone = COALESCE($4, timezone),
			notifications_enabled = COALESCE($5, notifications_enabled),
			email_notifications = COALESCE($6, email_notifications),
			sound_enabled = COALESCE($7, sound_enabled),
			compact_mode = COALESCE($8, compact_mode),
			auto_save = COALESCE($9, auto_save),
			auto_save_interval = COALESCE($10, auto_save_interval),
			updated_at = $11
		WHERE user_id = $1`

	// Convert pointer fields to values for COALESCE to work properly
	var theme, language, timezone interface{}
	var notificationsEnabled, emailNotifications, soundEnabled, compactMode, autoSave interface{}
	var autoSaveInterval interface{}

	if request.Theme != nil {
		theme = *request.Theme
	}
	if request.Language != nil {
		language = *request.Language
	}
	if request.Timezone != nil {
		timezone = *request.Timezone
	}
	if request.NotificationsEnabled != nil {
		notificationsEnabled = *request.NotificationsEnabled
	}
	if request.EmailNotifications != nil {
		emailNotifications = *request.EmailNotifications
	}
	if request.SoundEnabled != nil {
		soundEnabled = *request.SoundEnabled
	}
	if request.CompactMode != nil {
		compactMode = *request.CompactMode
	}
	if request.AutoSave != nil {
		autoSave = *request.AutoSave
	}
	if request.AutoSaveInterval != nil {
		autoSaveInterval = *request.AutoSaveInterval
	}

	_, err := r.db.Exec(ctx, query,
		userID,
		theme,
		language,
		timezone,
		notificationsEnabled,
		emailNotifications,
		soundEnabled,
		compactMode,
		autoSave,
		autoSaveInterval,
		time.Now(),
	)

	return err
}

func (r *userSettingsRepository) CreateOrUpdate(ctx context.Context, userID int, request *models.UserSettingsRequest) error {
	// First try to get existing settings
	existing, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if existing == nil {
		// Create new settings with defaults
		settings := &models.UserSettings{
			SettingsUID:          uuid.New(),
			UserID:               userID,
			Theme:                "projectnest-default",
			Language:             "en",
			Timezone:             "UTC",
			NotificationsEnabled: true,
			EmailNotifications:   true,
			SoundEnabled:         true,
			CompactMode:          false,
			AutoSave:             true,
			AutoSaveInterval:     30,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}

		// Override with request values if provided
		if request.Theme != nil {
			settings.Theme = *request.Theme
		}
		if request.Language != nil {
			settings.Language = *request.Language
		}
		if request.Timezone != nil {
			settings.Timezone = *request.Timezone
		}
		if request.NotificationsEnabled != nil {
			settings.NotificationsEnabled = *request.NotificationsEnabled
		}
		if request.EmailNotifications != nil {
			settings.EmailNotifications = *request.EmailNotifications
		}
		if request.SoundEnabled != nil {
			settings.SoundEnabled = *request.SoundEnabled
		}
		if request.CompactMode != nil {
			settings.CompactMode = *request.CompactMode
		}
		if request.AutoSave != nil {
			settings.AutoSave = *request.AutoSave
		}
		if request.AutoSaveInterval != nil {
			settings.AutoSaveInterval = *request.AutoSaveInterval
		}

		return r.Create(ctx, settings)
	}

	// Update existing settings
	return r.Update(ctx, userID, request)
}

func (r *userSettingsRepository) DeleteByUserID(ctx context.Context, userID int) error {
	query := `DELETE FROM user_settings WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}
