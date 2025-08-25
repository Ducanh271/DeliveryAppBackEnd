package models

import "time"

type Product struct {
	ID        int       `json:"id"`
	StoreID   int       `json:"store_id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
}
