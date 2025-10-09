package websocket

import (
	"time"
)

type Message struct {
	Type       string    `json:"type"`
	OrderID    int64     `json:"order_id"`
	FromUserID int64     `json:"from_user_id,omitempty"`
	ToUserID   int64     `json:"to_user_id,omitempty"`
	Content    string    `json:"content"`
	Latitude   float64   `json:"latitude,omitempty"`
	Longitude  float64   `json:"longitude,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}

// type Message struct {
// 	Type       string                 `json:"type"`
// 	Data       map[string]interface{} `json:"data"`
// 	FromUserID int64                  `json:"from_user_id,omitempty"`
// 	ToUserID   int64                  `json:"to_user_id, omitempty"`
// }
