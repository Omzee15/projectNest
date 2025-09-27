package routes

import (
	"net/http"
	"time"

	"lucid-lists-backend/internal/handlers"
	"lucid-lists-backend/internal/middleware"
	"lucid-lists-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(r *gin.Engine, projectHandler *handlers.ProjectHandler, listHandler *handlers.ListHandler, taskHandler *handlers.TaskHandler) {
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

		// Project routes
		projects := api.Group("/projects")
		{
			projects.GET("", projectHandler.GetProjects)
			projects.GET("/:uid", projectHandler.GetProject)
			projects.POST("", projectHandler.CreateProject)
			projects.PUT("/:uid", projectHandler.UpdateProject)
			projects.PATCH("/:uid", projectHandler.PartialUpdateProject)
			projects.DELETE("/:uid", projectHandler.DeleteProject)
		}

		// List routes
		lists := api.Group("/lists")
		{
			lists.POST("", listHandler.CreateList)
			lists.PUT("/:uid", listHandler.UpdateList)
			lists.PATCH("/:uid", listHandler.PartialUpdateList)
			lists.DELETE("/:uid", listHandler.DeleteList)
			lists.PUT("/:uid/position", listHandler.UpdatePosition)
		}

		// Task routes
		tasks := api.Group("/tasks")
		{
			tasks.POST("", taskHandler.CreateTask)
			tasks.PUT("/:uid", taskHandler.UpdateTask)
			tasks.PATCH("/:uid", taskHandler.PartialUpdateTask)
			tasks.DELETE("/:uid", taskHandler.DeleteTask)
			tasks.POST("/:uid/move", taskHandler.MoveTask)
		}
	}

	logger.WithComponent("router").Info("Routes setup completed successfully")
}
