-- Migration to standardize all foreign key relationships to use integer IDs instead of UUIDs
-- This will improve query performance and simplify joins

BEGIN;

-- 1. CANVAS table - change created_by and updated_by from user_uid to user id
ALTER TABLE canvas ADD COLUMN created_by_id INTEGER;
ALTER TABLE canvas ADD COLUMN updated_by_id INTEGER;

-- Populate the new columns
UPDATE canvas SET created_by_id = (SELECT id FROM users WHERE user_uid = canvas.created_by);
UPDATE canvas SET updated_by_id = (SELECT id FROM users WHERE user_uid = canvas.updated_by);

-- Drop old foreign key constraints and add new ones
ALTER TABLE canvas DROP CONSTRAINT IF EXISTS canvas_created_by_fkey;
ALTER TABLE canvas DROP CONSTRAINT IF EXISTS canvas_updated_by_fkey;
ALTER TABLE canvas ADD CONSTRAINT canvas_created_by_id_fkey FOREIGN KEY (created_by_id) REFERENCES users(id);
ALTER TABLE canvas ADD CONSTRAINT canvas_updated_by_id_fkey FOREIGN KEY (updated_by_id) REFERENCES users(id);

-- Drop old UUID columns and rename new ones
ALTER TABLE canvas DROP COLUMN created_by;
ALTER TABLE canvas DROP COLUMN updated_by;
ALTER TABLE canvas RENAME COLUMN created_by_id TO created_by;
ALTER TABLE canvas RENAME COLUMN updated_by_id TO updated_by;

-- 2. CHAT_CONVERSATIONS table - change created_by and updated_by
ALTER TABLE chat_conversations ADD COLUMN created_by_id INTEGER;
ALTER TABLE chat_conversations ADD COLUMN updated_by_id INTEGER;

UPDATE chat_conversations SET created_by_id = (SELECT id FROM users WHERE user_uid = chat_conversations.created_by);
UPDATE chat_conversations SET updated_by_id = (SELECT id FROM users WHERE user_uid = chat_conversations.updated_by);

ALTER TABLE chat_conversations DROP CONSTRAINT IF EXISTS chat_conversations_created_by_fkey;
ALTER TABLE chat_conversations DROP CONSTRAINT IF EXISTS chat_conversations_updated_by_fkey;
ALTER TABLE chat_conversations ADD CONSTRAINT chat_conversations_created_by_id_fkey FOREIGN KEY (created_by_id) REFERENCES users(id);
ALTER TABLE chat_conversations ADD CONSTRAINT chat_conversations_updated_by_id_fkey FOREIGN KEY (updated_by_id) REFERENCES users(id);

ALTER TABLE chat_conversations DROP COLUMN created_by;
ALTER TABLE chat_conversations DROP COLUMN updated_by;
ALTER TABLE chat_conversations RENAME COLUMN created_by_id TO created_by;
ALTER TABLE chat_conversations RENAME COLUMN updated_by_id TO updated_by;

-- 3. CHAT_MESSAGES table - change created_by
ALTER TABLE chat_messages ADD COLUMN created_by_id INTEGER;

UPDATE chat_messages SET created_by_id = (SELECT id FROM users WHERE user_uid = chat_messages.created_by);

ALTER TABLE chat_messages DROP CONSTRAINT IF EXISTS chat_messages_created_by_fkey;
ALTER TABLE chat_messages ADD CONSTRAINT chat_messages_created_by_id_fkey FOREIGN KEY (created_by_id) REFERENCES users(id);

ALTER TABLE chat_messages DROP COLUMN created_by;
ALTER TABLE chat_messages RENAME COLUMN created_by_id TO created_by;

-- 4. LIST table - change created_by and updated_by
ALTER TABLE list ADD COLUMN created_by_id INTEGER;
ALTER TABLE list ADD COLUMN updated_by_id INTEGER;

UPDATE list SET created_by_id = (SELECT id FROM users WHERE user_uid = list.created_by);
UPDATE list SET updated_by_id = (SELECT id FROM users WHERE user_uid = list.updated_by);

ALTER TABLE list DROP CONSTRAINT IF EXISTS list_created_by_fkey;
ALTER TABLE list DROP CONSTRAINT IF EXISTS list_updated_by_fkey;
ALTER TABLE list ADD CONSTRAINT list_created_by_id_fkey FOREIGN KEY (created_by_id) REFERENCES users(id);
ALTER TABLE list ADD CONSTRAINT list_updated_by_id_fkey FOREIGN KEY (updated_by_id) REFERENCES users(id);

ALTER TABLE list DROP COLUMN created_by;
ALTER TABLE list DROP COLUMN updated_by;
ALTER TABLE list RENAME COLUMN created_by_id TO created_by;
ALTER TABLE list RENAME COLUMN updated_by_id TO updated_by;

-- 5. NOTE_FOLDER table - change created_by and updated_by
ALTER TABLE note_folder ADD COLUMN created_by_id INTEGER;
ALTER TABLE note_folder ADD COLUMN updated_by_id INTEGER;

UPDATE note_folder SET created_by_id = (SELECT id FROM users WHERE user_uid = note_folder.created_by);
UPDATE note_folder SET updated_by_id = (SELECT id FROM users WHERE user_uid = note_folder.updated_by);

ALTER TABLE note_folder DROP CONSTRAINT IF EXISTS note_folder_created_by_fkey;
ALTER TABLE note_folder DROP CONSTRAINT IF EXISTS note_folder_updated_by_fkey;
ALTER TABLE note_folder ADD CONSTRAINT note_folder_created_by_id_fkey FOREIGN KEY (created_by_id) REFERENCES users(id);
ALTER TABLE note_folder ADD CONSTRAINT note_folder_updated_by_id_fkey FOREIGN KEY (updated_by_id) REFERENCES users(id);

ALTER TABLE note_folder DROP COLUMN created_by;
ALTER TABLE note_folder DROP COLUMN updated_by;
ALTER TABLE note_folder RENAME COLUMN created_by_id TO created_by;
ALTER TABLE note_folder RENAME COLUMN updated_by_id TO updated_by;

-- 6. NOTES table - change created_by and updated_by
ALTER TABLE notes ADD COLUMN created_by_id INTEGER;
ALTER TABLE notes ADD COLUMN updated_by_id INTEGER;

UPDATE notes SET created_by_id = (SELECT id FROM users WHERE user_uid = notes.created_by);
UPDATE notes SET updated_by_id = (SELECT id FROM users WHERE user_uid = notes.updated_by);

ALTER TABLE notes DROP CONSTRAINT IF EXISTS notes_created_by_fkey;
ALTER TABLE notes DROP CONSTRAINT IF EXISTS notes_updated_by_fkey;
ALTER TABLE notes ADD CONSTRAINT notes_created_by_id_fkey FOREIGN KEY (created_by_id) REFERENCES users(id);
ALTER TABLE notes ADD CONSTRAINT notes_updated_by_id_fkey FOREIGN KEY (updated_by_id) REFERENCES users(id);

ALTER TABLE notes DROP COLUMN created_by;
ALTER TABLE notes DROP COLUMN updated_by;
ALTER TABLE notes RENAME COLUMN created_by_id TO created_by;
ALTER TABLE notes RENAME COLUMN updated_by_id TO updated_by;

-- 7. PROJECT table - change created_by and updated_by
ALTER TABLE project ADD COLUMN created_by_id INTEGER;
ALTER TABLE project ADD COLUMN updated_by_id INTEGER;

UPDATE project SET created_by_id = (SELECT id FROM users WHERE user_uid = project.created_by);
UPDATE project SET updated_by_id = (SELECT id FROM users WHERE user_uid = project.updated_by);

ALTER TABLE project DROP CONSTRAINT IF EXISTS project_created_by_fkey;
ALTER TABLE project DROP CONSTRAINT IF EXISTS project_updated_by_fkey;
ALTER TABLE project ADD CONSTRAINT project_created_by_id_fkey FOREIGN KEY (created_by_id) REFERENCES users(id);
ALTER TABLE project ADD CONSTRAINT project_updated_by_id_fkey FOREIGN KEY (updated_by_id) REFERENCES users(id);

ALTER TABLE project DROP COLUMN created_by;
ALTER TABLE project DROP COLUMN updated_by;
ALTER TABLE project RENAME COLUMN created_by_id TO created_by;
ALTER TABLE project RENAME COLUMN updated_by_id TO updated_by;

-- 8. TASK table - change created_by and updated_by
ALTER TABLE task ADD COLUMN created_by_id INTEGER;
ALTER TABLE task ADD COLUMN updated_by_id INTEGER;

UPDATE task SET created_by_id = (SELECT id FROM users WHERE user_uid = task.created_by);
UPDATE task SET updated_by_id = (SELECT id FROM users WHERE user_uid = task.updated_by);

ALTER TABLE task DROP CONSTRAINT IF EXISTS fk_task_created_by;
ALTER TABLE task DROP CONSTRAINT IF EXISTS fk_task_updated_by;
ALTER TABLE task ADD CONSTRAINT task_created_by_id_fkey FOREIGN KEY (created_by_id) REFERENCES users(id);
ALTER TABLE task ADD CONSTRAINT task_updated_by_id_fkey FOREIGN KEY (updated_by_id) REFERENCES users(id);

ALTER TABLE task DROP COLUMN created_by;
ALTER TABLE task DROP COLUMN updated_by;
ALTER TABLE task RENAME COLUMN created_by_id TO created_by;
ALTER TABLE task RENAME COLUMN updated_by_id TO updated_by;

-- 9. TASK_ASSIGNEE table - change user_id from UUID to integer
ALTER TABLE task_assignee ADD COLUMN user_id_int INTEGER;

UPDATE task_assignee SET user_id_int = (SELECT id FROM users WHERE user_uid = task_assignee.user_id);

-- Drop old constraint, rename column, and add new constraint
ALTER TABLE task_assignee DROP COLUMN user_id;
ALTER TABLE task_assignee RENAME COLUMN user_id_int TO user_id;
ALTER TABLE task_assignee ADD CONSTRAINT task_assignee_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id);

-- 10. PROJECT_MEMBER table - change user_id from UUID to integer
ALTER TABLE project_member ADD COLUMN user_id_int INTEGER;

UPDATE project_member SET user_id_int = (SELECT id FROM users WHERE user_uid = project_member.user_id);

ALTER TABLE project_member DROP COLUMN user_id;
ALTER TABLE project_member RENAME COLUMN user_id_int TO user_id;
ALTER TABLE project_member ADD CONSTRAINT project_member_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id);

-- 11. WORKSPACE table - change created_by and updated_by
ALTER TABLE workspace ADD COLUMN created_by_id INTEGER;
ALTER TABLE workspace ADD COLUMN updated_by_id INTEGER;

UPDATE workspace SET created_by_id = (SELECT id FROM users WHERE user_uid = workspace.created_by);
UPDATE workspace SET updated_by_id = (SELECT id FROM users WHERE user_uid = workspace.updated_by);

ALTER TABLE workspace DROP CONSTRAINT IF EXISTS workspace_created_by_fkey;
ALTER TABLE workspace DROP CONSTRAINT IF EXISTS workspace_updated_by_fkey;
ALTER TABLE workspace ADD CONSTRAINT workspace_created_by_id_fkey FOREIGN KEY (created_by_id) REFERENCES users(id);
ALTER TABLE workspace ADD CONSTRAINT workspace_updated_by_id_fkey FOREIGN KEY (updated_by_id) REFERENCES users(id);

ALTER TABLE workspace DROP COLUMN created_by;
ALTER TABLE workspace DROP COLUMN updated_by;
ALTER TABLE workspace RENAME COLUMN created_by_id TO created_by;
ALTER TABLE workspace RENAME COLUMN updated_by_id TO updated_by;

-- 12. WORKSPACE_MEMBER table - change user_id from UUID to integer (if it exists)
-- Let's check if this table has a user_id column that needs to be updated
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_name = 'workspace_member' AND column_name = 'user_id' 
               AND data_type = 'uuid') THEN
        
        ALTER TABLE workspace_member ADD COLUMN user_id_int INTEGER;
        UPDATE workspace_member SET user_id_int = (SELECT id FROM users WHERE user_uid = workspace_member.user_id);
        ALTER TABLE workspace_member DROP COLUMN user_id;
        ALTER TABLE workspace_member RENAME COLUMN user_id_int TO user_id;
        ALTER TABLE workspace_member ADD CONSTRAINT workspace_member_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id);
    END IF;
END $$;

COMMIT;