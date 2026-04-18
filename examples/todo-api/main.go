package main

import (
	"log"

	"github.com/nestgo/nestgo/core"
	nesthttp "github.com/nestgo/nestgo/http"
	"github.com/nestgo/nestgo/middleware"

	"github.com/nestgo/nestgo/examples/todo-api/todos"

	"time"
)

func main() {
	// Create the NestGo application.
	app := core.New(
		core.WithAddress(":3000"),
		core.WithGlobalPrefix("/api"),
		core.WithCors("*"),
	)

	// Add rate limiting.
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)
	app.UseMiddleware(rateLimiter.Middleware())

	// Register modules.
	app.RegisterModule(&todos.TodosModule{})

	// Add OpenAPI documentation.
	openapi := nesthttp.NewOpenAPIGenerator("Todo API", "1.0.0")
	openapi.AddServer("http://localhost:3000", "Local Development")

	// Start the server.
	log.Println("🚀 Todo API starting on http://localhost:3000")
	log.Println("📚 API Base: http://localhost:3000/api")
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
