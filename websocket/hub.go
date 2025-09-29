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
	select {
	case client.Send <- data:
	default:
		close(client.Send)
		delete(h.Clients, userID)
		return nil
	}
	return nil
}

// xử lý message mà client gửi lên (ví dụ chat)
func (h *Hub) HandleMessage(sender *Client, msg *Message) {
	switch msg.Type {
	case "chat_message":
		// Giả sử Data có: { "to_user": 5, "content": "hi" }
		dataMap := msg.Data
		toID := int64(dataMap["to_user"].(float64))
		if toID == sender.ID {
			// Tránh gửi lại cho chính mình
			return
		}
		// Gửi tin nhắn đến người nhận
		msg.FromUserID = sender.ID // Thêm thông tin người gửi
		if err := h.SendToUser(toID, msg); err != nil {
			// Xử lý lỗi nếu cần (ví dụ: người nhận không online)
		}
	default:
		// Broadcast cho tất cả, trừ người gửi
		raw, _ := json.Marshal(msg)
		h.mu.RLock()
		for _, client := range h.Clients {
			if client.ID != sender.ID {
				select {
				case client.Send <- raw:
				default:
					close(client.Send)
					delete(h.Clients, client.ID)
				}
			}
		}
		h.mu.RUnlock()
	}
}
