package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/internal/services"
	"lucid-lists-backend/internal/utils"
	"google.golang.org/api/option"
)

type AIProjectCreationHandler struct {
	projectService *services.ProjectService
	listService    *services.ListService
	taskService    *services.TaskService
	canvasService  *services.CanvasService
	chatService    *services.ChatService
}

func NewAIProjectCreationHandler(
	projectService *services.ProjectService,
	listService *services.ListService,
	taskService *services.TaskService,
	canvasService *services.CanvasService,
	chatService *services.ChatService,
) *AIProjectCreationHandler {
	return &AIProjectCreationHandler{
		projectService: projectService,
		listService:    listService,
		taskService:    taskService,
		canvasService:  canvasService,
		chatService:    chatService,
	}
}

// AIProjectCreationRequest represents the request payload for AI project creation
type AIProjectCreationRequest struct {
	ProjectContent string `json:"project_content" binding:"required"`
}

// Internal structures for parsed project data
type ParsedProjectData struct {
	ProjectName    string       `json:"project_name"`
	Description    string       `json:"description"`
	Lists          []ParsedList `json:"lists"`
	Flowchart      string       `json:"flowchart"`
	DatabaseSchema string       `json:"database_schema"`
}

type ParsedList struct {
	Name  string       `json:"name"`
	Tasks []ParsedTask `json:"tasks"`
}

type ParsedTask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// CreateProjectFromAI handles AI-generated project creation
func (h *AIProjectCreationHandler) CreateProjectFromAI(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userUID, exists := c.Get("user_uid")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User UID not found")
		return
	}

	var req AIProjectCreationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Parse the project content using AI
	parsedData, err := h.parseProjectContentWithAI(req.ProjectContent)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to parse project content: "+err.Error())
		return
	}

	// Create project
	projectReq := &models.ProjectRequest{
		Name:        parsedData.ProjectName,
		Description: &parsedData.Description,
		Status:      "active",
		Color:       "#3B82F6", // Default blue color
	}

	ctx := context.Background()
	project, err := h.projectService.CreateProject(ctx, projectReq, userID.(int), userUID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create project: "+err.Error())
		return
	}

	// Create lists and tasks
	for _, listData := range parsedData.Lists {
		listReq := &models.ListRequest{
			ProjectUID: project.ProjectUID,
			Name:       listData.Name,
			Color:      "#6B7280", // Default gray color
			Position:   0,
		}

		list, err := h.listService.CreateList(ctx, listReq)
		if err != nil {
			continue // Skip failed lists but continue with others
		}

		// Create tasks for this list
		for _, taskData := range listData.Tasks {
			taskReq := &models.TaskRequest{
				ListUID:     list.ListUID,
				Title:       taskData.Title,
				Description: &taskData.Description,
				Priority:    stringPtr("medium"),
				Status:      "todo",
				Color:       "#F3F4F6", // Default light gray color
			}

			_, err := h.taskService.CreateTask(ctx, taskReq)
			if err != nil {
				continue // Skip failed tasks but continue with others
			}
		}
	}

	// Note: Canvas creation is optional and can be added later if needed
	// For now, we'll store flowchart and DB schema in the project description or as separate fields

	// Return success response
	utils.SuccessResponse(c, project, "Project created successfully with AI assistance")
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}

// parseProjectContentWithAI uses Google Gemini to parse the project content
func (h *AIProjectCreationHandler) parseProjectContentWithAI(content string) (*ParsedProjectData, error) {
	ctx := context.Background()

	// Get API key from environment
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}

	// Initialize Gemini client
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %v", err)
	}
	defer client.Close()

	// Get the model
	model := client.GenerativeModel("gemini-1.5-flash")

	// Create system prompt for parsing
	systemPrompt := `You are a project management assistant. Parse the following project description and extract structured information in JSON format.

Please analyze the content and extract:
1. project_name: A clear, concise name for the project
2. description: A brief description of the project
3. lists: An array of task lists/categories with their tasks
4. flowchart: If mentioned, create a mermaid flowchart script
5. database_schema: If mentioned, provide database schema information

Return ONLY valid JSON in this exact format:
{
  "project_name": "string",
  "description": "string", 
  "lists": [
    {
      "name": "string",
      "tasks": [
        {
          "title": "string",
          "description": "string"
        }
      ]
    }
  ],
  "flowchart": "string (mermaid script or empty)",
  "database_schema": "string (schema description or empty)"
}

Content to parse:`

	// Generate response
	resp, err := model.GenerateContent(ctx, genai.Text(systemPrompt+"\n\n"+content))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %v", err)
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response generated")
	}

	// Extract text from response
	var responseText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			responseText += string(txt)
		}
	}

	// Clean the response (remove markdown code blocks if present)
	responseText = strings.TrimSpace(responseText)
	responseText = strings.TrimPrefix(responseText, "```json")
	responseText = strings.TrimPrefix(responseText, "```")
	responseText = strings.TrimSuffix(responseText, "```")
	responseText = strings.TrimSpace(responseText)

	// Parse JSON response
	var parsedData ParsedProjectData
	if err := json.Unmarshal([]byte(responseText), &parsedData); err != nil {
		return nil, fmt.Errorf("failed to parse AI response as JSON: %v\nResponse: %s", err, responseText)
	}

	return &parsedData, nil
}