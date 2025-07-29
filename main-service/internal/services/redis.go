package services

import (
	"context"
	"encoding/json"
	"main-service/internal/models"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(url string) *RedisClient {
	opt, _ := redis.ParseURL(url)
	return &RedisClient{
		client: redis.NewClient(opt),
	}
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func (r *RedisClient) PublishSearchRequest(ctx context.Context, searchID string, request models.SearchRequest) error {
	message := map[string]interface{}{
		"search_id":  searchID,
		"from":       request.From,
		"to":         request.To,
		"date":       request.Date,
		"passengers": request.Passengers,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: "flight.search.requested",
		Values: map[string]interface{}{"data": string(data)},
	}).Err()
}

func (r *RedisClient) SubscribeToResults(ctx context.Context) (<-chan models.SearchResult, error) {
	results := make(chan models.SearchResult)

	go func() {
		defer close(results)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				entries, err := r.client.XRead(ctx, &redis.XReadArgs{
					Streams: []string{"flight.search.results", "$"},
					Block:   0,
				}).Result()

				if err != nil {
					continue
				}

				for _, stream := range entries {
					for _, message := range stream.Messages {
						var result models.SearchResult
						if err := json.Unmarshal([]byte(message.Values["data"].(string)), &result); err != nil {
							continue
						}
						results <- result
					}
				}
			}
		}
	}()

	return results, nil
}
