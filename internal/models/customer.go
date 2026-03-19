package models

import "time"

type Customer struct {
	ID       int       `json:"_id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Orders   []Order   `json:"orders"`
	SyncedAt time.Time `json:"synced_at"`
}
