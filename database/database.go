// Package database provides database abstraction with repository pattern,
// transaction support, and migration management.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

// Config holds database configuration.
type Config struct {
	Driver          string `env:"DB_DRIVER" default:"postgres"`
	Host            string `env:"DB_HOST" default:"localhost"`
	Port            int    `env:"DB_PORT" default:"5432"`
	User            string `env:"DB_USER" default:"postgres"`
	Password        string `env:"DB_PASSWORD"`
	Database        string `env:"DB_NAME" default:"nestgo"`
	SSLMode         string `env:"DB_SSL_MODE" default:"disable"`
	MaxOpenConns    int    `env:"DB_MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns    int    `env:"DB_MAX_IDLE_CONNS" default:"5"`
	ConnMaxLifetime int    `env:"DB_CONN_MAX_LIFETIME" default:"300"` // seconds
}

// DSN returns the data source name.
func (c *Config) DSN() string {
	switch c.Driver {
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			c.User, c.Password, c.Host, c.Port, c.Database)
	case "sqlite3", "sqlite":
		return c.Database
	default:
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
			c.Host, c.Port, c.User, c.Password, c.Database)
	}
}

// Database wraps a sql.DB connection.
type Database struct {
	db     *sql.DB
	config *Config
	logger *slog.Logger
}

// New creates a new database connection.
func New(cfg *Config, logger *slog.Logger) (*Database, error) {
	db, err := sql.Open(cfg.Driver, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("database: failed to open: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	// Test connection.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database: failed to ping: %w", err)
	}

	logger.Info("database connected",
		"driver", cfg.Driver,
		"host", cfg.Host,
		"database", cfg.Database,
	)

	return &Database{
		db:     db,
		config: cfg,
		logger: logger,
	}, nil
}

// DB returns the underlying sql.DB.
func (d *Database) DB() *sql.DB {
	return d.db
}

// Driver returns the configured database driver name.
func (d *Database) Driver() string {
	return d.config.Driver
}

// Logger returns the database logger.
func (d *Database) Logger() *slog.Logger {
	return d.logger
}

// Close closes the database connection.
func (d *Database) Close() error {
	d.logger.Info("closing database connection")
	return d.db.Close()
}

// Ping checks the database connection.
func (d *Database) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

// OnShutdown implements core.OnShutdown.
func (d *Database) OnShutdown() error {
	return d.Close()
}

// --- Transaction Support ---

// Tx wraps a database transaction.
type Tx struct {
	tx     *sql.Tx
	logger *slog.Logger
}

// BeginTx starts a new transaction.
func (d *Database) BeginTx(ctx context.Context) (*Tx, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("database: failed to begin transaction: %w", err)
	}
	return &Tx{tx: tx, logger: d.logger}, nil
}

// Commit commits the transaction.
func (t *Tx) Commit() error {
	return t.tx.Commit()
}

// Rollback rolls back the transaction.
func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}

// Tx returns the underlying sql.Tx.
func (t *Tx) Tx() *sql.Tx {
	return t.tx
}

// WithTransaction executes a function within a transaction.
func (d *Database) WithTransaction(ctx context.Context, fn func(tx *Tx) error) error {
	tx, err := d.BeginTx(ctx)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			d.logger.Error("transaction rollback failed", "error", rbErr)
		}
		return err
	}

	return tx.Commit()
}

// --- Repository Pattern ---

// Repository provides CRUD operations for an entity.
type Repository[T any] struct {
	db        *Database
	tableName string
	logger    *slog.Logger
}

// NewRepository creates a new repository for the given table.
func NewRepository[T any](db *Database, tableName string) *Repository[T] {
	return &Repository[T]{
		db:        db,
		tableName: tableName,
		logger:    db.logger.With("repository", tableName),
	}
}

// DB returns the underlying database.
func (r *Repository[T]) DB() *Database {
	return r.db
}

// Table returns the table name.
func (r *Repository[T]) Table() string {
	return r.tableName
}

// QueryRow executes a query returning a single row.
func (r *Repository[T]) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return r.db.db.QueryRowContext(ctx, query, args...)
}

// Query executes a query returning rows.
func (r *Repository[T]) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return r.db.db.QueryContext(ctx, query, args...)
}

// Exec executes a query without returning rows.
func (r *Repository[T]) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return r.db.db.ExecContext(ctx, query, args...)
}

// --- Migration (Deprecated) ---

// Deprecated: Migration is the legacy migration type.
// Use the database/migration package for production-grade migrations with
// transaction safety, rollback support, and CLI integration.
type Migration struct {
	Version     string
	Description string
	Up          func(db *sql.DB) error
	Down        func(db *sql.DB) error
}

// Deprecated: Migrator is the legacy migration manager.
// Use the database/migration.Engine for production-grade migrations.
type Migrator struct {
	db         *Database
	migrations []Migration
	logger     *slog.Logger
}

// Deprecated: NewMigrator creates a new legacy migrator.
// Use migration.NewEngine() from the database/migration package instead.
func NewMigrator(db *Database) *Migrator {
	return &Migrator{
		db:     db,
		logger: db.logger.With("component", "migrator"),
	}
}

// Add adds a migration.
func (m *Migrator) Add(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

// Up runs all pending migrations.
func (m *Migrator) Up() error {
	// Create migrations table if not exists.
	_, err := m.db.db.Exec(`
		CREATE TABLE IF NOT EXISTS _nestgo_migrations (
			version VARCHAR(255) PRIMARY KEY,
			description TEXT,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("migrator: failed to create migrations table: %w", err)
	}

	for _, migration := range m.migrations {
		// Check if already applied.
		var count int
		err := m.db.db.QueryRow("SELECT COUNT(*) FROM _nestgo_migrations WHERE version = $1", migration.Version).Scan(&count)
		if err != nil {
			return fmt.Errorf("migrator: failed to check migration status: %w", err)
		}
		if count > 0 {
			continue
		}

		m.logger.Info("applying migration", "version", migration.Version, "description", migration.Description)
		if err := migration.Up(m.db.db); err != nil {
			return fmt.Errorf("migrator: migration %s failed: %w", migration.Version, err)
		}

		_, err = m.db.db.Exec("INSERT INTO _nestgo_migrations (version, description) VALUES ($1, $2)",
			migration.Version, migration.Description)
		if err != nil {
			return fmt.Errorf("migrator: failed to record migration: %w", err)
		}
	}

	return nil
}

// Down rolls back the last migration.
func (m *Migrator) Down() error {
	if len(m.migrations) == 0 {
		return nil
	}

	last := m.migrations[len(m.migrations)-1]
	m.logger.Info("rolling back migration", "version", last.Version)

	if last.Down != nil {
		if err := last.Down(m.db.db); err != nil {
			return fmt.Errorf("migrator: rollback %s failed: %w", last.Version, err)
		}
	}

	_, err := m.db.db.Exec("DELETE FROM _nestgo_migrations WHERE version = $1", last.Version)
	return err
}
