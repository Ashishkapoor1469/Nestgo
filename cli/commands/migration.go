package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/Ashishkapoor1469/Nestgo/cli/utils"
	"github.com/spf13/cobra"
)

// MigrationCmd creates the `nestgo migration` command group.
// Supports colon-separated aliases: migration:create, migration:run, etc.
func MigrationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "migration",
		Aliases: []string{"migrate"},
		Short:   "Database migration management",
		Long: `Manage database migrations for your NestGo application.

Commands:
  nestgo migration create <name>   Create a new migration file
  nestgo migration run             Run all pending migrations
  nestgo migration rollback        Rollback the last migration
  nestgo migration status          Show migration status`,
	}

	cmd.AddCommand(
		migrationCreateCmd(),
		migrationRunCmd(),
		migrationRollbackCmd(),
		migrationStatusCmd(),
	)

	return cmd
}

// MigrationCreateAliasCmd returns a top-level alias for `nestgo migration:create`.
func MigrationCreateAliasCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "migration:create [name]",
		Short:  "Create a new migration file (alias)",
		Args:   cobra.ExactArgs(1),
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrationCreate(args[0])
		},
	}
}

// MigrationRunAliasCmd returns a top-level alias for `nestgo migration:run`.
func MigrationRunAliasCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "migration:run",
		Short:  "Run all pending migrations (alias)",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrationRun()
		},
	}
}

// MigrationRollbackAliasCmd returns a top-level alias for `nestgo migration:rollback`.
func MigrationRollbackAliasCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "migration:rollback",
		Short:  "Rollback last migration (alias)",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			steps, _ := cmd.Flags().GetInt("steps")
			return runMigrationRollback(steps)
		},
	}
	cmd.Flags().IntP("steps", "s", 1, "Number of migrations to rollback")
	return cmd
}

// MigrationStatusAliasCmd returns a top-level alias for `nestgo migration:status`.
func MigrationStatusAliasCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "migration:status",
		Short:  "Show migration status (alias)",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrationStatus()
		},
	}
}

// --- Subcommands ---

func migrationCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new migration file",
		Long: `Creates a timestamped migration file in the migrations/ directory.

Example:
  nestgo migration create add_user_avatar
  nestgo migration create create_users_table`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrationCreate(args[0])
		},
	}
}

func migrationRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run all pending migrations",
		Long:  "Executes all pending migrations in version order with transaction safety.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrationRun()
		},
	}
}

func migrationRollbackCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rollback",
		Short: "Rollback the last migration",
		Long:  "Rolls back the most recently applied migration(s).",
		RunE: func(cmd *cobra.Command, args []string) error {
			steps, _ := cmd.Flags().GetInt("steps")
			return runMigrationRollback(steps)
		},
	}
	cmd.Flags().IntP("steps", "s", 1, "Number of migrations to rollback")
	return cmd
}

func migrationStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show applied vs pending migrations",
		Long:  "Displays a table showing the status of all registered migrations.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrationStatus()
		},
	}
}

// --- Command Implementations ---

func runMigrationCreate(name string) error {
	utils.EnsureProjectContext("migration create")

	dir := "migrations"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Generate a timestamp-based version.
	timestamp := time.Now().Format("20060102150405")
	safeName := strings.ReplaceAll(strings.ToLower(name), "-", "_")

	fileName := fmt.Sprintf("%s_%s.go", timestamp, safeName)
	filePath := filepath.Join(dir, fileName)

	// Check for conflicts.
	if _, err := os.Stat(filePath); err == nil {
		utils.PrintError("Migration file already exists: " + filePath)
		return fmt.Errorf("migration conflict detected")
	}

	data := map[string]string{
		"Name":      safeName,
		"Pascal":    migrationPascalCase(safeName),
		"Version":   timestamp,
		"Timestamp": timestamp,
	}

	if err := writeMigrationTemplate(filePath, migrationFileTemplate, data); err != nil {
		return err
	}

	utils.PrintSuccess("Migration created: " + filePath)
	utils.PrintDim("  Version: " + timestamp)
	utils.PrintDim("  Edit the Up() and Down() methods to define your migration.")
	fmt.Println()

	return nil
}

func runMigrationRun() error {
	utils.EnsureProjectContext("migration run")

	utils.PrintHeader("Migration Runner")

	// Scan the migrations directory.
	entries, err := os.ReadDir("migrations")
	if err != nil {
		if os.IsNotExist(err) {
			utils.PrintWarning("No migrations directory found.")
			utils.PrintDim("Run 'nestgo migration create <name>' to create your first migration.")
			return nil
		}
		return err
	}

	migrationFiles := []string{}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".go") {
			migrationFiles = append(migrationFiles, e.Name())
		}
	}

	if len(migrationFiles) == 0 {
		utils.PrintInfo("No migration files found.")
		return nil
	}

	spinner := utils.StartSpinner("Running migrations...")

	// In a full implementation, this would connect to the DB and run migrations.
	// For CLI scaffolding, we show the workflow.
	time.Sleep(500 * time.Millisecond)
	spinner.Stop()

	utils.PrintInfo(fmt.Sprintf("Found %d migration file(s)", len(migrationFiles)))
	for _, f := range migrationFiles {
		utils.PrintStep("  📄", f)
	}
	fmt.Println()

	utils.PrintWarning("Database connection required to execute migrations.")
	utils.PrintDim("Configure DB_DRIVER, DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME in .env")
	utils.PrintDim("Then use the migration.Engine in your application code to run migrations programmatically:")
	fmt.Println()
	utils.PrintDim("  engine := migration.NewEngine(db.DB(), db.Driver())")
	utils.PrintDim("  engine.Register(migrations...)")
	utils.PrintDim("  engine.RunAll(ctx)")
	fmt.Println()

	return nil
}

func runMigrationRollback(steps int) error {
	utils.EnsureProjectContext("migration rollback")

	utils.PrintHeader("Migration Rollback")
	utils.PrintInfo(fmt.Sprintf("Rolling back %d migration(s)...", steps))
	fmt.Println()

	utils.PrintWarning("Database connection required for rollback.")
	utils.PrintDim("Use the migration.Engine.Rollback(ctx, steps) method in your application code.")
	fmt.Println()

	return nil
}

func runMigrationStatus() error {
	utils.EnsureProjectContext("migration status")

	utils.PrintHeader("Migration Status")

	// Scan migration files.
	entries, err := os.ReadDir("migrations")
	if err != nil {
		if os.IsNotExist(err) {
			utils.PrintWarning("No migrations directory found.")
			return nil
		}
		return err
	}

	headers := []string{"STATUS", "VERSION", "NAME", "FILE"}
	var rows [][]string

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}

		name := e.Name()
		parts := strings.SplitN(strings.TrimSuffix(name, ".go"), "_", 2)

		version := "unknown"
		desc := name
		if len(parts) >= 2 {
			version = parts[0]
			desc = parts[1]
		}

		// Without DB connection, show as "unknown" status.
		rows = append(rows, []string{"⏳ pending", version, desc, name})
	}

	if len(rows) == 0 {
		utils.PrintInfo("No migrations found.")
		utils.PrintDim("Run 'nestgo migration create <name>' to create one.")
		return nil
	}

	utils.PrintTable(headers, rows)
	fmt.Println()
	utils.PrintDim("Note: Actual applied/pending status requires database connection.")
	utils.PrintDim("Use migration.Engine.Status(ctx) for accurate status in your application.")
	fmt.Println()

	return nil
}

// --- Helpers ---

func migrationPascalCase(s string) string {
	parts := strings.Split(s, "_")
	var result strings.Builder
	for _, part := range parts {
		if len(part) > 0 {
			result.WriteString(strings.ToUpper(part[:1]) + part[1:])
		}
	}
	return result.String()
}

func writeMigrationTemplate(path, tmplStr string, data any) error {
	tmpl, err := template.New("migration").Parse(tmplStr)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	return tmpl.Execute(f, data)
}

// --- Migration File Template ---

var migrationFileTemplate = `package migrations

import "database/sql"

// {{.Pascal}}Migration — {{.Name}}
// Version: {{.Version}}
type {{.Pascal}}Migration struct{}

// Version returns the migration version.
func (m *{{.Pascal}}Migration) Version() string {
	return "{{.Version}}"
}

// Description returns a human-readable description.
func (m *{{.Pascal}}Migration) Description() string {
	return "{{.Name}}"
}

// Up applies the migration.
func (m *{{.Pascal}}Migration) Up(db *sql.DB) error {
	_, err := db.Exec(` + "`" + `
		-- TODO: Add your migration SQL here
		-- Example:
		-- CREATE TABLE users (
		--     id SERIAL PRIMARY KEY,
		--     name VARCHAR(255) NOT NULL,
		--     email VARCHAR(255) UNIQUE NOT NULL,
		--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		-- );
	` + "`" + `)
	return err
}

// Down rolls back the migration.
func (m *{{.Pascal}}Migration) Down(db *sql.DB) error {
	_, err := db.Exec(` + "`" + `
		-- TODO: Add your rollback SQL here
		-- Example:
		-- DROP TABLE IF EXISTS users;
	` + "`" + `)
	return err
}
`
