package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/cliffordotieno/ai-context-gap-tracker/internal/config"
	_ "github.com/lib/pq"
)

// DB wraps sql.DB
type DB struct {
	*sql.DB
}

// NewConnection creates a new database connection
func NewConnection(cfg config.DatabaseConfig) (*DB, error) {
	db, err := sql.Open("postgres", cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &DB{DB: db}, nil
}

// runMigrations runs database migrations
func runMigrations(db *sql.DB) error {
	migrations := []string{
		createContextTable,
		createRulesTable,
		createAuditTable,
		createSessionsTable,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Database schema definitions
const (
	createContextTable = `
		CREATE TABLE IF NOT EXISTS contexts (
			id SERIAL PRIMARY KEY,
			session_id VARCHAR(255) NOT NULL,
			turn_number INTEGER NOT NULL,
			user_input TEXT NOT NULL,
			entities JSONB,
			topics JSONB,
			timeline JSONB,
			assertions JSONB,
			ambiguities JSONB,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(session_id, turn_number)
		);
		CREATE INDEX IF NOT EXISTS idx_contexts_session_id ON contexts(session_id);
		CREATE INDEX IF NOT EXISTS idx_contexts_turn_number ON contexts(turn_number);
		CREATE INDEX IF NOT EXISTS idx_contexts_entities ON contexts USING GIN(entities);
		CREATE INDEX IF NOT EXISTS idx_contexts_topics ON contexts USING GIN(topics);
	`

	createRulesTable = `
		CREATE TABLE IF NOT EXISTS rules (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			description TEXT,
			rule_type VARCHAR(50) NOT NULL,
			conditions JSONB NOT NULL,
			actions JSONB NOT NULL,
			priority INTEGER DEFAULT 0,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_rules_type ON rules(rule_type);
		CREATE INDEX IF NOT EXISTS idx_rules_priority ON rules(priority);
		CREATE INDEX IF NOT EXISTS idx_rules_active ON rules(is_active);
	`

	createAuditTable = `
		CREATE TABLE IF NOT EXISTS audit_logs (
			id SERIAL PRIMARY KEY,
			session_id VARCHAR(255) NOT NULL,
			turn_number INTEGER NOT NULL,
			response_text TEXT NOT NULL,
			certainty_level VARCHAR(50) NOT NULL,
			flags JSONB,
			assumptions JSONB,
			contradictions JSONB,
			retry_count INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_audit_session_id ON audit_logs(session_id);
		CREATE INDEX IF NOT EXISTS idx_audit_certainty ON audit_logs(certainty_level);
		CREATE INDEX IF NOT EXISTS idx_audit_flags ON audit_logs USING GIN(flags);
	`

	createSessionsTable = `
		CREATE TABLE IF NOT EXISTS sessions (
			id VARCHAR(255) PRIMARY KEY,
			user_id VARCHAR(255),
			context_graph JSONB,
			memory_state JSONB,
			last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
		CREATE INDEX IF NOT EXISTS idx_sessions_last_activity ON sessions(last_activity);
	`
)