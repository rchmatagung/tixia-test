package models

type SearchRequest struct {
    From       string `json:"from" validate:"required,len=3"`
    To         string `json:"to" validate:"required,len=3"`
    Date       string `json:"date" validate:"required"`
    Passengers int    `json:"passengers" validate:"required,min=1"`
}

type SearchResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Data    struct {
        SearchID string `json:"search_id"`
        Status   string `json:"status"`
    } `json:"data"`
}

type SearchResult struct {
    SearchID string        `json:"search_id"`
    Status   string        `json:"status"`
    Results  []Flight      `json:"results,omitempty"`
    Total    int           `json:"total_results,omitempty"`
}

type Flight struct {
	ID           string  `json:"id"`
    Airline      string  `json:"airline"`
    FlightNumber string  `json:"flight_number"`
    From         string  `json:"from"`
    To           string  `json:"to"`
    Departure    string  `json:"departure_time"`
    Arrival      string  `json:"arrival_time"`
    Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
	Available    bool    `json:"available"`
}