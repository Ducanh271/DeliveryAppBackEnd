package websocket

type Message struct {
	Type       string                 `json:"type"`
	Data       map[string]interface{} `json:"data"`
	FromUserID int64                  `json:"from_user_id,omitempty"`
}
