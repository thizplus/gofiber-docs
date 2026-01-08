package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
	"gofiber-template/interfaces/api/routes"
	"gofiber-template/pkg/di"
	"gofiber-template/pkg/logger"
)

func main() {
	// Initialize logger
	if err := logger.Init(logger.DefaultConfig()); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	logger.Info("Logger initialized", "log_file", "logs/app.log")

	// Initialize DI container
	container := di.NewContainer()

	// Initialize all dependencies
	if err := container.Initialize(); err != nil {
		log.Fatal("Failed to initialize container:", err)
	}

	// Setup graceful shutdown
	setupGracefulShutdown(container)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(),
		AppName:      container.GetConfig().App.Name,
	})

	// Setup middleware
	app.Use(middleware.LoggerMiddleware())
	app.Use(middleware.CorsMiddleware())

	// Create handlers from services
	services := container.GetHandlerServices()
	h := handlers.NewHandlers(services, container.GetConfig())

	// Setup routes
	routes.SetupRoutes(app, h)

	// Setup admin routes for API statistics
	api := app.Group("/api/v1")
	routes.SetupAdminRoutes(api, container.GetAPILoggerService())

	// Start server
	port := container.GetConfig().App.Port
	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("🌍 Environment: %s", container.GetConfig().App.Env)
	log.Printf("📚 Health check: http://localhost:%s/health", port)
	log.Printf("📖 API docs: http://localhost:%s/api/v1", port)
	log.Printf("🔌 WebSocket: ws://localhost:%s/ws", port)

	log.Fatal(app.Listen(":" + port))
}

func setupGracefulShutdown(container *di.Container) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("\n🛑 Gracefully shutting down...")

		if err := container.Cleanup(); err != nil {
			log.Printf("❌ Error during cleanup: %v", err)
		}

		log.Println("👋 Shutdown complete")
		os.Exit(0)
	}()
}