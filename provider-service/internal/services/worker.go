package services

import (
    "context"
    "encoding/json"
    "provider-service/internal/models"

    "github.com/redis/go-redis/v9"
    "github.com/rs/zerolog"
)

type Worker struct {
    redis    *redis.Client
    mockAPI  *MockAPI
    logger   *zerolog.Logger
    group    string
    consumer string
}

func NewWorker(redis *redis.Client, logger *zerolog.Logger) *Worker {
    return &Worker{
        redis:    redis,
        mockAPI:  NewMockAPI(),
        logger:   logger,
        group:    "flight-search-group",
        consumer: "provider-service",
    }
}

func (w *Worker) Start(ctx context.Context) error {
    // Create consumer group if not exists
    if err := w.createConsumerGroup(ctx); err != nil {
        return err
    }

    return w.processMessages(ctx)
}

func (w *Worker) createConsumerGroup(ctx context.Context) error {
    // Try to create the group
    err := w.redis.XGroupCreateMkStream(ctx, "flight.search.requested", w.group, "$").Err()
    if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
        return err
    }
    return nil
}

func (w *Worker) processMessages(ctx context.Context) error {
    w.logger.Info().Msg("Worker started processing messages")

    for {
        select {
        case <-ctx.Done():
            w.logger.Info().Msg("Worker shutting down")
            return nil
        default:
            entries, err := w.redis.XReadGroup(ctx, &redis.XReadGroupArgs{
                Group:    w.group,
                Consumer: w.consumer,
                Streams:  []string{"flight.search.requested", ">"},
                Block:    0,
            }).Result()

            if err != nil {
                w.logger.Error().Err(err).Msg("Error reading from Redis stream")
                continue
            }

            for _, stream := range entries {
                for _, message := range stream.Messages {
                    w.processMessage(ctx, &message)
                }
            }
        }
    }
}

func (w *Worker) processMessage(ctx context.Context, msg *redis.XMessage) {
    var request models.SearchRequest
    if err := json.Unmarshal([]byte(msg.Values["data"].(string)), &request); err != nil {
        w.logger.Error().Err(err).Msg("Failed to unmarshal message")
        return
    }

    w.logger.Info().
        Str("search_id", request.SearchID).
        Str("from", request.From).
        Str("to", request.To).
        Msg("Processing search request")

    // Simulate API call
    flights := w.mockAPI.SearchFlights(request)

    // Prepare result
    result := models.SearchResult{
        SearchID: request.SearchID,
        Status:   "completed",
        Results:  flights,
    }

    // Publish result
    data, err := json.Marshal(result)
    if err != nil {
        w.logger.Error().Err(err).Msg("Failed to marshal result")
        return
    }

    if err := w.redis.XAdd(ctx, &redis.XAddArgs{
        Stream: "flight.search.results",
        Values: map[string]interface{}{"data": string(data)},
    }).Err(); err != nil {
        w.logger.Error().Err(err).Msg("Failed to publish result")
        return
    }

    // Acknowledge message
    if err := w.redis.XAck(ctx, "flight.search.requested", w.group, msg.ID).Err(); err != nil {
        w.logger.Error().Err(err).Msg("Failed to acknowledge message")
    }

    w.logger.Info().
        Str("search_id", request.SearchID).
        Int("total_flights", len(flights)).
        Msg("Search completed and results published")
}