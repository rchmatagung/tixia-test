package models

type SearchRequest struct {
    SearchID   string `json:"search_id"`
    From       string `json:"from"`
    To         string `json:"to"`
    Date       string `json:"date"`
    Passengers int    `json:"passengers"`
}

type SearchResult struct {
    SearchID string   `json:"search_id"`
    Status   string   `json:"status"`
    Results  []Flight `json:"results"`
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