-- Migration: Add user_settings table
-- This table stores user preferences including theme, language, and other settings

CREATE TABLE IF NOT EXISTS user_settings (
    id SERIAL PRIMARY KEY,
    settings_uid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    user_id INTEGER NOT NULL,
    theme VARCHAR(100) NOT NULL DEFAULT 'projectnest-default',
    language VARCHAR(2) NOT NULL DEFAULT 'en',
    timezone VARCHAR(100) NOT NULL DEFAULT 'UTC',
    notifications_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    email_notifications BOOLEAN NOT NULL DEFAULT TRUE,
    sound_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    compact_mode BOOLEAN NOT NULL DEFAULT FALSE,
    auto_save BOOLEAN NOT NULL DEFAULT TRUE,
    auto_save_interval INTEGER NOT NULL DEFAULT 30 CHECK (auto_save_interval >= 10 AND auto_save_interval <= 600),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Foreign key constraint (assuming users table has an id column)
    CONSTRAINT fk_user_settings_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    -- Unique constraint to ensure one settings record per user
    CONSTRAINT unique_user_settings UNIQUE (user_id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_user_settings_user_id ON user_settings(user_id);
CREATE INDEX IF NOT EXISTS idx_user_settings_settings_uid ON user_settings(settings_uid);
CREATE INDEX IF NOT EXISTS idx_user_settings_theme ON user_settings(theme);

-- Add comments for documentation
COMMENT ON TABLE user_settings IS 'Stores user preferences and settings';
COMMENT ON COLUMN user_settings.theme IS 'Theme name (e.g., projectnest-default, projectnest-dark, github-dark)';
COMMENT ON COLUMN user_settings.language IS 'ISO 639-1 language code (e.g., en, es, fr)';
COMMENT ON COLUMN user_settings.timezone IS 'IANA timezone identifier (e.g., America/New_York, UTC)';
COMMENT ON COLUMN user_settings.auto_save_interval IS 'Auto-save interval in seconds (10-600)';