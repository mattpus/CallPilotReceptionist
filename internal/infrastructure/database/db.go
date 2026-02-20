package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/CallPilotReceptionist/pkg/config"
	"github.com/CallPilotReceptionist/pkg/logger"
)

type DB struct {
	*sql.DB
	logger *logger.Logger
}

func NewDB(cfg *config.DatabaseConfig, log *logger.Logger) (*DB, error) {
	db, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Database connection established", nil)

	return &DB{
		DB:     db,
		logger: log,
	}, nil
}

func (db *DB) Close() error {
	db.logger.Info("Closing database connection", nil)
	return db.DB.Close()
}

func (db *DB) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	return db.PingContext(ctx)
}
