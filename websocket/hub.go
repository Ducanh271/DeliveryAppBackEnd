package websocket

import (
	"database/sql"
	"encoding/json"
	// "fmt"
	"log"
	"sync"
	"time"

	"example.com/delivery-app/models"
)

type Hub struct {
	DB         *sql.DB
	Clients    map[int64]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	mu         sync.RWMutex
}

func NewHub(db *sql.DB) *Hub {
	return &Hub{
		DB:         db,
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
		if msg.ToUserID == 0 || msg.ToUserID == sender.ID {
			return
		}
		log.Printf("%v", sender.ID)
		msg.CreatedAt = time.Now()
		messageModel := &models.Message{
			OrderID:    msg.OrderID,
			FromUserID: sender.ID,
			ToUserID:   msg.ToUserID,
			Content:    msg.Content,
			CreatedAt:  time.Now(),
		}
		message := &Message{
			Type:       msg.Type,
			OrderID:    msg.OrderID,
			FromUserID: sender.ID,
			ToUserID:   msg.ToUserID,
			Content:    msg.Content,
			CreatedAt:  time.Now(),
		}
		if err := models.SaveMessage(h.DB, messageModel); err != nil {
			log.Printf("Failed to save message: %v", err)
		}
		if err := h.SendToUser(msg.ToUserID, message); err != nil {
			log.Printf("Receiver is not online: %v", err)
		}
	case "location_update":
		if msg.OrderID == 0 {
			return
		}
		customerID, err := models.GetUserIDFromOrderID(h.DB, msg.OrderID)
		if err != nil {
			log.Printf("Can't find this order: %d, %v", msg.OrderID, err)
			return
		}
		// Gửi vị trí cho customer
		locationMsg := &Message{
			Type:      "location_update",
			OrderID:   msg.OrderID,
			Content:   "Update shipper's location",
			Latitude:  msg.Latitude,
			Longitude: msg.Longitude,
			CreatedAt: time.Now(),
		}

		if err := h.SendToUser(customerID, locationMsg); err != nil {
			log.Printf("Customer %d not online: %v", customerID, err)
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
