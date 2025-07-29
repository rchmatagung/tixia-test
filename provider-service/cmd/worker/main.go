package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "provider-service/internal/config"
    "provider-service/internal/services"
    "syscall"

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
    worker := services.NewWorker(redisClient, &logger)

    // Setup graceful shutdown
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sigChan
        logger.Info().Msg("Shutting down worker...")
        cancel()
    }()

    // Start worker
    if err := worker.Start(ctx); err != nil {
        log.Fatal("Worker failed:", err)
    }
}