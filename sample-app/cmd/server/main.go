package main

import (
	"log"

	"github.com/Ashishkapoor1469/Nestgo/core"
	"github.com/Ashishkapoor1469/Nestgo/config"
	"github.com/sample-app/internal/modules"
	appconfig "github.com/sample-app/internal/config"
)

func main() {
	// Load configuration.
	cfg := config.MustLoad[appconfig.AppConfig](".")

	// Create the application.
	app := core.New(
		core.WithAddress(":"+cfg.Port),
		core.WithGlobalPrefix("/api"),
	)

	// Register root module.
	app.RegisterModule(&modules.AppModule{})

	// Start the server.
	log.Printf("🚀 Starting %s on :%s", "sample-app", cfg.Port)
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
