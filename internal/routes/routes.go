package routes

import (
	"net/http"
	"time"

	"lucid-lists-backend/internal/handlers"
	"lucid-lists-backend/internal/middleware"
	"lucid-lists-backend/internal/repositories"
	"lucid-lists-backend/internal/services"
	"lucid-lists-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(r *gin.Engine, projectHandler *handlers.ProjectHandler, listHandler *handlers.ListHandler, taskHandler *handlers.TaskHandler, canvasHandler *handlers.CanvasHandler, noteHandler *handlers.NoteHandler, folderHandler *handlers.NoteFolderHandler, chatHandler *handlers.ChatHandler, authHandler *handlers.AuthHandler, aiProjectHandler *handlers.AIProjectCreationHandler, settingsHandler *handlers.UserSettingsHandler, authService *services.AuthService, projectRepo repositories.ProjectRepository, db *pgxpool.Pool) {
	// Add middleware
	r.Use(middleware.RequestLogging())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		logger.WithComponent("health").Info("Health check accessed")
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Log API group initialization
		logger.WithComponent("router").Info("Initializing API routes")

		// Authentication routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes - require authentication
		protected := api.Group("")
		protected.Use(middleware.Authentication(authService))
		{
			// Project routes
			projects := protected.Group("/projects")
			{
				// Public routes (no ownership validation needed)
				projects.GET("", projectHandler.GetProjectsWithProgress)
				projects.POST("", projectHandler.CreateProject)

				// Project-specific routes (require ownership validation)
				projectOwned := projects.Group("/:uid")
				projectOwned.Use(middleware.ProjectOwnership(projectRepo, db))
				{
					projectOwned.GET("", projectHandler.GetProject)
					projectOwned.GET("/progress", projectHandler.GetProjectProgress)
					projectOwned.PUT("", projectHandler.UpdateProject)
					projectOwned.PATCH("", projectHandler.PartialUpdateProject)
					projectOwned.DELETE("", projectHandler.DeleteProject)

					// Project members routes
					projectOwned.GET("/members", projectHandler.GetProjectMembers)
					projectOwned.POST("/members", projectHandler.AddProjectMember)

					// Phase 3: Canvas routes (require project ownership)
					projectOwned.GET("/canvas", canvasHandler.GetCanvas)
					projectOwned.POST("/canvas", canvasHandler.UpdateCanvas)
					projectOwned.PUT("/canvas", canvasHandler.UpdateCanvas)
					projectOwned.DELETE("/canvas", canvasHandler.DeleteCanvas)

					// Phase 3: Notes routes (require project ownership)
					projectOwned.GET("/notes", noteHandler.GetNotesByProject)
					projectOwned.POST("/notes", noteHandler.CreateNote)

					// Phase 3: Folder routes (require project ownership)
					projectOwned.GET("/folders", folderHandler.GetFolders)
					projectOwned.POST("/folders", folderHandler.CreateFolder)

					// Phase 3: Chat routes (require project ownership)
					projectOwned.GET("/chat/conversations", chatHandler.GetConversations)
					projectOwned.POST("/chat/conversations", chatHandler.CreateConversation)
				}
			}

			// List routes
			lists := protected.Group("/lists")
			{
				lists.POST("", listHandler.CreateList)
				lists.PUT("/:uid", listHandler.UpdateList)
				lists.PATCH("/:uid", listHandler.PartialUpdateList)
				lists.DELETE("/:uid", listHandler.DeleteList)
				lists.PUT("/:uid/position", listHandler.UpdatePosition)
			}

			// Task routes
			tasks := protected.Group("/tasks")
			{
				tasks.POST("", taskHandler.CreateTask)
				tasks.PUT("/:uid", taskHandler.UpdateTask)
				tasks.PATCH("/:uid", taskHandler.PartialUpdateTask)
				tasks.DELETE("/:uid", taskHandler.DeleteTask)
				tasks.POST("/:uid/move", taskHandler.MoveTask)
			}

			// Phase 3: Note routes (individual note operations)
			// Note: These are protected through authentication middleware
			notes := protected.Group("/notes")
			{
				notes.GET("/:uid", noteHandler.GetNote)
				notes.PUT("/:uid", noteHandler.UpdateNote)
				notes.PATCH("/:uid", noteHandler.PartialUpdateNote)
				notes.DELETE("/:uid", noteHandler.DeleteNote)
				notes.POST("/:uid/move-to-folder", noteHandler.MoveNoteToFolder)
			}

			// Phase 3: Folder routes (individual folder operations)
			// Note: These are protected through authentication middleware
			folders := protected.Group("/folders")
			{
				folders.PUT("/:uid", folderHandler.UpdateFolder)
				folders.DELETE("/:uid", folderHandler.DeleteFolder)
			}

			// Phase 3: Chat routes (individual conversation operations)
			// Note: These are protected through authentication middleware
			chat := protected.Group("/chat")
			{
				chat.GET("/conversations/:conversationUid", chatHandler.GetConversationWithMessages)
				chat.DELETE("/conversations/:conversationUid", chatHandler.DeleteConversation)
				chat.POST("/messages", chatHandler.CreateMessage)
			}

			// User settings routes
			settings := protected.Group("/settings")
			{
				settings.GET("", settingsHandler.GetUserSettings)
				settings.PUT("", settingsHandler.UpdateUserSettings)
				settings.PATCH("", settingsHandler.UpdateUserSettings) // Support both PUT and PATCH
				settings.POST("/reset", settingsHandler.ResetUserSettings)
			}

			// AI project creation routes
			ai := protected.Group("/ai")
			{
				ai.POST("/create-project", aiProjectHandler.CreateProjectFromAI)
			}
		}
	}

	logger.WithComponent("router").Info("Routes setup completed successfully")
}
