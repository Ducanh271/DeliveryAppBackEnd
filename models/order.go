package models

import (
	"database/sql"
	"time"
)

type Order struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	PaymentStatus string    `json:"payment_status"` // unpaid || paid || refund
	OrderStatus   string    `json:"order_status"`   // pending || processing ||shipped || delivered || canceled
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	TotalAmount   float64   `json:"total_amount"`
	ThumbnailID   int       `json:"thumbnail_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
type OrderItem struct {
	ID        int64   `json:"id"`
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int64   `json:"quantity"`
	Price     float64 `json:"price"`
}

type CreateOrderRequest struct {
	Latitude  float64                  `json:"latitude"`
	Longitude float64                  `json:"longitude"`
	Products  []CreateOrderItemRequest `json:"products"`
}
type CreateOrderItemRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int64 `json:"quantity"`
}
type OrderSummaryResponse struct {
	Order
	Thumbnail string `json:"thumbnail"` // lấy ảnh sản phẩm đầu tiên để hiển thị list
}
type OrdersOfUserResponse struct {
	Orders []OrderSummaryResponse `json:"orders"`
}
type GetOrderDetailResponse struct {
	Order      Order                 `json:"order"`
	OrderItems []OrderItemDetailResp `json:"items"`
}

type OrderItemDetailResp struct {
	ProductID    int64   `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductImage string  `json:"product_image"`
	Quantity     int64   `json:"quantity"`
	Price        float64 `json:"price"`
	Subtotal     float64 `json:"subtotal"`
}

func AddNewOrderToOrderTx(tx *sql.Tx, order *Order) (int64, error) {
	query := "insert into orders (user_id, payment_status, order_status, latitude, longitude, total_amount, thumbnail_id, created_at, updated_at) values (?, ?, ?, ?, ?, ?, ?,?, ?)"
	result, err := tx.Exec(query, order.UserID, order.PaymentStatus, order.OrderStatus, order.Latitude, order.Longitude, order.TotalAmount, order.ThumbnailID, order.CreatedAt, order.UpdatedAt)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
func AddNewOrderItemsTx(tx *sql.Tx, orderItem *OrderItem) error {
	query := "insert into order_items (order_id, product_id, quantity, price) values (?, ?, ?, ?)"
	_, err := tx.Exec(query, orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
	return err
}
func GetOrdersByUserID(db *sql.DB, userID int64) ([]OrderSummaryResponse, error) {
	query := `
		SELECT o.id, o.user_id, o.payment_status, o.order_status,
		       o.latitude, o.longitude, o.total_amount, 
		       o.thumbnail_id, o.created_at, o.updated_at,
		       i.url AS thumbnail
		FROM orders o
		LEFT JOIN Images i ON o.thumbnail_id = i.id
		WHERE o.user_id = ?
		ORDER BY o.id DESC
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []OrderSummaryResponse
	for rows.Next() {
		var order Order
		var thumbnail sql.NullString

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.PaymentStatus,
			&order.OrderStatus,
			&order.Latitude,
			&order.Longitude,
			&order.TotalAmount,
			&order.ThumbnailID,
			&order.CreatedAt,
			&order.UpdatedAt,
			&thumbnail,
		)
		if err != nil {
			return nil, err
		}

		resp := OrderSummaryResponse{
			Order:     order,
			Thumbnail: "",
		}
		if thumbnail.Valid {
			resp.Thumbnail = thumbnail.String
		}

		orders = append(orders, resp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func GetDetailOrder(db *sql.DB, orderID int64, userID int64) (*GetOrderDetailResponse, error) {
	var order Order
	orderQuery := `select id, user_id, payment_status, order_status, latitude, longitude, total_amount, thumbnail_id, created_at, updated_at from orders where id = ? and user_id = ?`
	err := db.QueryRow(orderQuery, orderID).Scan(
		&order.ID,
		&order.UserID,
		&order.PaymentStatus,
		&order.OrderStatus,
		&order.Latitude,
		&order.Longitude,
		&order.TotalAmount,
		&order.ThumbnailID,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	// --- Lấy chi tiết từng order_item ---
	itemQuery := `
		SELECT 
			o.product_id, 
			p.name, 
			o.quantity, 
			o.price,
			(SELECT url 
			 FROM Images i 
			 JOIN ProductImages pi ON i.id = pi.image_id 
			 WHERE pi.product_id = o.product_id AND pi.is_main = true 
			 LIMIT 1) AS image_url
		FROM order_items o
		JOIN Products p ON o.product_id = p.id
		WHERE o.order_id = ?`

	rows, err := db.Query(itemQuery, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderItemDetailResp
	for rows.Next() {
		var item OrderItemDetailResp
		err := rows.Scan(
			&item.ProductID,
			&item.ProductName,
			&item.Quantity,
			&item.Price,
			&item.ProductImage,
		)
		if err != nil {
			return nil, err
		}
		item.Subtotal = float64(item.Quantity) * item.Price
		items = append(items, item)
	}

	resp := &GetOrderDetailResponse{
		Order:      order,
		OrderItems: items,
	}
	return resp, nil
}
