package services

import (
	"encoding/json"
	"provider-service/internal/models"
)

var masterFlightsJSON = []byte(`
[
  {"id":"flight-uuid","airline":"Garuda Indonesia","flight_number":"GA123","from":"CGK","to":"DPS","departure_time":"2025-07-10 14:00","arrival_time":"2025-07-10 17:00","price":1000000,"currency":"IDR","available":true},
  {"id":"flight-uuid-2","airline":"Lion Air","flight_number":"JT123","from":"CGK","to":"SUB","departure_time":"2025-07-11 14:00","arrival_time":"2025-07-11 17:00","price":800000,"currency":"IDR","available":true},
  {"id":"flight-uuid-3","airline":"AirAsia Indonesia","flight_number":"QZ456","from":"CGK","to":"DPS","departure_time":"2025-07-12 08:30","arrival_time":"2025-07-12 11:30","price":750000,"currency":"IDR","available":true},
  {"id":"flight-uuid-4","airline":"Citilink","flight_number":"QG789","from":"CGK","to":"SUB","departure_time":"2025-07-13 16:45","arrival_time":"2025-07-13 19:45","price":650000,"currency":"IDR","available":true},
  {"id":"flight-uuid-5","airline":"Batik Air","flight_number":"ID234","from":"CGK","to":"DPS","departure_time":"2025-07-14 20:15","arrival_time":"2025-07-14 23:15","price":900000,"currency":"IDR","available":true},
  {"id":"flight-uuid-6","airline":"Sriwijaya Air","flight_number":"SJ567","from":"CGK","to":"SUB","departure_time":"2025-07-15 06:00","arrival_time":"2025-07-15 09:00","price":700000,"currency":"IDR","available":true},
  {"id":"flight-uuid-7","airline":"Garuda Indonesia","flight_number":"GA456","from":"DPS","to":"CGK","departure_time":"2025-07-16 10:30","arrival_time":"2025-07-16 13:30","price":1100000,"currency":"IDR","available":true},
  {"id":"flight-uuid-8","airline":"Lion Air","flight_number":"JT789","from":"SUB","to":"CGK","departure_time":"2025-07-17 12:15","arrival_time":"2025-07-17 15:15","price":850000,"currency":"IDR","available":true},
  {"id":"flight-uuid-9","airline":"AirAsia Indonesia","flight_number":"QZ123","from":"DPS","to":"SUB","departure_time":"2025-07-10 18:00","arrival_time":"2025-07-10 19:30","price":500000,"currency":"IDR","available":true},
  {"id":"flight-uuid-10","airline":"Citilink","flight_number":"QG456","from":"SUB","to":"DPS","departure_time":"2025-07-11 07:45","arrival_time":"2025-07-11 09:15","price":550000,"currency":"IDR","available":true},
  {"id":"flight-uuid-11","airline":"Batik Air","flight_number":"ID789","from":"CGK","to":"DPS","departure_time":"2025-07-12 22:30","arrival_time":"2025-07-13 01:30","price":950000,"currency":"IDR","available":true},
  {"id":"flight-uuid-12","airline":"Sriwijaya Air","flight_number":"SJ234","from":"DPS","to":"SUB","departure_time":"2025-07-13 14:30","arrival_time":"2025-07-13 16:00","price":600000,"currency":"IDR","available":true},
  {"id":"flight-uuid-13","airline":"Garuda Indonesia","flight_number":"GA789","from":"SUB","to":"DPS","departure_time":"2025-07-14 09:15","arrival_time":"2025-07-14 10:45","price":1200000,"currency":"IDR","available":true},
  {"id":"flight-uuid-14","airline":"Lion Air","flight_number":"JT456","from":"DPS","to":"CGK","departure_time":"2025-07-15 15:45","arrival_time":"2025-07-15 18:45","price":900000,"currency":"IDR","available":true},
  {"id":"flight-uuid-15","airline":"AirAsia Indonesia","flight_number":"QZ789","from":"SUB","to":"CGK","departure_time":"2025-07-16 11:30","arrival_time":"2025-07-16 14:30","price":700000,"currency":"IDR","available":true},
  {"id":"flight-uuid-16","airline":"Batik Air","flight_number":"ID567","from":"CGK","to":"SUB","departure_time":"2025-07-10 10:00","arrival_time":"2025-07-10 13:00","price":850000,"currency":"IDR","available":true},
  {"id":"flight-uuid-17","airline":"Sriwijaya Air","flight_number":"SJ890","from":"DPS","to":"CGK","departure_time":"2025-07-11 19:30","arrival_time":"2025-07-11 22:30","price":950000,"currency":"IDR","available":true},
  {"id":"flight-uuid-18","airline":"Citilink","flight_number":"QG234","from":"SUB","to":"DPS","departure_time":"2025-07-12 15:00","arrival_time":"2025-07-12 16:30","price":600000,"currency":"IDR","available":true},
  {"id":"flight-uuid-19","airline":"Lion Air","flight_number":"JT567","from":"CGK","to":"DPS","departure_time":"2025-07-13 08:00","arrival_time":"2025-07-13 11:00","price":800000,"currency":"IDR","available":true},
  {"id":"flight-uuid-20","airline":"AirAsia Indonesia","flight_number":"QZ890","from":"DPS","to":"SUB","departure_time":"2025-07-14 16:45","arrival_time":"2025-07-14 18:15","price":550000,"currency":"IDR","available":true},
  {"id":"flight-uuid-21","airline":"Garuda Indonesia","flight_number":"GA567","from":"SUB","to":"CGK","departure_time":"2025-07-15 13:00","arrival_time":"2025-07-15 16:00","price":1000000,"currency":"IDR","available":true},
  {"id":"flight-uuid-22","airline":"Batik Air","flight_number":"ID890","from":"DPS","to":"CGK","departure_time":"2025-07-16 07:15","arrival_time":"2025-07-16 10:15","price":900000,"currency":"IDR","available":true},
  {"id":"flight-uuid-23","airline":"Citilink","flight_number":"QG567","from":"CGK","to":"DPS","departure_time":"2025-07-17 09:00","arrival_time":"2025-07-17 12:00","price":750000,"currency":"IDR","available":true},
  {"id":"flight-uuid-24","airline":"Sriwijaya Air","flight_number":"SJ567","from":"SUB","to":"DPS","departure_time":"2025-07-17 20:00","arrival_time":"2025-07-17 21:30","price":650000,"currency":"IDR","available":true}
]
`)

type MockAPI struct{}

func NewMockAPI() *MockAPI {
	return &MockAPI{}
}

// func init() {
// 	_ = json.Unmarshal(masterFlightsJSON, &)
// }

func (m *MockAPI) SearchFlights(request models.SearchRequest) []models.Flight {
	var masterFlights []models.Flight
	_ = json.Unmarshal(masterFlightsJSON, &masterFlights)

	var candidates []models.Flight
	for _, f := range masterFlights {
		if f.From == request.From && f.To == request.To && extractDate(f.Departure) == request.Date {
			candidates = append(candidates, models.Flight{
				ID:           f.ID,
				Airline:      f.Airline,
				FlightNumber: f.FlightNumber,
				From:         f.From,
				To:           f.To,
				Departure:    f.Departure,
				Arrival:      f.Arrival,
				Price:        f.Price,
				Currency:     f.Currency,
				Available:    f.Available,
			})
		}
	}

	return candidates
}

func extractDate(ts string) string {
	// ts = "2025-07-10 14:00"
	if len(ts) >= 10 {
		return ts[:10] // "2025-07-10"
	}
	return ts
}
