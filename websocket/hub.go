package websocket

import (
	"encoding/json"
	"sync"
)

type Hub struct {
	Clients    map[int64]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[int64]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.ID] = client
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			h.mu.RLock()
			for _, client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client.ID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// gửi thông báo tới 1 user cụ thể
func (h *Hub) SendToUser(userID int64, msg *Message) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	client, ok := h.Clients[userID]
	if !ok {
		return nil
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	client.Send <- data
	return nil
}

// xử lý message mà client gửi lên (ví dụ chat)
func (h *Hub) HandleMessage(sender *Client, msg *Message) {
	switch msg.Type {
	case "chat_message":
		// giả sử Data có: { "to_user": 5, "content": "hi" }
		dataMap := msg.Data.(map[string]interface{})
		toID := int64(dataMap["to_user"].(float64))
		h.SendToUser(toID, msg)
	default:
		// broadcast cho tất cả
		raw, _ := json.Marshal(msg)
		h.Broadcast <- raw
	}
}
