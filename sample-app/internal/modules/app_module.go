package modules

import (
	"log/slog"

	"github.com/Ashishkapoor1469/Nestgo/common"
	"github.com/Ashishkapoor1469/Nestgo/database"
	"github.com/Ashishkapoor1469/Nestgo/di"

	_ "github.com/lib/pq"
	
	"github.com/sample-app/internal/modules/auth"
	"github.com/sample-app/internal/modules/users"
)

// AppModule is the root module of the application.
type AppModule struct{}

func (m *AppModule) Module() common.ModuleConfig {
	return common.ModuleConfig{
		Name: "app",
		Imports: []common.Module{
			&auth.AuthModule{},
			&users.UsersModule{},
		},
		Controllers: []common.Controller{},
		Providers: []di.Provider{
			{
				// Database Factory using pure Go functional DI injection.
				// The DI container automatically provides *slog.Logger.
				Factory: func(logger *slog.Logger) (*database.Database, error) {
					cfg := &database.Config{
						Driver:   "postgres",
						Host:     "localhost",
						Port:     5432,
						User:     "postgres",
						Password: "password123", // Update with your local pg connection
						Database: "nestgo_demo",
						SSLMode:  "disable",
					}

					db, err := database.New(cfg, logger)
					if err != nil {
						logger.Error("Failed to connect to database. Make sure postgres is running!", "err", err)
						// We don't return error to let the framework boot in degraded mode for the sample if PG is off
						// return nil, err
					}

					// Let's create our users table automatically on boot for this sample!
					if db != nil {
						_, err := db.DB().Exec(`
							CREATE TABLE IF NOT EXISTS users (
								id SERIAL PRIMARY KEY,
								name VARCHAR(255) NOT NULL,
								email VARCHAR(255) UNIQUE NOT NULL,
								password_hash VARCHAR(255) NOT NULL,
								created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
								updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
							)
						`)
						if err != nil {
							logger.Error("Failed to auto-migrate users table", "err", err)
						} else {
							logger.Info("Users table migration verified ✅")
						}
					}

					return db, nil
				},
			},
		},
	}
}
