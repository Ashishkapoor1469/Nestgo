// Package migration provides a production-ready database migration engine
// with transaction-safe execution, ordered migration running, rollback support,
// and status tracking via a schema_migrations table.
package migration

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"
)

// Migration defines the interface that all migrations must implement.
type Migration interface {
	// Version returns the migration version string (e.g., "20260422120000").
	Version() string
	// Description returns a human-readable description of the migration.
	Description() string
	// Up applies the migration changes.
	Up(db *sql.DB) error
	// Down rolls back the migration changes.
	Down(db *sql.DB) error
}

// MigrationStatus represents the status of a single migration.
type MigrationStatus struct {
	Version     string
	Description string
	Applied     bool
	AppliedAt   *time.Time
}

// SimpleMigration is a convenience implementation of the Migration interface.
type SimpleMigration struct {
	Ver  string
	Desc string
	UpFn func(db *sql.DB) error
	DownFn func(db *sql.DB) error
}

func (m *SimpleMigration) Version() string     { return m.Ver }
func (m *SimpleMigration) Description() string { return m.Desc }
func (m *SimpleMigration) Up(db *sql.DB) error {
	if m.UpFn != nil {
		return m.UpFn(db)
	}
	return nil
}
func (m *SimpleMigration) Down(db *sql.DB) error {
	if m.DownFn != nil {
		return m.DownFn(db)
	}
	return nil
}

// Engine manages database migrations.
type Engine struct {
	db         *sql.DB
	driver     string
	migrations []Migration
	logger     *slog.Logger
	tableName  string
}

// EngineOption configures the migration engine.
type EngineOption func(*Engine)

// WithLogger sets a custom logger for the engine.
func WithLogger(logger *slog.Logger) EngineOption {
	return func(e *Engine) {
		e.logger = logger
	}
}

// WithTableName sets a custom migrations tracking table name.
func WithTableName(name string) EngineOption {
	return func(e *Engine) {
		e.tableName = name
	}
}

// NewEngine creates a new migration engine.
func NewEngine(db *sql.DB, driver string, opts ...EngineOption) *Engine {
	e := &Engine{
		db:        db,
		driver:    strings.ToLower(driver),
		tableName: "schema_migrations",
		logger:    slog.Default(),
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Register adds one or more migrations to the engine.
func (e *Engine) Register(migrations ...Migration) {
	e.migrations = append(e.migrations, migrations...)
	// Always keep migrations sorted by version.
	sort.Slice(e.migrations, func(i, j int) bool {
		return e.migrations[i].Version() < e.migrations[j].Version()
	})
}

// placeholder returns the correct SQL placeholder for the configured driver.
func (e *Engine) placeholder(n int) string {
	switch e.driver {
	case "mysql":
		return "?"
	default: // postgres, sqlite
		return fmt.Sprintf("$%d", n)
	}
}

// EnsureTable creates the schema_migrations table if it does not exist.
func (e *Engine) EnsureTable(ctx context.Context) error {
	var createSQL string
	switch e.driver {
	case "mysql":
		createSQL = fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				version VARCHAR(255) PRIMARY KEY,
				description TEXT,
				applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
		`, e.tableName)
	default: // postgres, sqlite
		createSQL = fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				version VARCHAR(255) PRIMARY KEY,
				description TEXT,
				applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)
		`, e.tableName)
	}

	_, err := e.db.ExecContext(ctx, createSQL)
	if err != nil {
		return fmt.Errorf("migration: failed to create %s table: %w", e.tableName, err)
	}
	return nil
}

// Applied returns the list of already-applied migration versions from the database.
func (e *Engine) Applied(ctx context.Context) (map[string]time.Time, error) {
	if err := e.EnsureTable(ctx); err != nil {
		return nil, err
	}

	query := fmt.Sprintf("SELECT version, applied_at FROM %s ORDER BY version", e.tableName)
	rows, err := e.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("migration: failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]time.Time)
	for rows.Next() {
		var version string
		var appliedAt time.Time
		if err := rows.Scan(&version, &appliedAt); err != nil {
			return nil, fmt.Errorf("migration: failed to scan applied migration: %w", err)
		}
		applied[version] = appliedAt
	}
	return applied, rows.Err()
}

// Pending returns migrations that have not yet been applied, in order.
func (e *Engine) Pending(ctx context.Context) ([]Migration, error) {
	applied, err := e.Applied(ctx)
	if err != nil {
		return nil, err
	}

	var pending []Migration
	for _, m := range e.migrations {
		if _, ok := applied[m.Version()]; !ok {
			pending = append(pending, m)
		}
	}
	return pending, nil
}

// RunAll runs all pending migrations in order.
// Each migration runs within its own transaction for safety.
// Returns the number of migrations applied and any error.
func (e *Engine) RunAll(ctx context.Context) (int, error) {
	pending, err := e.Pending(ctx)
	if err != nil {
		return 0, err
	}

	if len(pending) == 0 {
		e.logger.Info("migration: no pending migrations")
		return 0, nil
	}

	applied := 0
	for _, m := range pending {
		if err := e.runSingle(ctx, m); err != nil {
			return applied, fmt.Errorf("migration %s (%s) failed: %w", m.Version(), m.Description(), err)
		}
		applied++
		e.logger.Info("migration applied",
			"version", m.Version(),
			"description", m.Description(),
		)
	}

	return applied, nil
}

// runSingle runs a single migration within a transaction.
func (e *Engine) runSingle(ctx context.Context, m Migration) error {
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Run the migration's Up function.
	if err := m.Up(e.db); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("up failed: %w", err)
	}

	// Record the migration in the tracking table.
	insertSQL := fmt.Sprintf(
		"INSERT INTO %s (version, description) VALUES (%s, %s)",
		e.tableName,
		e.placeholder(1),
		e.placeholder(2),
	)
	if _, err := tx.ExecContext(ctx, insertSQL, m.Version(), m.Description()); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return tx.Commit()
}

// Rollback rolls back the last N applied migrations (default 1).
func (e *Engine) Rollback(ctx context.Context, steps int) (int, error) {
	if steps <= 0 {
		steps = 1
	}

	applied, err := e.Applied(ctx)
	if err != nil {
		return 0, err
	}

	if len(applied) == 0 {
		e.logger.Info("migration: nothing to rollback")
		return 0, nil
	}

	// Get applied migrations sorted by version descending (most recent first).
	type appliedEntry struct {
		version   string
		appliedAt time.Time
	}
	var entries []appliedEntry
	for v, at := range applied {
		entries = append(entries, appliedEntry{v, at})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].version > entries[j].version
	})

	// Limit to requested steps.
	if steps > len(entries) {
		steps = len(entries)
	}

	rolledBack := 0
	for i := 0; i < steps; i++ {
		version := entries[i].version

		// Find the migration.
		var mig Migration
		for _, m := range e.migrations {
			if m.Version() == version {
				mig = m
				break
			}
		}
		if mig == nil {
			return rolledBack, fmt.Errorf("migration %s not found in registry", version)
		}

		// Run rollback in a transaction.
		tx, err := e.db.BeginTx(ctx, nil)
		if err != nil {
			return rolledBack, fmt.Errorf("failed to begin transaction: %w", err)
		}

		if err := mig.Down(e.db); err != nil {
			_ = tx.Rollback()
			return rolledBack, fmt.Errorf("rollback of %s failed: %w", version, err)
		}

		deleteSQL := fmt.Sprintf(
			"DELETE FROM %s WHERE version = %s",
			e.tableName,
			e.placeholder(1),
		)
		if _, err := tx.ExecContext(ctx, deleteSQL, version); err != nil {
			_ = tx.Rollback()
			return rolledBack, fmt.Errorf("failed to remove migration record: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return rolledBack, fmt.Errorf("failed to commit rollback: %w", err)
		}

		rolledBack++
		e.logger.Info("migration rolled back",
			"version", version,
			"description", mig.Description(),
		)
	}

	return rolledBack, nil
}

// Status returns the status of all registered migrations.
func (e *Engine) Status(ctx context.Context) ([]MigrationStatus, error) {
	applied, err := e.Applied(ctx)
	if err != nil {
		return nil, err
	}

	var statuses []MigrationStatus
	for _, m := range e.migrations {
		s := MigrationStatus{
			Version:     m.Version(),
			Description: m.Description(),
		}
		if at, ok := applied[m.Version()]; ok {
			s.Applied = true
			s.AppliedAt = &at
		}
		statuses = append(statuses, s)
	}

	return statuses, nil
}
