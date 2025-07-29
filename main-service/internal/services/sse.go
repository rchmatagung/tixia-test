package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"main-service/internal/models"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type SSEService struct {
	clients   map[string]map[string]chan models.SearchResult
	results   map[string]models.SearchResult
	resultsMu sync.RWMutex
	mutex     sync.RWMutex
	logger    *zerolog.Logger
}

func NewSSEService(logger *zerolog.Logger) *SSEService {
	return &SSEService{
		clients: make(map[string]map[string]chan models.SearchResult),
		results: make(map[string]models.SearchResult),
		logger:  logger,
	}
}

func (s *SSEService) AddClient(searchID string, clientID string) chan models.SearchResult {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.clients[searchID]; !exists {
		s.clients[searchID] = make(map[string]chan models.SearchResult)
	}

	ch := make(chan models.SearchResult, 10)
	s.clients[searchID][clientID] = ch

	s.logger.Info().
		Str("search_id", searchID).
		Str("client_id", clientID).
		Msg("Client connected to SSE")

	return ch
}

func (s *SSEService) RemoveClient(searchID string, clientID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if clients, exists := s.clients[searchID]; exists {
		close(clients[clientID])
		delete(clients, clientID)

		if len(clients) == 0 {
			delete(s.clients, searchID)
		}
	}

	s.logger.Info().
		Str("search_id", searchID).
		Str("client_id", clientID).
		Msg("Client disconnected from SSE")
}

func (s *SSEService) Broadcast(searchID string, result models.SearchResult) {

	s.resultsMu.Lock()
	s.results[searchID] = result
	s.resultsMu.Unlock()

	fmt.Printf(">>> [BROADCAST] saved result for %s\n", searchID)

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if clients, exists := s.clients[searchID]; exists {
		for clientID, ch := range clients {
			select {
			case ch <- result:
				s.logger.Debug().
					Str("search_id", searchID).
					Str("client_id", clientID).
					Msg("Sent SSE update to client")
			default:
				s.logger.Warn().
					Str("search_id", searchID).
					Str("client_id", clientID).
					Msg("Client channel full, dropping message")
			}
		}
	} else {
		s.logger.Debug().
			Str("search_id", searchID).
			Msg("No SSE clients connected for this search")
	}
}

func (s *SSEService) StreamHandler(c *fiber.Ctx, searchID string) error {
	clientID := c.Get("X-Client-ID")
	if clientID == "" {
		clientID = "anonymous"
	}

	s.logger.Info().
		Str("search_id", searchID).
		Str("client_id", clientID).
		Msg("New SSE client connected")

	ch := s.AddClient(searchID, clientID)
	defer s.RemoveClient(searchID, clientID)

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		s.resultsMu.RLock()
		if res, ok := s.results[searchID]; ok {
			s.resultsMu.RUnlock()
			fmt.Printf(">>> [CACHE HIT] sending cached result for %s\n", searchID)
			data, _ := json.Marshal(res)
			w.WriteString("data: " + string(data) + "\n\n")
			w.Flush()
			return
		}
		s.resultsMu.RUnlock()
		fmt.Printf(">>> [CACHE MISS] no result yet, waiting for channel %s\n", searchID)

		w.WriteString("data: {\"search_id\":\"" + searchID + "\",\"status\":\"processing\",\"results\":[]}\n\n")
		w.Flush()

		for result := range ch {
			data, _ := json.Marshal(result)
			w.WriteString("data: " + string(data) + "\n\n")
			w.Flush()

			if result.Status == "completed" {
				return
			}
		}
	})

	return nil
}
