package models

import (
	"time"

	"github.com/google/uuid"
)

// Database models - these map directly to the database schema
// These use integer IDs internally for efficiency
// Now supporting both pgx (`db` tags) and GORM (`gorm` tags)

type User struct {
	ID        int       `db:"id" gorm:"primaryKey;column:id"`
	UserUID   uuid.UUID `db:"user_uid" gorm:"type:uuid;default:gen_random_uuid();column:user_uid;uniqueIndex"`
	Email     string    `db:"email" gorm:"column:email;uniqueIndex;not null"`
	Password  string    `db:"password_hash" gorm:"column:password_hash;not null"`
	Name      string    `db:"name" gorm:"column:name;not null"`
	CreatedAt time.Time `db:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `db:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	IsActive  bool      `db:"is_active" gorm:"column:is_active;default:true"`
}

type Workspace struct {
	ID           int        `db:"id" gorm:"primaryKey;column:id"`
	WorkspaceUID uuid.UUID  `db:"workspace_uid" gorm:"type:uuid;default:gen_random_uuid();column:workspace_uid;uniqueIndex"`
	Name         string     `db:"name" gorm:"column:name;not null"`
	Description  *string    `db:"description" gorm:"column:description;type:text"`
	CreatedAt    time.Time  `db:"created_at" gorm:"column:created_at;autoCreateTime"`
	CreatedBy    *int       `db:"created_by" gorm:"column:created_by"`
	UpdatedAt    *time.Time `db:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	UpdatedBy    *int       `db:"updated_by" gorm:"column:updated_by"`
	IsActive     bool       `db:"is_active" gorm:"column:is_active;default:true"`
}

type Project struct {
	ID               int        `db:"id" gorm:"primaryKey;column:id"`
	ProjectUID       uuid.UUID  `db:"project_uid" gorm:"type:uuid;default:gen_random_uuid();column:project_uid;uniqueIndex"`
	UserID           int        `db:"user_id" gorm:"column:user_id;not null;index"`
	Name             string     `db:"name" gorm:"column:name;not null"`
	Description      *string    `db:"description" gorm:"column:description;type:text"`
	Status           string     `db:"status" gorm:"column:status;default:'active'"`
	Color            string     `db:"color" gorm:"column:color;default:'#FFFFFF'"`
	Position         *int       `db:"position" gorm:"column:position"`
	StartDate        *time.Time `db:"start_date" gorm:"column:start_date"`
	EndDate          *time.Time `db:"end_date" gorm:"column:end_date"`
	IsPrivate        bool       `db:"is_private" gorm:"column:is_private;default:false;index"`
	DbmlContent      *string    `db:"dbml_content" gorm:"column:dbml_content;type:text"`
	DbmlLayoutData   *string    `db:"dbml_layout_data" gorm:"column:dbml_layout_data;type:jsonb"`
	FlowchartContent *string    `db:"flowchart_content" gorm:"column:flowchart_content;type:text"`
	CreatedAt        time.Time  `db:"created_at" gorm:"column:created_at;autoCreateTime"`
	CreatedBy        *int       `db:"created_by" gorm:"column:created_by"`
	UpdatedAt        *time.Time `db:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	UpdatedBy        *int       `db:"updated_by" gorm:"column:updated_by"`
	IsActive         bool       `db:"is_active" gorm:"column:is_active;default:true;index"`
}

type ProjectMember struct {
	ID        int       `db:"id" gorm:"primaryKey;column:id"`
	ProjectID int       `db:"project_id" gorm:"column:project_id;not null"`
	UserID    int       `db:"user_id" gorm:"column:user_id;not null"`
	Role      string    `db:"role" gorm:"column:role;default:'member'"` // 'owner' or 'member'
	JoinedAt  time.Time `db:"joined_at" gorm:"column:joined_at;autoCreateTime"`
}

type List struct {
	ID        int        `db:"id" gorm:"primaryKey;column:id"`
	ListUID   uuid.UUID  `db:"list_uid" gorm:"type:uuid;default:gen_random_uuid();column:list_uid;uniqueIndex"`
	ProjectID int        `db:"project_id" gorm:"column:project_id;not null;index"`
	Name      string     `db:"name" gorm:"column:name;not null"`
	Color     string     `db:"color" gorm:"column:color;default:'#FFFFFF'"`
	Position  int        `db:"position" gorm:"column:position;index"`
	CreatedAt time.Time  `db:"created_at" gorm:"column:created_at;autoCreateTime"`
	CreatedBy *int       `db:"created_by" gorm:"column:created_by"`
	UpdatedAt *time.Time `db:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	UpdatedBy *int       `db:"updated_by" gorm:"column:updated_by"`
	IsActive  bool       `db:"is_active" gorm:"column:is_active;default:true"`
}

type Task struct {
	ID          int        `db:"id" gorm:"primaryKey;column:id"`
	TaskUID     uuid.UUID  `db:"task_uid" gorm:"type:uuid;default:gen_random_uuid();column:task_uid;uniqueIndex"`
	ListID      int        `db:"list_id" gorm:"column:list_id;not null;index"`
	Title       string     `db:"title" gorm:"column:title;not null"`
	Description *string    `db:"description" gorm:"column:description;type:text"`
	Priority    *string    `db:"priority" gorm:"column:priority"`
	Status      string     `db:"status" gorm:"column:status;default:'todo'"`
	Color       string     `db:"color" gorm:"column:color;default:'#FFFFFF'"`
	Position    *int       `db:"position" gorm:"column:position"`
	IsCompleted bool       `db:"is_completed" gorm:"column:is_completed;default:false"`
	DueDate     *time.Time `db:"due_date" gorm:"column:due_date"`
	CompletedAt *time.Time `db:"completed_at" gorm:"column:completed_at"`
	CreatedAt   time.Time  `db:"created_at" gorm:"column:created_at;autoCreateTime"`
	CreatedBy   *int       `db:"created_by" gorm:"column:created_by"`
	UpdatedAt   *time.Time `db:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	UpdatedBy   *int       `db:"updated_by" gorm:"column:updated_by"`
	IsActive    bool       `db:"is_active" gorm:"column:is_active;default:true;index"`
}

// TaskComment - Comments on tasks
type TaskComment struct {
	ID         int        `db:"id" gorm:"primaryKey;column:id"`
	CommentUID uuid.UUID  `db:"comment_uid" gorm:"type:uuid;default:gen_random_uuid();column:comment_uid;uniqueIndex"`
	TaskID     int        `db:"task_id" gorm:"column:task_id;not null;index"`
	UserID     int        `db:"user_id" gorm:"column:user_id;not null"`
	Content    string     `db:"content" gorm:"column:content;type:text;not null"`
	CreatedAt  time.Time  `db:"created_at" gorm:"column:created_at;autoCreateTime;index"`
	UpdatedAt  *time.Time `db:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	IsActive   bool       `db:"is_active" gorm:"column:is_active;default:true"`

	// GORM relationships
	Task *Task `gorm:"foreignKey:TaskID;references:ID"`
	User *User `gorm:"foreignKey:UserID;references:ID"`
}

// TableName specifies the table name for TaskComment
func (TaskComment) TableName() string {
	return "task_comment"
}

// TaskCategory - Categories for tasks (project-specific)
type TaskCategory struct {
	ID          int        `db:"id" gorm:"primaryKey;column:id"`
	CategoryUID uuid.UUID  `db:"category_uid" gorm:"type:uuid;default:gen_random_uuid();column:category_uid;uniqueIndex"`
	ProjectID   int        `db:"project_id" gorm:"column:project_id;not null;index"`
	Name        string     `db:"name" gorm:"column:name;not null"`
	Color       string     `db:"color" gorm:"column:color;default:'#808080'"`
	Description *string    `db:"description" gorm:"column:description;type:text"`
	CreatedAt   time.Time  `db:"created_at" gorm:"column:created_at;autoCreateTime"`
	CreatedBy   *int       `db:"created_by" gorm:"column:created_by"`
	UpdatedAt   *time.Time `db:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	UpdatedBy   *int       `db:"updated_by" gorm:"column:updated_by"`
	IsActive    bool       `db:"is_active" gorm:"column:is_active;default:true"`

	// GORM relationships
	Project *Project `gorm:"foreignKey:ProjectID;references:ID"`
}

// TableName specifies the table name for TaskCategory
func (TaskCategory) TableName() string {
	return "task_category"
}

// TaskCategoryMap - Many-to-many relationship between tasks and categories
type TaskCategoryMap struct {
	ID         int       `db:"id" gorm:"primaryKey;column:id"`
	TaskID     int       `db:"task_id" gorm:"column:task_id;not null;index"`
	CategoryID int       `db:"category_id" gorm:"column:category_id;not null;index"`
	AssignedAt time.Time `db:"assigned_at" gorm:"column:assigned_at;autoCreateTime"`
	AssignedBy *int      `db:"assigned_by" gorm:"column:assigned_by"`

	// GORM relationships
	Task     *Task         `gorm:"foreignKey:TaskID;references:ID"`
	Category *TaskCategory `gorm:"foreignKey:CategoryID;references:ID"`
}

// TableName specifies the table name for TaskCategoryMap
func (TaskCategoryMap) TableName() string {
	return "task_category_map"
}

// Phase 3: Brainstorming & Planning Layer Models

type BrainstormCanvas struct {
	ID        int        `db:"id" gorm:"primaryKey;column:id"`
	CanvasUID uuid.UUID  `db:"canvas_uid" gorm:"type:uuid;default:gen_random_uuid();column:canvas_uid;uniqueIndex"`
	ProjectID int        `db:"project_id" gorm:"column:project_id;not null;index"`
	StateJSON string     `db:"state_json" gorm:"column:state_json;type:jsonb"` // JSONB stored as string
	CreatedAt time.Time  `db:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt *time.Time `db:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	CreatedBy *int       `db:"created_by" gorm:"column:created_by"`
	UpdatedBy *int       `db:"updated_by" gorm:"column:updated_by"`
	IsActive  bool       `db:"is_active" gorm:"column:is_active;default:true"`
}

type NoteFolder struct {
	ID             int        `db:"id" gorm:"primaryKey;column:id"`
	FolderUID      uuid.UUID  `db:"folder_uid" gorm:"type:uuid;default:gen_random_uuid();column:folder_uid;uniqueIndex"`
	ProjectID      int        `db:"project_id" gorm:"column:project_id;not null;index"`
	ParentFolderID *int       `db:"parent_folder_id" gorm:"column:parent_folder_id;index"`
	Name           string     `db:"name" gorm:"column:name;not null"`
	Position       *int       `db:"position" gorm:"column:position"`
	CreatedAt      time.Time  `db:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      *time.Time `db:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	CreatedBy      *int       `db:"created_by" gorm:"column:created_by"`
	UpdatedBy      *int       `db:"updated_by" gorm:"column:updated_by"`
	IsActive       bool       `db:"is_active" gorm:"column:is_active;default:true"`
}

type Note struct {
	ID          int        `db:"id" gorm:"primaryKey;column:id"`
	NoteUID     uuid.UUID  `db:"note_uid" gorm:"type:uuid;default:gen_random_uuid();column:note_uid;uniqueIndex"`
	ProjectID   int        `db:"project_id" gorm:"column:project_id;not null;index"`
	FolderID    *int       `db:"folder_id" gorm:"column:folder_id;index"`
	Title       string     `db:"title" gorm:"column:title;not null"`
	ContentJSON string     `db:"content_json" gorm:"column:content_json;type:jsonb"` // JSONB stored as string for rich content
	Position    *int       `db:"position" gorm:"column:position;index"`
	CreatedAt   time.Time  `db:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   *time.Time `db:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	CreatedBy   *int       `db:"created_by" gorm:"column:created_by"`
	UpdatedBy   *int       `db:"updated_by" gorm:"column:updated_by"`
	IsActive    bool       `db:"is_active" gorm:"column:is_active;default:true"`
}

// Chat Conversation models for DevSprint-AI
type ChatConversation struct {
	ID              int        `db:"id" gorm:"primaryKey;column:id"`
	ConversationUID uuid.UUID  `db:"conversation_uid" gorm:"type:uuid;default:gen_random_uuid();column:conversation_uid;uniqueIndex"`
	ProjectID       int        `db:"project_id" gorm:"column:project_id;not null;index"`
	Name            string     `db:"name" gorm:"column:name;not null"`
	CreatedAt       time.Time  `db:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       *time.Time `db:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	CreatedBy       *int       `db:"created_by" gorm:"column:created_by"`
	UpdatedBy       *int       `db:"updated_by" gorm:"column:updated_by"`
	IsActive        bool       `db:"is_active" gorm:"column:is_active;default:true"`
}

type ChatMessage struct {
	ID             int       `db:"id" gorm:"primaryKey;column:id"`
	MessageUID     uuid.UUID `db:"message_uid" gorm:"type:uuid;default:gen_random_uuid();column:message_uid;uniqueIndex"`
	ConversationID int       `db:"conversation_id" gorm:"column:conversation_id;not null;index"`
	MessageType    string    `db:"message_type" gorm:"column:message_type;not null"` // 'user' or 'ai'
	Content        string    `db:"content" gorm:"column:content;type:text;not null"`
	CreatedAt      time.Time `db:"created_at" gorm:"column:created_at;autoCreateTime;index"`
	CreatedBy      *int      `db:"created_by" gorm:"column:created_by"`
}

type UserSettings struct {
	ID                   int       `db:"id" gorm:"primaryKey;column:id"`
	SettingsUID          uuid.UUID `db:"settings_uid" gorm:"type:uuid;default:gen_random_uuid();column:settings_uid;uniqueIndex"`
	UserID               int       `db:"user_id" gorm:"column:user_id;not null;uniqueIndex"`
	Theme                string    `db:"theme" gorm:"column:theme;default:'projectnest-default';index"` // theme name
	Language             string    `db:"language" gorm:"column:language;default:'en'"`                  // preferred language
	Timezone             string    `db:"timezone" gorm:"column:timezone;default:'UTC'"`                 // user timezone
	NotificationsEnabled bool      `db:"notifications_enabled" gorm:"column:notifications_enabled;default:true"`
	EmailNotifications   bool      `db:"email_notifications" gorm:"column:email_notifications;default:true"`
	SoundEnabled         bool      `db:"sound_enabled" gorm:"column:sound_enabled;default:true"`
	CompactMode          bool      `db:"compact_mode" gorm:"column:compact_mode;default:false"`          // compact UI mode
	AutoSave             bool      `db:"auto_save" gorm:"column:auto_save;default:true"`                 // auto-save feature
	AutoSaveInterval     int       `db:"auto_save_interval" gorm:"column:auto_save_interval;default:30"` // in seconds
	CreatedAt            time.Time `db:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt            time.Time `db:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// Additional models for complete database coverage
type TaskAssignee struct {
	ID         int       `db:"id" gorm:"primaryKey;column:id"`
	TaskID     int       `db:"task_id" gorm:"column:task_id;not null;index"`
	UserID     int       `db:"user_id" gorm:"column:user_id;not null;index"`
	AssignedAt time.Time `db:"assigned_at" gorm:"column:assigned_at;autoCreateTime"`
}

type WorkspaceMember struct {
	ID          int       `db:"id" gorm:"primaryKey;column:id"`
	WorkspaceID int       `db:"workspace_id" gorm:"column:workspace_id;not null;index"`
	UserID      int       `db:"user_id" gorm:"column:user_id;not null;index"`
	Role        string    `db:"role" gorm:"column:role;default:'member'"`
	JoinedAt    time.Time `db:"joined_at" gorm:"column:joined_at;autoCreateTime"`
}

// Notes model (legacy table that may need to be cleaned up)
type Notes struct {
	ID        int        `db:"id" gorm:"primaryKey;column:id"`
	NoteUID   uuid.UUID  `db:"note_uid" gorm:"type:uuid;default:gen_random_uuid();column:note_uid;uniqueIndex"`
	ProjectID int        `db:"project_id" gorm:"column:project_id;not null;index"`
	FolderID  *int       `db:"folder_id" gorm:"column:folder_id;index"`
	Title     string     `db:"title" gorm:"column:title;not null"`
	Content   string     `db:"content" gorm:"column:content;type:text;default:''"`
	CreatedAt time.Time  `db:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt *time.Time `db:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	CreatedBy *int       `db:"created_by" gorm:"column:created_by"`
	UpdatedBy *int       `db:"updated_by" gorm:"column:updated_by"`
	IsActive  bool       `db:"is_active" gorm:"column:is_active;default:true"`
}

// Canvas model
type Canvas struct {
	ID         int        `db:"id" gorm:"primaryKey;column:id"`
	CanvasUID  uuid.UUID  `db:"canvas_uid" gorm:"type:uuid;default:gen_random_uuid();column:canvas_uid;uniqueIndex"`
	ProjectID  int        `db:"project_id" gorm:"column:project_id;not null;index"`
	CanvasData string     `db:"canvas_data" gorm:"column:canvas_data;type:jsonb"`
	CreatedAt  time.Time  `db:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  *time.Time `db:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	CreatedBy  *int       `db:"created_by" gorm:"column:created_by"`
	UpdatedBy  *int       `db:"updated_by" gorm:"column:updated_by"`
	IsActive   bool       `db:"is_active" gorm:"column:is_active;default:true"`
}
