package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Order struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	ShipperID     int64     `json:"shipper_id"`
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
type OrderForShipper struct {
	OrderID int64
	Address string
}
type GetOrdersForShipperRes struct {
	Orders []Order
}

type OrderItemDetailResp struct {
	ProductID    int64   `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductImage string  `json:"product_image"`
	Quantity     int64   `json:"quantity"`
	Price        float64 `json:"price"`
	Subtotal     float64 `json:"subtotal"`
}
type ReceiveOrderRequest struct {
	OrderID int64 `json:"order_id"`
}
type UpdateOrderRequest struct {
	OrderID       int64  `json:"order_id"`
	PaymentStatus string `json:"payment_status"`
	OrderStatus   string `json:"order_status"`
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

// func GetAllOrder by admin
func GetAllOrders(db *sql.DB, page, limit int) ([]OrderSummaryResponse, int, error) {
	offset := (page - 1) * limit
	var total int
	err := db.QueryRow("select count(*) from orders").Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	query := `SELECT o.id, o.user_id, o.payment_status, o.order_status,
		       o.latitude, o.longitude, o.total_amount, 
		       o.thumbnail_id, o.created_at, o.updated_at,
		       i.url AS thumbnail
		FROM orders o
		LEFT JOIN Images i ON o.thumbnail_id = i.id
		ORDER BY o.id DESC limit ? offset ?`

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
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
			return nil, 0, err
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
		return nil, 0, err
	}

	return orders, total, nil
}

// get all order by shipper
func GetOrdersByShipper(db *sql.DB, page, limit int) ([]OrderSummaryResponse, int, error) {
	offset := (page - 1) * limit
	var total int
	err := db.QueryRow("select count(*) from orders").Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	query := `SELECT o.id, o.payment_status, o.order_status, 
		       o.latitude, o.longitude, o.total_amount, 
		       o.thumbnail_id, o.created_at, o.updated_at,
		       i.url AS thumbnail
		FROM orders o
		LEFT JOIN Images i ON o.thumbnail_id = i.id
		where o.order_status = "processing"
		ORDER BY o.id DESC limit ? offset ?`

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []OrderSummaryResponse
	for rows.Next() {
		var order Order
		var thumbnail sql.NullString

		err := rows.Scan(
			&order.ID,
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
			return nil, 0, err
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
		return nil, 0, err
	}

	return orders, total, nil
}
func GetOrderByID(db *sql.DB, orderID int64) (*Order, error) {
	query := `
		SELECT id, user_id , payment_status, order_status,
		       latitude, longitude, total_amount, thumbnail_id, created_at, updated_at
		FROM orders
		WHERE id = ?
	`

	row := db.QueryRow(query, orderID)

	var o Order
	err := row.Scan(
		&o.ID,
		&o.UserID,
		&o.PaymentStatus,
		&o.OrderStatus,
		&o.Latitude,
		&o.Longitude,
		&o.TotalAmount,
		&o.ThumbnailID,
		&o.CreatedAt,
		&o.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Không tìm thấy order
			return nil, nil
		}
		// Lỗi khác
		return nil, err
	}

	return &o, nil
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
func GetDetailOrder(db *sql.DB, orderID int64, userID int64, role string) (*GetOrderDetailResponse, error) {
	var order Order
	if role == "customer" || role == "shipper" {
		var exists bool
		err := db.QueryRow(
			"select exists (select 1 from orders where (id = ? and (user_id = ? or shipper_id = ?)))", orderID, userID, userID,
		).Scan(&exists)
		if err != nil {
			return nil, err
		}
	}

	orderQuery := `select id, user_id, payment_status, order_status, latitude, longitude, total_amount, thumbnail_id, created_at, updated_at from orders where id = ? and user_id`

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

// func for shipper

func CheckShipperOrder(db *sql.DB, shipperID int64, orderID int64) (bool, error) {
	query := "select shipper_id from orders where id = ?"
	var shipperIDdb int64
	err := db.QueryRow(query, orderID).Scan(&shipperIDdb)

	if err == sql.ErrNoRows {
		// không có order này trong DB
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if shipperID != shipperIDdb {
		return false, fmt.Errorf("shipper %d can't reach order %d", shipperID, orderID)
	}
	return true, nil
}

func UpdateStatusOrder(db *sql.DB, orderID int64, paymentStatus *string, orderStatus *string) error {
	query := "update orders set "
	args := []interface{}{}

	if paymentStatus != nil {
		query += "payment_status = ?"
		args = append(args, *paymentStatus)
	}
	if orderStatus != nil {
		if len(args) > 0 {
			query += ", "
		}
		query += "order_status = ?"
		args = append(args, *orderStatus)
	}

	query += " where id = ?"
	args = append(args, orderID)

	_, err := db.Exec(query, args...)
	return err
}

func UpdateShipperForOrder(db *sql.DB, orderID int64, shipperID int64) error {
	query := "update orders set shipper_id = ?, order_status = 'shipping' where id = ? and order_status= 'processing' "
	_, err := db.Exec(query, shipperID, orderID)
	return err
}

// func check current numbers of orders of shipper
func CheckNumberOfOrdersShipper(db *sql.DB, userID int64) (int, error) {
	var num int
	query := "select count(*) from orders where shipper_id = ? and order_status = 'shipping'"
	err := db.QueryRow(query, userID).Scan(&num)
	return num, err
}
