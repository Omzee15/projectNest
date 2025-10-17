package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"lucid-lists-backend/internal/config"
	"lucid-lists-backend/internal/database"
	"lucid-lists-backend/internal/handlers"
	"lucid-lists-backend/internal/middleware"
	"lucid-lists-backend/internal/repositories"
	"lucid-lists-backend/internal/routes"
	"lucid-lists-backend/internal/services"
	"lucid-lists-backend/pkg/logger"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logger.WithComponent("main").Warn("No .env file found")
	}

	// Initialize logger
	logger.Init()
	log := logger.WithComponent("main")
	log.Info("Starting Lucid Lists Backend Server")

	// Load configuration
	cfg := config.Load()
	log.Info("Configuration loaded")

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Info("Database connected successfully")

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	projectRepo := repositories.NewProjectRepository(db, userRepo)
	listRepo := repositories.NewListRepository(db)
	taskRepo := repositories.NewTaskRepository(db)
	// Phase 3: Brainstorming & Planning Layer repositories
	canvasRepo := repositories.NewCanvasRepository(db)
	noteRepo := repositories.NewNoteRepository(db)
	folderRepo := repositories.NewNoteFolderRepository(db)
	chatRepo := repositories.NewChatRepository(db)
	settingsRepo := repositories.NewUserSettingsRepository(db)

	// Initialize services
	projectService := services.NewProjectService(projectRepo, listRepo, taskRepo, userRepo)
	listService := services.NewListService(listRepo, taskRepo, projectRepo)
	taskService := services.NewTaskService(taskRepo, listRepo)
	// Phase 3: Brainstorming & Planning Layer services
	canvasService := services.NewCanvasService(canvasRepo, projectRepo)
	noteService := services.NewNoteService(noteRepo, projectRepo)
	folderService := services.NewNoteFolderService(folderRepo, projectRepo)
	chatService := services.NewChatService(chatRepo, projectRepo, userRepo)
	settingsService := services.NewUserSettingsService(settingsRepo)

	// Initialize handlers
	projectHandler := handlers.NewProjectHandler(projectService)
	listHandler := handlers.NewListHandler(listService)
	taskHandler := handlers.NewTaskHandler(taskService)
	// Phase 3: Brainstorming & Planning Layer handlers
	canvasHandler := handlers.NewCanvasHandler(canvasService)
	noteHandler := handlers.NewNoteHandler(noteService)
	folderHandler := handlers.NewNoteFolderHandler(folderService)
	chatHandler := handlers.NewChatHandler(chatService)
	settingsHandler := handlers.NewUserSettingsHandler(settingsService)

	// AI project creation handler
	aiProjectHandler := handlers.NewAIProjectCreationHandler(projectService, listService, taskService, canvasService, chatService)

	// Authentication layer
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	// Setup router
	router := setupRouter(cfg)

	// Setup routes
	routes.SetupRoutes(router, projectHandler, listHandler, taskHandler, canvasHandler, noteHandler, folderHandler, chatHandler, authHandler, aiProjectHandler, settingsHandler, authService, projectRepo, db)

	// Create server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Infof("Server starting on port %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Info("Server exited")
}

func setupRouter(cfg *config.Config) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware - CORS must be first to handle preflight requests
	router.Use(gin.Recovery())
	router.Use(middleware.CORSWithLogging(cfg.CORSAllowedOrigins))

	logger.WithComponent("router").Info("Router setup completed successfully")
	return router
}
