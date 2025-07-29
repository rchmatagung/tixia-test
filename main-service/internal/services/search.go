package services

import (
    "context"
    "main-service/internal/models"

    "github.com/google/uuid"
    "github.com/rs/zerolog"
)

type SearchService struct {
    redis  *RedisClient
    logger *zerolog.Logger
}

func NewSearchService(redis *RedisClient, logger *zerolog.Logger) *SearchService {
    return &SearchService{
        redis:  redis,
        logger: logger,
    }
}

func (s *SearchService) InitiateSearch(ctx context.Context, request models.SearchRequest) (string, error) {
    searchID := uuid.New().String()
    
    s.logger.Info().
        Str("search_id", searchID).
        Str("from", request.From).
        Str("to", request.To).
        Msg("Initiating flight search")

    if err := s.redis.PublishSearchRequest(ctx, searchID, request); err != nil {
        s.logger.Error().Err(err).Str("search_id", searchID).Msg("Failed to publish search request")
        return "", err
    }

    return searchID, nil
}

func (s *SearchService) SubscribeToResults(ctx context.Context, sseService *SSEService) error {
    results, err := s.redis.SubscribeToResults(ctx)
    if err != nil {
        return err
    }

    for result := range results {
        s.logger.Info().
            Str("search_id", result.SearchID).
            Str("status", result.Status).
            Int("total_results", len(result.Results)).
            Msg("Received search results")
        
        sseService.Broadcast(result.SearchID, result)
    }

    return nil
}