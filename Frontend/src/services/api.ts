// API service layer for backend integration

import { 
  Project, List, Task, ProjectWithLists, ApiResponse,
  ProjectRequest, ProjectUpdateRequest,
  ListRequest, ListUpdateRequest,
  TaskRequest, TaskUpdateRequest,
  ProjectProgress, ProjectWithProgress,
  ProjectMember, TaskAssigneeWithUser, UserSearchResult,
  // Phase 3 types
  BrainstormCanvas, Note, NotesResponse,
  CanvasRequest, NoteRequest, NoteUpdateRequest,
  NoteFolder, NoteFoldersResponse, NoteFolderRequest,
  // Chat types
  ChatConversation, ChatConversationWithMessages, ChatConversationRequest,
  ChatMessage, ChatMessageRequest,
  // AI types
  AIProjectCreationRequest, AIProjectCreationResponse,
  // Task Comments and Categories
  TaskComment, TaskCommentRequest, TaskCommentUpdateRequest,
  TaskCategory, TaskCategoryRequest, TaskCategoryUpdateRequest
} from '@/types';

const BACKEND_PORT = import.meta.env.VITE_BACKEND_PORT || '8080';
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || `http://localhost:${BACKEND_PORT}/api`;

class ApiService {
  private log(level: 'info' | 'error' | 'warn', message: string, data?: any) {
    const timestamp = new Date().toISOString();
    const prefix = `[API Service] ${timestamp} [${level.toUpperCase()}]`;
    
    if (level === 'error') {
      console.error(prefix, message, data || '');
    } else if (level === 'warn') {
      console.warn(prefix, message, data || '');
    } else {
      console.log(prefix, message, data || '');
    }
  }

  private async request<T>(endpoint: string, options?: RequestInit): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;
    const method = options?.method || 'GET';
    
    this.log('info', `Making ${method} request to ${url}`, {
      headers: options?.headers,
      body: options?.body ? JSON.parse(options.body as string) : undefined
    });

    // Get token from localStorage for authentication
    const token = localStorage.getItem('token');
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...(options?.headers as Record<string, string>),
    };

    // Add authorization header if token exists and not an auth endpoint
    if (token && !endpoint.startsWith('/auth')) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    try {
      const response = await fetch(url, {
        headers,
        ...options,
      });

      this.log('info', `Response received from ${url}`, {
        status: response.status,
        statusText: response.statusText,
        ok: response.ok
      });

      if (!response.ok) {
        const errorText = await response.text();
        this.log('error', `API Error: ${response.status} ${response.statusText}`, {
          url,
          method,
          errorBody: errorText
        });

        // Handle authentication failures
        if (response.status === 401) {
          this.log('warn', 'Authentication failed - clearing tokens and redirecting to login');
          
          // Clear authentication data
          localStorage.removeItem('token');
          localStorage.removeItem('user');
          
          // Show toast message
          if (typeof window !== 'undefined' && window.location.pathname !== '/login') {
            // Import toast dynamically to avoid circular dependencies
            import('@/hooks/use-toast').then(({ toast }) => {
              toast({
                title: "Session Expired",
                description: "Your session has expired. Please login again.",
                variant: "destructive",
              });
            });
            
            // Redirect to login page
            window.location.href = '/login';
          }
          
          throw new Error(`Authentication failed: ${errorText}`);
        }

        throw new Error(`API Error: ${response.status} ${response.statusText} - ${errorText}`);
      }

      const data = await response.json();
      this.log('info', `Successful response from ${url}`, data);
      return data;
    } catch (error) {
      this.log('error', `Request failed to ${url}`, {
        error: error instanceof Error ? error.message : String(error),
        method,
        body: options?.body
      });
      throw error;
    }
  }

  // Projects
  async getProjects(): Promise<ApiResponse<Project[]>> {
    return this.request('/projects');
  }

  async getProjectsWithProgress(): Promise<ApiResponse<ProjectWithProgress[]>> {
    return this.request('/projects?include_progress=true');
  }

  async getProjectProgress(id: string): Promise<ApiResponse<ProjectProgress>> {
    return this.request(`/projects/${id}/progress`);
  }

  async getProject(id: string): Promise<ApiResponse<ProjectWithLists>> {
    return this.request(`/projects/${id}`);
  }

  async createProject(project: ProjectRequest): Promise<ApiResponse<Project>> {
    return this.request('/projects', {
      method: 'POST',
      body: JSON.stringify(project),
    });
  }

  async updateProject(id: string, project: ProjectRequest): Promise<ApiResponse<Project>> {
    return this.request(`/projects/${id}`, {
      method: 'PUT',
      body: JSON.stringify(project),
    });
  }

  async partialUpdateProject(id: string, updates: ProjectUpdateRequest): Promise<ApiResponse<Project>> {
    return this.request(`/projects/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(updates),
    });
  }

  async deleteProject(id: string): Promise<ApiResponse<void>> {
    return this.request(`/projects/${id}`, {
      method: 'DELETE',
    });
  }

  // Project Members
  async getProjectMembers(projectId: string): Promise<ApiResponse<ProjectMember[]>> {
    return this.request(`/projects/${projectId}/members`);
  }

  async addProjectMember(projectId: string, email: string, role: string = 'member', canWrite: boolean = false): Promise<ApiResponse<any>> {
    return this.request(`/projects/${projectId}/members`, {
      method: 'POST',
      body: JSON.stringify({ email, role, can_write: canWrite }),
    });
  }

  async updateMemberWriteAccess(projectId: string, memberUid: string, canWrite: boolean): Promise<ApiResponse<any>> {
    return this.request(`/projects/${projectId}/members/${memberUid}/write-access`, {
      method: 'PUT',
      body: JSON.stringify({ can_write: canWrite }),
    });
  }

  async searchUsers(query: string, projectId?: string): Promise<ApiResponse<UserSearchResult[]>> {
    const params = new URLSearchParams();
    if (query) params.set('q', query);
    if (projectId) params.set('project_uid', projectId);
    return this.request(`/users/search?${params.toString()}`);
  }

  // Lists
  async createList(listData: ListRequest): Promise<ApiResponse<List>> {
    return this.request('/lists', {
      method: 'POST',
      body: JSON.stringify(listData),
    });
  }

  async updateList(uid: string, listData: ListRequest): Promise<ApiResponse<List>> {
    return this.request(`/lists/${uid}`, {
      method: 'PUT',
      body: JSON.stringify(listData),
    });
  }

  async partialUpdateList(uid: string, updates: ListUpdateRequest): Promise<ApiResponse<List>> {
    return this.request(`/lists/${uid}`, {
      method: 'PATCH',
      body: JSON.stringify(updates),
    });
  }

  async deleteList(uid: string): Promise<ApiResponse<void>> {
    return this.request(`/lists/${uid}`, {
      method: 'DELETE',
    });
  }

  async updateListPosition(uid: string, position: number): Promise<ApiResponse<List>> {
    const requestBody = { position };
    this.log('info', `Updating list position for ${uid} to position ${position}`, requestBody);
    this.log('info', `Request body type check:`, { 
      position, 
      positionType: typeof position, 
      positionIsNumber: typeof position === 'number',
      jsonString: JSON.stringify(requestBody)
    });
    
    return this.request(`/lists/${uid}/position`, {
      method: 'PUT',
      body: JSON.stringify(requestBody),
    });
  }

  // Tasks
  async createTask(taskData: TaskRequest): Promise<ApiResponse<Task>> {
    return this.request('/tasks', {
      method: 'POST',
      body: JSON.stringify(taskData),
    });
  }

  async updateTask(uid: string, taskData: TaskRequest): Promise<ApiResponse<Task>> {
    return this.request(`/tasks/${uid}`, {
      method: 'PUT',
      body: JSON.stringify(taskData),
    });
  }

  async partialUpdateTask(uid: string, updates: TaskUpdateRequest): Promise<ApiResponse<Task>> {
    return this.request(`/tasks/${uid}`, {
      method: 'PATCH',
      body: JSON.stringify(updates),
    });
  }

  async deleteTask(uid: string): Promise<ApiResponse<void>> {
    return this.request(`/tasks/${uid}`, {
      method: 'DELETE',
    });
  }

  async moveTask(taskUid: string, listUid: string): Promise<ApiResponse<any>> {
    return this.request(`/tasks/${taskUid}/move`, {
      method: 'POST',
      body: JSON.stringify({ list_uid: listUid }),
    });
  }

  // Task Assignees
  async getTaskAssignees(taskUid: string): Promise<ApiResponse<TaskAssigneeWithUser[]>> {
    return this.request(`/tasks/${taskUid}/assignees`);
  }

  async assignTaskToUser(taskUid: string, userUid: string): Promise<ApiResponse<void>> {
    return this.request(`/tasks/${taskUid}/assignees`, {
      method: 'POST',
      body: JSON.stringify({ user_uid: userUid }),
    });
  }

  async unassignTaskFromUser(taskUid: string, userUid: string): Promise<ApiResponse<void>> {
    return this.request(`/tasks/${taskUid}/assignees/${userUid}`, {
      method: 'DELETE',
    });
  }

  async bulkAssignTask(taskUid: string, userUids: string[]): Promise<ApiResponse<void>> {
    return this.request(`/tasks/${taskUid}/assignees/bulk`, {
      method: 'POST',
      body: JSON.stringify({ user_uids: userUids }),
    });
  }

  // Phase 3: Canvas API methods
  async getCanvas(projectUid: string): Promise<ApiResponse<BrainstormCanvas>> {
    return this.request(`/projects/${projectUid}/canvas`);
  }

  async updateCanvas(projectUid: string, canvasData: CanvasRequest): Promise<ApiResponse<any>> {
    return this.request(`/projects/${projectUid}/canvas`, {
      method: 'POST',
      body: JSON.stringify(canvasData),
    });
  }

  async deleteCanvas(projectUid: string): Promise<ApiResponse<any>> {
    return this.request(`/projects/${projectUid}/canvas`, {
      method: 'DELETE',
    });
  }

  // Phase 3: Notes API methods
  async getNotesByProject(projectUid: string): Promise<ApiResponse<NotesResponse>> {
    return this.request(`/projects/${projectUid}/notes`);
  }

  async getNote(noteUid: string): Promise<ApiResponse<Note>> {
    return this.request(`/notes/${noteUid}`);
  }

  async createNote(projectUid: string, noteData: NoteRequest): Promise<ApiResponse<Note>> {
    return this.request(`/projects/${projectUid}/notes`, {
      method: 'POST',
      body: JSON.stringify(noteData),
    });
  }

  async updateNote(noteUid: string, noteData: NoteRequest): Promise<ApiResponse<Note>> {
    return this.request(`/notes/${noteUid}`, {
      method: 'PUT',
      body: JSON.stringify(noteData),
    });
  }

  async partialUpdateNote(noteUid: string, noteData: NoteUpdateRequest): Promise<ApiResponse<Note>> {
    return this.request(`/notes/${noteUid}`, {
      method: 'PATCH',
      body: JSON.stringify(noteData),
    });
  }

  async deleteNote(noteUid: string): Promise<ApiResponse<any>> {
    return this.request(`/notes/${noteUid}`, {
      method: 'DELETE',
    });
  }

    // Folder operations
  async getFolders(projectUid: string): Promise<ApiResponse<NoteFoldersResponse>> {
    return this.request(`/projects/${projectUid}/folders`);
  }

  async createFolder(projectUid: string, folderData: NoteFolderRequest): Promise<ApiResponse<NoteFolder>> {
    return this.request(`/projects/${projectUid}/folders`, {
      method: 'POST',
      body: JSON.stringify(folderData),
    });
  }

  async updateFolder(folderUid: string, folderData: NoteFolderRequest): Promise<ApiResponse<NoteFolder>> {
    return this.request(`/folders/${folderUid}`, {
      method: 'PUT',
      body: JSON.stringify(folderData),
    });
  }

  async deleteFolder(folderUid: string): Promise<ApiResponse<any>> {
    return this.request(`/folders/${folderUid}`, {
      method: 'DELETE',
    });
  }

  async moveNoteToFolder(noteUid: string, folderId?: number): Promise<ApiResponse<Note>> {
    return this.request(`/notes/${noteUid}/move-to-folder`, {
      method: 'POST',
      body: JSON.stringify({ folder_id: folderId }),
    });
  }

  // Chat Conversations
  async getChatConversations(projectUid: string): Promise<ApiResponse<ChatConversation[]>> {
    return this.request(`/projects/${projectUid}/chat/conversations`);
  }

  async createChatConversation(projectUid: string, req: ChatConversationRequest): Promise<ApiResponse<ChatConversation>> {
    return this.request(`/projects/${projectUid}/chat/conversations`, {
      method: 'POST',
      body: JSON.stringify(req),
    });
  }

  async getChatConversationWithMessages(conversationUid: string): Promise<ApiResponse<ChatConversationWithMessages>> {
    return this.request(`/chat/conversations/${conversationUid}`);
  }

  async deleteChatConversation(conversationUid: string): Promise<ApiResponse<any>> {
    return this.request(`/chat/conversations/${conversationUid}`, {
      method: 'DELETE',
    });
  }

  async createChatMessage(req: ChatMessageRequest): Promise<ApiResponse<ChatMessage>> {
    return this.request('/chat/messages', {
      method: 'POST',
      body: JSON.stringify(req),
    });
  }

  // AI Project Creation
  async createProjectFromAI(projectContent: string): Promise<ApiResponse<Project>> {
    return this.request<ApiResponse<Project>>('/ai/create-project', {
      method: 'POST',
      body: JSON.stringify({ project_content: projectContent }),
    });
  }

  // Task Comments
  async getTaskComments(taskUid: string): Promise<ApiResponse<TaskComment[]>> {
    return this.request(`/tasks/${taskUid}/comments`);
  }

  async createTaskComment(req: TaskCommentRequest): Promise<ApiResponse<TaskComment>> {
    return this.request('/tasks/comments', {
      method: 'POST',
      body: JSON.stringify(req),
    });
  }

  async updateTaskComment(commentUid: string, req: TaskCommentUpdateRequest): Promise<ApiResponse<TaskComment>> {
    return this.request(`/tasks/comments/${commentUid}`, {
      method: 'PUT',
      body: JSON.stringify(req),
    });
  }

  async deleteTaskComment(commentUid: string): Promise<ApiResponse<any>> {
    return this.request(`/tasks/comments/${commentUid}`, {
      method: 'DELETE',
    });
  }

  // Task Categories
  async getProjectCategories(projectUid: string): Promise<ApiResponse<TaskCategory[]>> {
    return this.request(`/projects/${projectUid}/categories`);
  }

  async createTaskCategory(req: TaskCategoryRequest): Promise<ApiResponse<TaskCategory>> {
    return this.request('/categories', {
      method: 'POST',
      body: JSON.stringify(req),
    });
  }

  async updateTaskCategory(categoryUid: string, req: TaskCategoryUpdateRequest): Promise<ApiResponse<TaskCategory>> {
    return this.request(`/categories/${categoryUid}`, {
      method: 'PUT',
      body: JSON.stringify(req),
    });
  }

  async deleteTaskCategory(categoryUid: string): Promise<ApiResponse<any>> {
    return this.request(`/categories/${categoryUid}`, {
      method: 'DELETE',
    });
  }

  async getTaskCategories(taskUid: string): Promise<ApiResponse<TaskCategory[]>> {
    return this.request(`/tasks/${taskUid}/categories`);
  }

  async assignCategoriesToTask(taskUid: string, categoryUids: string[]): Promise<ApiResponse<any>> {
    return this.request(`/tasks/${taskUid}/categories`, {
      method: 'POST',
      body: JSON.stringify(categoryUids),
    });
  }

  async removeCategoryFromTask(taskUid: string, categoryUid: string): Promise<ApiResponse<any>> {
    return this.request(`/tasks/${taskUid}/categories/${categoryUid}`, {
      method: 'DELETE',
    });
  }

  // Authentication methods
  async post<T>(endpoint: string, data: any): Promise<{ data: T }> {
    const response = await this.request<any>(endpoint, {
      method: 'POST',
      body: JSON.stringify(data),
    });
    // Backend wraps responses in { data: T, success: boolean, message?: string }
    // So we need to return the wrapped format for consistency
    return { data: response.data };
  }
}

export const apiService = new ApiService();