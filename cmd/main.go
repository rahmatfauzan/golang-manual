package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/rahmatfauzan/golang-manual/internal/config"
	"github.com/rahmatfauzan/golang-manual/internal/database"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db := database.ConnectDB(cfg)
	defer db.Close()

	app := fiber.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Tahqiq API is running",
			"env":     cfg.APP_ENV,
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		if err := db.Ping(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unhealthy",
				"error":  err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	go func() {
		log.Printf("Server Running on http://localhost:%s", cfg.APP_PORT)
		if err := app.Listen(":" + cfg.APP_PORT); err != nil {
			log.Printf("failed to start server:%v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server exited")

}
