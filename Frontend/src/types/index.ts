// API Types matching the database schema

export interface Workspace {
  id: number;
  workspace_uid: string;
  name: string;
  description?: string;
  created_at: string;
  created_by: string;
  updated_at?: string;
  updated_by?: string;
  is_active: boolean;
}

export interface Project {
  id: number;
  project_uid: string;
  name: string;
  description?: string;
  status: string;
  color: string;
  position?: number;
  start_date?: string;
  end_date?: string;
  is_private: boolean;
  dbml_content?: string;
  dbml_layout_data?: string;
  flowchart_content?: string;
  created_at: string;
  created_by: string;
  updated_at?: string;
  updated_by?: string;
  is_active: boolean;
}

export interface List {
  id: number;
  list_uid: string;
  project_id: number;
  name: string;
  color: string;
  position: number;
  created_at: string;
  created_by: string;
  updated_at?: string;
  updated_by?: string;
  is_active: boolean;
}

export interface Task {
  id: number;
  task_uid: string;
  list_id: number;
  title: string;
  description?: string;
  priority?: string;
  status: string;
  color: string;
  position?: number;
  is_completed: boolean;
  due_date?: string;
  completed_at?: string;
  created_at: string;
  created_by: string;
  updated_at?: string;
  updated_by?: string;
  is_active: boolean;
  assignees?: TaskAssigneeWithUser[];
}

export interface TaskAssignee {
  id: number;
  task_id: number;
  user_id: string;
  assigned_at: string;
}

export interface TaskAssigneeWithUser {
  user_uid: string;
  name: string;
  email: string;
  assigned_at: string;
}

export interface WorkspaceMember {
  id: number;
  workspace_id: number;
  user_id: string;
  role: string;
  joined_at: string;
}

export interface ProjectMember {
  user_uid: string;
  email: string;
  name: string;
  role: string;
  joined_at: string;
}

// Frontend UI Types
export interface ProjectWithLists extends Project {
  lists: ListWithTasks[] | null;
}

export interface ListWithTasks extends List {
  tasks: Task[] | null;
}

// API Response Types
export interface ApiResponse<T> {
  data: T;
  success: boolean;
  message?: string;
}

export interface ApiError {
  error: string;
  message: string;
  statusCode: number;
}

// Request Types for API calls
export interface ProjectRequest {
  name: string;
  description?: string;
  status?: string;
  color?: string;
  position?: number;
  start_date?: string;
  end_date?: string;
  is_private?: boolean;
  dbml_content?: string;
  dbml_layout_data?: string;
  flowchart_content?: string;
}

export interface ProjectUpdateRequest {
  name?: string;
  description?: string;
  status?: string;
  color?: string;
  position?: number;
  start_date?: string;
  end_date?: string;
  is_private?: boolean;
  dbml_content?: string;
  dbml_layout_data?: string;
  flowchart_content?: string;
}

export interface ListRequest {
  project_uid: string;
  name: string;
  color?: string;
  position?: number;
}

export interface ListUpdateRequest {
  name?: string;
  color?: string;
  position?: number;
}

export interface TaskRequest {
  list_uid: string;
  title: string;
  description?: string;
  priority?: string;
  status?: string;
  color?: string;
  position?: number;
  is_completed?: boolean;
  due_date?: string;
}

export interface TaskUpdateRequest {
  title?: string;
  description?: string;
  priority?: string;
  status?: string;
  color?: string;
  position?: number;
  is_completed?: boolean;
  due_date?: string;
}

// Progress tracking types for Phase 2
export interface ProjectProgress {
  total_tasks: number;
  completed_tasks: number;
  todo_tasks: number;
  progress: number;
}

export interface ProjectWithProgress extends Project {
  task_stats: ProjectProgress;
}

// Color constants
export const COLORS = {
  WHITE: '#FFFFFF',
  RED: '#EF4444',
  BLUE: '#3B82F6',
  GREEN: '#10B981',
  YELLOW: '#F59E0B',
  ORANGE: '#F97316',
  PURPLE: '#8B5CF6',
} as const;

export const COLOR_OPTIONS = [
  { name: 'White', value: COLORS.WHITE },
  { name: 'Red', value: COLORS.RED },
  { name: 'Blue', value: COLORS.BLUE },
  { name: 'Green', value: COLORS.GREEN },
  { name: 'Yellow', value: COLORS.YELLOW },
  { name: 'Orange', value: COLORS.ORANGE },
  { name: 'Purple', value: COLORS.PURPLE },
] as const;

// Phase 3: Brainstorming & Planning Layer Types

export interface BrainstormCanvas {
  canvas_uid: string;
  project_id: number;
  state_json: string;
  created_at: string;
  updated_at?: string;
}

export interface Note {
  note_uid: string;
  project_id: number;
  title: string;
  content: NoteContent;
  folder_id?: number;
  position?: number;
  created_at: string;
  updated_at?: string;
}

export interface NotesResponse {
  notes: Note[];
  total: number;
}

// Folder types
export interface NoteFolder {
  id: number;
  folder_uid: string;
  project_id: number;
  parent_folder_id?: number;
  name: string;
  position?: number;
  created_at: string;
  updated_at?: string;
}

export interface NoteFoldersResponse {
  folders: NoteFolder[];
  total: number;
}

export interface NoteFolderRequest {
  name: string;
  parent_folder_id?: number;
  position?: number;
}

// Phase 3 Request types
export interface CanvasRequest {
  state_json: string;
}

// Note content structures matching backend
export interface NoteChecklistItem {
  id: string;
  text: string;
  completed: boolean;
}

export interface NoteBlockMetadata {
  level?: number;
  style?: string;
}

export interface NoteBlock {
  id: string;
  type: string; // "text", "checklist", "heading"
  content?: string;
  metadata?: NoteBlockMetadata;
  items?: NoteChecklistItem[];
  children?: NoteBlock[];
}

export interface NoteContent {
  blocks: NoteBlock[];
}

export interface NoteRequest {
  title: string;
  content: NoteContent;
  folder_id?: number;
  position?: number;
}

export interface NoteUpdateRequest {
  title?: string;
  content?: NoteContent;
  folder_id?: number;
  position?: number;
}

// DEV AI Chat Types
export interface ChatMessage {
  message_uid: string;
  conversation_uid: string;
  message_type: 'user' | 'ai';
  content: string;
  created_at: string;
  isLoading?: boolean; // Frontend only for UI state
}

export interface ChatConversation {
  conversation_uid: string;
  project_uid: string;
  name: string;
  created_at: string;
  updated_at?: string;
}

export interface ChatConversationWithMessages extends ChatConversation {
  messages: ChatMessage[];
}

export interface ChatConversationRequest {
  name: string;
}

export interface ChatMessageRequest {
  conversation_uid: string;
  content: string;
  message_type: 'user' | 'ai';
}

// AI Project Creation Types
export interface AIListRequest {
  name: string;
  description: string;
  position: number;
}

export interface AITaskRequest {
  title: string;
  description: string;
  list_name: string;
  priority: string;
  status: string;
  position: number;
}

// AI Project Creation Types
export interface AIProjectCreationRequest {
  project_content: string;
}

export interface AIProjectCreationResponse {
  message: string;
  project: Project;
}

// Task Comments Types
export interface TaskComment {
  comment_uid: string;
  task_id: number;
  user_id: number;
  content: string;
  created_at: string;
  updated_at?: string;
  user_name?: string; // May be included in response
}

export interface TaskCommentRequest {
  task_uid: string;
  content: string;
}

export interface TaskCommentUpdateRequest {
  content: string;
}

// Task Categories Types
export interface TaskCategory {
  category_uid: string;
  project_id: number;
  name: string;
  color: string;
  description?: string;
  created_at: string;
  updated_at?: string;
}

export interface TaskCategoryRequest {
  project_uid: string;
  name: string;
  color?: string;
  description?: string;
}

export interface TaskCategoryUpdateRequest {
  name?: string;
  color?: string;
  description?: string;
}