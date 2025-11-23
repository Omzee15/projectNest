-- Migration: Add support for collaborative editing with Yjs CRDT
-- This migration adds fields to support real-time collaborative editing using Yjs

-- Add collaborative editing fields to the note table
ALTER TABLE note
ADD COLUMN IF NOT EXISTS yjs_state_vector BYTEA,
ADD COLUMN IF NOT EXISTS yjs_update BYTEA,
ADD COLUMN IF NOT EXISTS last_modified_by_user_uid UUID,
ADD COLUMN IF NOT EXISTS version INT DEFAULT 1,
ADD COLUMN IF NOT EXISTS content_html TEXT DEFAULT '';

-- Create index for faster lookups by project
CREATE INDEX IF NOT EXISTS idx_note_project_id ON note(project_id);
CREATE INDEX IF NOT EXISTS idx_note_is_active ON note(is_active);

-- Add comments for documentation
COMMENT ON COLUMN note.yjs_state_vector IS 'Yjs state vector for CRDT synchronization';
COMMENT ON COLUMN note.yjs_update IS 'Latest Yjs update/snapshot for the document';
COMMENT ON COLUMN note.last_modified_by_user_uid IS 'UUID of the last user who modified the note';
COMMENT ON COLUMN note.version IS 'Document version counter for conflict detection';
COMMENT ON COLUMN note.content_html IS 'Rendered HTML content for search and preview purposes';

-- Create table for tracking active collaborative sessions
CREATE TABLE IF NOT EXISTS note_sessions (
    id SERIAL PRIMARY KEY,
    session_uid UUID DEFAULT uuid_generate_v4() UNIQUE,
    note_id INT NOT NULL REFERENCES note(id) ON DELETE CASCADE,
    user_uid UUID NOT NULL,
    user_name VARCHAR NOT NULL,
    user_email VARCHAR,
    connected_at TIMESTAMP DEFAULT NOW(),
    last_seen_at TIMESTAMP DEFAULT NOW(),
    cursor_position JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true
);

-- Create indexes for session lookups
CREATE INDEX IF NOT EXISTS idx_note_sessions_note_id ON note_sessions(note_id);
CREATE INDEX IF NOT EXISTS idx_note_sessions_user_uid ON note_sessions(user_uid);
CREATE INDEX IF NOT EXISTS idx_note_sessions_is_active ON note_sessions(is_active);

-- Create table for tracking note edit history/snapshots
CREATE TABLE IF NOT EXISTS note_snapshots (
    id SERIAL PRIMARY KEY,
    snapshot_uid UUID DEFAULT uuid_generate_v4() UNIQUE,
    note_id INT NOT NULL REFERENCES note(id) ON DELETE CASCADE,
    content_json JSONB NOT NULL,
    content_html TEXT,
    yjs_state_vector BYTEA,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by_user_uid UUID,
    version INT NOT NULL,
    snapshot_type VARCHAR(50) DEFAULT 'auto', -- 'auto', 'manual', 'restore_point'
    description TEXT
);

-- Create indexes for snapshot queries
CREATE INDEX IF NOT EXISTS idx_note_snapshots_note_id ON note_snapshots(note_id);
CREATE INDEX IF NOT EXISTS idx_note_snapshots_created_at ON note_snapshots(created_at);
CREATE INDEX IF NOT EXISTS idx_note_snapshots_version ON note_snapshots(version);

-- Add comment
COMMENT ON TABLE note_sessions IS 'Tracks active collaborative editing sessions for real-time presence';
COMMENT ON TABLE note_snapshots IS 'Stores historical snapshots of note content for version history and recovery';

-- Update existing notes to have default values
UPDATE note 
SET content_html = '', version = 1 
WHERE content_html IS NULL OR version IS NULL;
