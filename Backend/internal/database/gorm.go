package database

import (
	"fmt"
	"lucid-lists-backend/internal/config"
	"lucid-lists-backend/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// ConnectGORM establishes a GORM database connection and runs automigrations
func ConnectGORM(cfg *config.Config) (*gorm.DB, error) {
	log := logger.WithComponent("database-gorm")

	var dsn string
	if cfg.DatabaseURL != "" {
		dsn = cfg.DatabaseURL
		log.Info("Using DATABASE_URL for GORM database connection")
	} else {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)
		log.Info("Using individual DB components for GORM database connection")
	}

	// Configure GORM with custom logger and disable prepared statements
	gormConfig := &gorm.Config{
		Logger:                                   gormlogger.Default.LogMode(gormlogger.Info),
		PrepareStmt:                              false, // Disable prepared statements to avoid conflicts
		DisableForeignKeyConstraintWhenMigrating: true,  // Let existing foreign keys stay
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database with GORM: %w", err)
	}

	log.Info("GORM database connected successfully")

	// Run automigrations
	if err := AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to run automigrations: %w", err)
	}

	log.Info("Database migrations completed successfully")
	return db, nil
}

// AutoMigrate runs all database migrations
func AutoMigrate(db *gorm.DB) error {
	log := logger.WithComponent("automigrate")
	log.Info("Running database migrations...")

	// Drop existing tables to recreate with correct foreign keys
	log.Info("Dropping existing tables if any...")
	db.Exec("DROP TABLE IF EXISTS task_category_map CASCADE")
	db.Exec("DROP TABLE IF EXISTS task_category CASCADE")
	db.Exec("DROP TABLE IF EXISTS task_comment CASCADE")

	// Create new tables manually using raw SQL to avoid GORM trying to recreate existing tables
	// This approach prevents GORM from following foreign key references and trying to create Task, User, etc.

	// 1. Create task_comment table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS task_comment (
			id SERIAL PRIMARY KEY,
			comment_uid UUID DEFAULT gen_random_uuid() UNIQUE NOT NULL,
			task_id INTEGER NOT NULL REFERENCES task(id) ON DELETE CASCADE,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP,
			is_active BOOLEAN DEFAULT true
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create task_comment table: %w", err)
	}
	log.Info("task_comment table ready")

	// 2. Create task_category table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS task_category (
			id SERIAL PRIMARY KEY,
			category_uid UUID DEFAULT gen_random_uuid() UNIQUE NOT NULL,
			project_id INTEGER NOT NULL REFERENCES project(id) ON DELETE CASCADE,
			name VARCHAR(100) NOT NULL,
			color VARCHAR(7) DEFAULT '#808080',
			description TEXT,
			created_at TIMESTAMP DEFAULT NOW(),
			created_by INTEGER REFERENCES users(id),
			updated_at TIMESTAMP,
			updated_by INTEGER REFERENCES users(id),
			is_active BOOLEAN DEFAULT true
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create task_category table: %w", err)
	}
	log.Info("task_category table ready")

	// 3. Create task_category_map table (many-to-many join)
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS task_category_map (
			id SERIAL PRIMARY KEY,
			task_id INTEGER NOT NULL REFERENCES task(id) ON DELETE CASCADE,
			category_id INTEGER NOT NULL REFERENCES task_category(id) ON DELETE CASCADE,
			assigned_at TIMESTAMP DEFAULT NOW(),
			assigned_by INTEGER REFERENCES users(id),
			UNIQUE(task_id, category_id)
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create task_category_map table: %w", err)
	}
	log.Info("task_category_map table ready")

	// 4. Create indexes for better performance
	db.Exec("CREATE INDEX IF NOT EXISTS idx_task_comment_task_id ON task_comment(task_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_task_comment_created_at ON task_comment(created_at)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_task_category_project_id ON task_category(project_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_task_category_map_task_id ON task_category_map(task_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_task_category_map_category_id ON task_category_map(category_id)")

	log.Info("All indexes created successfully")
	log.Info("New table migrations completed successfully")

	return nil
}
