package models

import (
	"database/sql"
	"time"
)

type Message struct {
	ID         int64     `json:"id"`
	OrderID    int64     `json:"order_id"`
	FromUserID int64     `json:"from_user_id"`
	ToUserID   int64     `json:"to_user_id"`
	Content    string    `json:"content"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

func SaveMessage(db *sql.DB, msg *Message) error {
	query := `insert into messages(order_id, sender_id, receiver_id, content, is_read) values (?, ?, ?, ?, ?)`
	result, err := db.Exec(query, msg.OrderID, msg.FromUserID, msg.ToUserID, msg.Content, false)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err == nil {
		msg.ID = id
	}
	msg.CreatedAt = time.Now()
	return nil
}

// GetMessagesByOrder lấy tất cả tin nhắn theo order_id
func GetMessagesByOrder(db *sql.DB, orderID int64) ([]Message, error) {
	rows, err := db.Query(`
		SELECT id, order_id, sender_id, receiver_id, content, is_read, created_at
		FROM messages
		WHERE order_id = ?
		ORDER BY created_at ASC
	`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.OrderID, &m.FromUserID, &m.ToUserID, &m.Content, &m.IsRead, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, nil
}

func GetUnreadCountByUserID(db *sql.DB, userID int64, orderID int64) (int, error) {
	var unRead int
	query := `select count(*) from messages where user_id = ? and order_id = ? and is_read = false`
	err := db.QueryRow(query, userID, orderID).Scan(&unRead)
	if err != nil {
		return 0, err
	}
	return unRead, nil
}
