package main

import (
	"fmt"
	"log"

	"github.com/Ashishkapoor1469/Nestgo/core"
	"github.com/Ashishkapoor1469/Nestgo/examples/basic-api/internal/modules"
)

func main() {
	app := core.New(
		core.WithAddress(":3000"),
		core.WithGlobalPrefix("/api"),
	)

	app.RegisterModule(&modules.AppModule{})

	fmt.Println("\n  🚀 Basic API running at http://localhost:3000")
	fmt.Println("  📋 Endpoints:")
	fmt.Println("       GET  /api/health")
	fmt.Println("       GET  /api/users")
	fmt.Println("       GET  /api/users/:id")
	fmt.Println("       POST /api/users")
	fmt.Println()

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
