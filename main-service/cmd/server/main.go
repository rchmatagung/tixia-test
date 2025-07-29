package main

import (
    "context"
    "log"
    "main-service/internal/config"
    "main-service/internal/handlers"
    "main-service/internal/services"
    "os"
    "os/signal"
    "syscall"

    "github.com/gofiber/fiber/v2"
    fiberlog "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/recover"
    "github.com/rs/zerolog"
)

func main() {
    // Initialize logger
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
    logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

    // Load config
    cfg := config.Load()

    // Initialize Redis
    redisClient := services.NewRedisClient(cfg.RedisURL)
    defer redisClient.Close()

    // Initialize services
    searchService := services.NewSearchService(redisClient, &logger)
    sseService := services.NewSSEService(&logger)

    // Initialize handlers
    searchHandler := handlers.NewSearchHandler(searchService, sseService, &logger)

    // Setup Fiber app
    app := fiber.New(fiber.Config{
        AppName: "Flight Search Main Service",
    })

    app.Use(recover.New())
    app.Use(fiberlog.New())

    // Routes
    api := app.Group("/api")
    flights := api.Group("/flights")
    flights.Post("/search", searchHandler.SearchFlights)
    flights.Get("/search/:search_id/stream", searchHandler.StreamResults)

    // Start SSE listener
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go func() {
        if err := searchService.SubscribeToResults(ctx, sseService); err != nil {
            logger.Error().Err(err).Msg("Failed to subscribe to results")
        }
    }()

    // Start server
    go func() {
        if err := app.Listen(":" + cfg.Port); err != nil {
            log.Fatal("Failed to start server:", err)
        }
    }()

    // Wait for interrupt signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    logger.Info().Msg("Shutting down server...")
    cancel()
    _ = app.Shutdown()
}