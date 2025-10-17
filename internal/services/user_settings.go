package services

import (
	"context"
	"fmt"

	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/internal/repositories"
	"lucid-lists-backend/internal/utils"
)

type UserSettingsService struct {
	repo repositories.UserSettingsRepository
}

func NewUserSettingsService(repo repositories.UserSettingsRepository) *UserSettingsService {
	return &UserSettingsService{
		repo: repo,
	}
}

// InitializeDefaultSettings creates default settings for a new user
func (s *UserSettingsService) InitializeDefaultSettings(ctx context.Context, userID int) error {
	// Check if settings already exist
	existing, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check existing settings: %w", err)
	}

	if existing != nil {
		// Settings already exist, no need to initialize
		return nil
	}

	// Create default settings
	defaultRequest := &models.UserSettingsRequest{
		Theme:                utils.StringPtr("projectnest-default"),
		Language:             utils.StringPtr("en"),
		Timezone:             utils.StringPtr("UTC"),
		NotificationsEnabled: utils.BoolPtr(true),
		EmailNotifications:   utils.BoolPtr(true),
		SoundEnabled:         utils.BoolPtr(true),
		CompactMode:          utils.BoolPtr(false),
		AutoSave:             utils.BoolPtr(true),
		AutoSaveInterval:     utils.IntPtr(30),
	}

	return s.repo.CreateOrUpdate(ctx, userID, defaultRequest)
}

func (s *UserSettingsService) GetUserSettings(ctx context.Context, userID int) (*models.UserSettingsResponse, error) {
	settings, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user settings: %w", err)
	}

	// If no settings exist, initialize default settings
	if settings == nil {
		if err := s.InitializeDefaultSettings(ctx, userID); err != nil {
			return nil, fmt.Errorf("failed to initialize default settings: %w", err)
		}

		// Fetch the newly created settings
		settings, err = s.repo.GetByUserID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get newly created settings: %w", err)
		}
	}

	// If still no settings (shouldn't happen), return defaults
	if settings == nil {
		return &models.UserSettingsResponse{
			Theme:                "projectnest-default",
			Language:             "en",
			Timezone:             "UTC",
			NotificationsEnabled: true,
			EmailNotifications:   true,
			SoundEnabled:         true,
			CompactMode:          false,
			AutoSave:             true,
			AutoSaveInterval:     30,
		}, nil
	}

	// Convert to response DTO
	response := &models.UserSettingsResponse{
		SettingsUID:          settings.SettingsUID,
		Theme:                settings.Theme,
		Language:             settings.Language,
		Timezone:             settings.Timezone,
		NotificationsEnabled: settings.NotificationsEnabled,
		EmailNotifications:   settings.EmailNotifications,
		SoundEnabled:         settings.SoundEnabled,
		CompactMode:          settings.CompactMode,
		AutoSave:             settings.AutoSave,
		AutoSaveInterval:     settings.AutoSaveInterval,
		CreatedAt:            settings.CreatedAt,
		UpdatedAt:            settings.UpdatedAt,
	}

	return response, nil
}

func (s *UserSettingsService) UpdateUserSettings(ctx context.Context, userID int, request *models.UserSettingsRequest) (*models.UserSettingsResponse, error) {
	// Validate the request
	if err := utils.ValidateStruct(request); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Validate theme if provided
	if request.Theme != nil && *request.Theme != "" {
		validThemes := []string{
			"projectnest-default",
			"projectnest-dark",
			"solarized-light",
		}

		isValid := false
		for _, theme := range validThemes {
			if *request.Theme == theme {
				isValid = true
				break
			}
		}

		if !isValid {
			return nil, fmt.Errorf("invalid theme: %s", *request.Theme)
		}
	}

	// Validate language if provided
	if request.Language != nil && *request.Language != "" {
		validLanguages := []string{"en", "es", "fr", "de", "it", "pt", "ru", "ja", "ko", "zh"}
		isValid := false
		for _, lang := range validLanguages {
			if *request.Language == lang {
				isValid = true
				break
			}
		}

		if !isValid {
			return nil, fmt.Errorf("invalid language: %s", *request.Language)
		}
	}

	// Update or create settings
	err := s.repo.CreateOrUpdate(ctx, userID, request)
	if err != nil {
		return nil, fmt.Errorf("failed to update user settings: %w", err)
	}

	// Return updated settings
	return s.GetUserSettings(ctx, userID)
}

// ResetToDefaults resets user settings to default values
func (s *UserSettingsService) ResetToDefaults(ctx context.Context, userID int) (*models.UserSettingsResponse, error) {
	// Delete existing settings first (if any)
	if err := s.repo.DeleteByUserID(ctx, userID); err != nil {
		// Don't fail if deletion fails - settings might not exist
		fmt.Printf("Warning: failed to delete existing settings: %v\n", err)
	}

	// Initialize default settings
	if err := s.InitializeDefaultSettings(ctx, userID); err != nil {
		return nil, fmt.Errorf("failed to initialize default settings: %w", err)
	}

	// Return the new default settings
	return s.GetUserSettings(ctx, userID)
}
