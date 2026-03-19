package models

import "time"

type Order struct {
	OrderID  int       `json:"order_id"`
	Product  string    `json:"product"`
	Amount   float64   `json:"amount"`
	Status   string    `json:"status"`
	PlacedAt time.Time `json:"placed_at"`
}
