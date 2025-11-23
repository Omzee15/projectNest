-- ProjectNest Database Schema
-- Generated from DBML script
-- Database: PostgreSQL

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create sequences for auto-increment IDs
CREATE SEQUENCE IF NOT EXISTS users_id_seq;
CREATE SEQUENCE IF NOT EXISTS project_id_seq;
CREATE SEQUENCE IF NOT EXISTS project_member_id_seq;
CREATE SEQUENCE IF NOT EXISTS list_id_seq;
CREATE SEQUENCE IF NOT EXISTS task_id_seq;
CREATE SEQUENCE IF NOT EXISTS task_assignee_id_seq;
CREATE SEQUENCE IF NOT EXISTS canvas_id_seq;
CREATE SEQUENCE IF NOT EXISTS brainstorm_canvas_id_seq;
CREATE SEQUENCE IF NOT EXISTS note_id_seq;
CREATE SEQUENCE IF NOT EXISTS notes_id_seq;
CREATE SEQUENCE IF NOT EXISTS note_folder_id_seq;
CREATE SEQUENCE IF NOT EXISTS chat_conversations_id_seq;
CREATE SEQUENCE IF NOT EXISTS chat_messages_id_seq;
CREATE SEQUENCE IF NOT EXISTS user_settings_id_seq;
CREATE SEQUENCE IF NOT EXISTS workspace_id_seq;
CREATE SEQUENCE IF NOT EXISTS workspace_member_id_seq;

-- Drop tables if they exist (in correct order due to foreign keys)
DROP TABLE IF EXISTS task_assignee CASCADE;
DROP TABLE IF EXISTS project_member CASCADE;
DROP TABLE IF EXISTS workspace_member CASCADE;
DROP TABLE IF EXISTS chat_messages CASCADE;
DROP TABLE IF EXISTS chat_conversations CASCADE;
DROP TABLE IF EXISTS task CASCADE;
DROP TABLE IF EXISTS list CASCADE;
DROP TABLE IF EXISTS note CASCADE;
DROP TABLE IF EXISTS notes CASCADE;
DROP TABLE IF EXISTS note_folder CASCADE;
DROP TABLE IF EXISTS canvas CASCADE;
DROP TABLE IF EXISTS brainstorm_canvas CASCADE;
DROP TABLE IF EXISTS project CASCADE;
DROP TABLE IF EXISTS workspace CASCADE;
DROP TABLE IF EXISTS user_settings CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS test CASCADE;

-- Create users table first (referenced by other tables)
CREATE TABLE users (
    id INT PRIMARY KEY DEFAULT nextval('users_id_seq'),
    user_uid UUID NOT NULL DEFAULT uuid_generate_v4() UNIQUE,
    email VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true
);

-- Create user_settings table
CREATE TABLE user_settings (
    id INT PRIMARY KEY DEFAULT nextval('user_settings_id_seq'),
    settings_uid UUID NOT NULL DEFAULT uuid_generate_v4() UNIQUE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    theme VARCHAR NOT NULL DEFAULT 'projectnest-default',
    language VARCHAR NOT NULL DEFAULT 'en',
    timezone VARCHAR NOT NULL DEFAULT 'UTC',
    notifications_enabled BOOLEAN NOT NULL DEFAULT true,
    email_notifications BOOLEAN NOT NULL DEFAULT true,
    sound_enabled BOOLEAN NOT NULL DEFAULT true,
    compact_mode BOOLEAN NOT NULL DEFAULT false,
    auto_save BOOLEAN NOT NULL DEFAULT true,
    auto_save_interval INT NOT NULL DEFAULT 30,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create workspace table
CREATE TABLE workspace (
    id INT PRIMARY KEY DEFAULT nextval('workspace_id_seq'),
    workspace_uid UUID DEFAULT uuid_generate_v4() UNIQUE,
    name VARCHAR NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    created_by INT REFERENCES users(id),
    updated_by INT REFERENCES users(id)
);

-- Create workspace_member table
CREATE TABLE workspace_member (
    id INT PRIMARY KEY DEFAULT nextval('workspace_member_id_seq'),
    workspace_id INT REFERENCES workspace(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR DEFAULT 'member',
    joined_at TIMESTAMP DEFAULT NOW()
);

-- Create project table
CREATE TABLE project (
    id INT PRIMARY KEY DEFAULT nextval('project_id_seq'),
    project_uid UUID DEFAULT uuid_generate_v4() UNIQUE,
    name VARCHAR NOT NULL,
    description TEXT,
    status VARCHAR DEFAULT 'active',
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    color VARCHAR DEFAULT '#FFFFFF',
    position INT,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_private BOOLEAN DEFAULT false,
    dbml_content TEXT,
    flowchart_content TEXT,
    dbml_layout_data JSONB,
    created_by INT REFERENCES users(id),
    updated_by INT REFERENCES users(id)
);

-- Create project_member table
CREATE TABLE project_member (
    id INT PRIMARY KEY DEFAULT nextval('project_member_id_seq'),
    project_id INT REFERENCES project(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR DEFAULT 'member',
    joined_at TIMESTAMP DEFAULT NOW()
);

-- Create list table
CREATE TABLE list (
    id INT PRIMARY KEY DEFAULT nextval('list_id_seq'),
    list_uid UUID DEFAULT uuid_generate_v4() UNIQUE,
    project_id INT REFERENCES project(id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    position INT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    color VARCHAR DEFAULT '#FFFFFF',
    created_by INT REFERENCES users(id),
    updated_by INT REFERENCES users(id)
);

-- Create task table
CREATE TABLE task (
    id INT PRIMARY KEY DEFAULT nextval('task_id_seq'),
    task_uid UUID DEFAULT uuid_generate_v4() UNIQUE,
    list_id INT REFERENCES list(id) ON DELETE CASCADE,
    title VARCHAR NOT NULL,
    description TEXT,
    priority VARCHAR,
    status VARCHAR DEFAULT 'todo',
    due_date TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    color VARCHAR DEFAULT '#FFFFFF',
    position INT,
    is_completed BOOLEAN DEFAULT false,
    created_by INT REFERENCES users(id),
    updated_by INT REFERENCES users(id)
);

-- Create task_assignee table
CREATE TABLE task_assignee (
    id INT PRIMARY KEY DEFAULT nextval('task_assignee_id_seq'),
    task_id INT REFERENCES task(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT NOW()
);

-- Create canvas table
CREATE TABLE canvas (
    id INT PRIMARY KEY DEFAULT nextval('canvas_id_seq'),
    canvas_uid UUID NOT NULL DEFAULT uuid_generate_v4() UNIQUE,
    project_id INT NOT NULL REFERENCES project(id) ON DELETE CASCADE,
    canvas_data JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    created_by INT REFERENCES users(id),
    updated_by INT REFERENCES users(id)
);

-- Create brainstorm_canvas table
CREATE TABLE brainstorm_canvas (
    id INT PRIMARY KEY DEFAULT nextval('brainstorm_canvas_id_seq'),
    canvas_uid UUID DEFAULT uuid_generate_v4() UNIQUE,
    project_id INT REFERENCES project(id) ON DELETE CASCADE,
    state_json JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMP,
    updated_by UUID,
    is_active BOOLEAN DEFAULT true
);

-- Create note_folder table
CREATE TABLE note_folder (
    id INT PRIMARY KEY DEFAULT nextval('note_folder_id_seq'),
    folder_uid UUID DEFAULT uuid_generate_v4() UNIQUE,
    project_id INT REFERENCES project(id) ON DELETE CASCADE,
    parent_folder_id INT REFERENCES note_folder(id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    position INT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    created_by INT REFERENCES users(id),
    updated_by INT REFERENCES users(id)
);

-- Create note table
CREATE TABLE note (
    id INT PRIMARY KEY DEFAULT nextval('note_id_seq'),
    note_uid UUID DEFAULT uuid_generate_v4() UNIQUE,
    project_id INT REFERENCES project(id) ON DELETE CASCADE,
    title VARCHAR NOT NULL,
    content_json JSONB NOT NULL DEFAULT '{"blocks": []}',
    position INT,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMP,
    updated_by UUID,
    is_active BOOLEAN DEFAULT true,
    folder_id INT REFERENCES note_folder(id) ON DELETE SET NULL
);

-- Create notes table (appears to be a duplicate/different version)
CREATE TABLE notes (
    id INT PRIMARY KEY DEFAULT nextval('notes_id_seq'),
    note_uid UUID NOT NULL DEFAULT uuid_generate_v4() UNIQUE,
    project_id INT NOT NULL REFERENCES project(id) ON DELETE CASCADE,
    folder_id INT REFERENCES note_folder(id) ON DELETE SET NULL,
    title VARCHAR NOT NULL,
    content TEXT DEFAULT '',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    created_by INT REFERENCES users(id),
    updated_by INT REFERENCES users(id)
);

-- Create chat_conversations table
CREATE TABLE chat_conversations (
    id INT PRIMARY KEY DEFAULT nextval('chat_conversations_id_seq'),
    conversation_uid UUID NOT NULL DEFAULT uuid_generate_v4() UNIQUE,
    project_id INT NOT NULL REFERENCES project(id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    created_by INT REFERENCES users(id),
    updated_by INT REFERENCES users(id)
);

-- Create chat_messages table
CREATE TABLE chat_messages (
    id INT PRIMARY KEY DEFAULT nextval('chat_messages_id_seq'),
    message_uid UUID NOT NULL DEFAULT uuid_generate_v4() UNIQUE,
    conversation_id INT NOT NULL REFERENCES chat_conversations(id) ON DELETE CASCADE,
    message_type VARCHAR NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INT REFERENCES users(id)
);

-- Create test table (as specified in DBML)
CREATE TABLE test (
    test INT
);

-- Create indexes for better performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_user_uid ON users(user_uid);
CREATE INDEX idx_project_user_id ON project(user_id);
CREATE INDEX idx_project_project_uid ON project(project_uid);
CREATE INDEX idx_list_project_id ON list(project_id);
CREATE INDEX idx_task_list_id ON task(list_id);
CREATE INDEX idx_canvas_project_id ON canvas(project_id);
CREATE INDEX idx_note_project_id ON note(project_id);
CREATE INDEX idx_notes_project_id ON notes(project_id);
CREATE INDEX idx_chat_conversations_project_id ON chat_conversations(project_id);
CREATE INDEX idx_chat_messages_conversation_id ON chat_messages(conversation_id);
CREATE INDEX idx_user_settings_user_id ON user_settings(user_id);

-- Create triggers for updated_at timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply triggers to relevant tables
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_project_updated_at BEFORE UPDATE ON project FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_list_updated_at BEFORE UPDATE ON list FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_task_updated_at BEFORE UPDATE ON task FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_canvas_updated_at BEFORE UPDATE ON canvas FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_note_updated_at BEFORE UPDATE ON note FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_notes_updated_at BEFORE UPDATE ON notes FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_note_folder_updated_at BEFORE UPDATE ON note_folder FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_chat_conversations_updated_at BEFORE UPDATE ON chat_conversations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_workspace_updated_at BEFORE UPDATE ON workspace FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_user_settings_updated_at BEFORE UPDATE ON user_settings FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add comments to tables
COMMENT ON TABLE users IS 'User accounts and authentication information';
COMMENT ON TABLE project IS 'Projects created by users';
COMMENT ON TABLE list IS 'Task lists within projects';
COMMENT ON TABLE task IS 'Individual tasks within lists';
COMMENT ON TABLE canvas IS 'Canvas data for visual project planning';
COMMENT ON TABLE note IS 'Notes and documentation within projects';
COMMENT ON TABLE chat_conversations IS 'Chat conversations within projects';
COMMENT ON TABLE chat_messages IS 'Individual chat messages';
COMMENT ON TABLE user_settings IS 'User preferences and settings';

-- Success message
SELECT 'ProjectNest database schema created successfully!' as result;