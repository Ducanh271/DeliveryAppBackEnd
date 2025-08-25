package models

import "time"

type Order struct {
	ID         int       `json:"id"`
	CustomerID int       `json:"customer_id"`
	StoreID    int       `json:"store_id"`
	DriverID   int       `json:"driver_id,omitempty"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"` // pending | accepted | delivering | completed | cancelled
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
