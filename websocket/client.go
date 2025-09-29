package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID   int64           // ID người dùng
	Conn *websocket.Conn // kết nối WebSocket
	Send chan []byte     // kênh gửi dữ liệu ra cho client
}

func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		// parse message
		var m Message
		if err := json.Unmarshal(msg, &m); err != nil {
			continue
		}

		// xử lý theo loại message (ví dụ chat hoặc gì khác)
		hub.HandleMessage(c, &m)
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		c.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}
