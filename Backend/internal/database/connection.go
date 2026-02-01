package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"lucid-lists-backend/internal/config"
	"lucid-lists-backend/pkg/logger"
)

func Connect(cfg *config.Config) (*pgxpool.Pool, error) {
	var dsn string

	// Use DATABASE_URL if available (for production/Render), otherwise build from components
	log := logger.WithComponent("database")
	if cfg.DatabaseURL != "" {
		dsn = cfg.DatabaseURL
		log.Info("Using DATABASE_URL for database connection")
	} else {
		// Build connection string from individual components (for local development)
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)
		log.Info("Using individual DB components for database connection")
	}

	// Configure connection pool
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Set pool settings
	config.MaxConns = 10
	config.MinConns = 2
	// Disable prepared statement cache to avoid "prepared statement does not exist" errors
	// This can happen when connections are recycled in the pool
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}
