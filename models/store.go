package models

import "time"

type Store struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	OwnerID   int       `json:"owner_id"` // liên kết với User
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}
