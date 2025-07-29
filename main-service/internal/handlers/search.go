package handlers

import (
    "main-service/internal/models"
    "main-service/internal/services"

    "github.com/gofiber/fiber/v2"
    "github.com/go-playground/validator/v10"
    "github.com/rs/zerolog"
)

type SearchHandler struct {
    searchService *services.SearchService
    sseService    *services.SSEService
    logger        *zerolog.Logger
    validator     *validator.Validate
}

func NewSearchHandler(searchService *services.SearchService, sseService *services.SSEService, logger *zerolog.Logger) *SearchHandler {
    return &SearchHandler{
        searchService: searchService,
        sseService:    sseService,
        logger:        logger,
        validator:     validator.New(),
    }
}

func (h *SearchHandler) SearchFlights(c *fiber.Ctx) error {
    var req models.SearchRequest
    
    if err := c.BodyParser(&req); err != nil {
        h.logger.Error().Err(err).Msg("Failed to parse request body")
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Invalid request body",
        })
    }

    if err := h.validator.Struct(req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": err.Error(),
        })
    }

    searchID, err := h.searchService.InitiateSearch(c.Context(), req)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to initiate search",
        })
    }

    return c.JSON(models.SearchResponse{
        Success: true,
        Message: "Search request submitted",
        Data: struct {
            SearchID string `json:"search_id"`
            Status   string `json:"status"`
        }{
            SearchID: searchID,
            Status:   "processing",
        },
    })
}

func (h *SearchHandler) StreamResults(c *fiber.Ctx) error {
    searchID := c.Params("search_id")
    if searchID == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Search ID is required",
        })
    }

    return h.sseService.StreamHandler(c, searchID)
}