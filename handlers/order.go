package handlers

import (
	"database/sql"
	"example.com/delivery-app/models"
	"example.com/delivery-app/websocket"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

const MaxOrders = 10

type shipperInfoRes struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type OrdersShipperResponse struct {
	OrderID   int64   `json:"order_id"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}
type ReceiveOrderRequest struct {
	OrderID int64
}

func CreateOrderWithItems(db *sql.DB, order *models.Order, items []models.OrderItem) error {
	// bắt đầu transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// tạo order
	orderID, err := models.AddNewOrderToOrderTx(tx, order)
	if err != nil {
		tx.Rollback()
		return err
	}

	// tạo order_items
	for _, item := range items {
		item.OrderID = orderID
		if err := models.AddNewOrderItemsTx(tx, &item); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
func CreateOrderHandler(c *gin.Context, db *sql.DB) {
	var req models.CreateOrderRequest
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// chuyển từ OrderItemRequest -> OrderItem
	var items []models.OrderItem
	var totalAmount float64
	for _, p := range req.Products {
		items = append(items, models.OrderItem{
			ProductID: p.ProductID,
			Quantity:  p.Quantity,
			Price:     models.GetPriceProduct(db, p.ProductID), // giả sử bạn tính giá sau, hoặc join từ bảng products
		})
		totalAmount += float64(p.Quantity) * models.GetPriceProduct(db, p.ProductID)
	}
	firstProductID := req.Products[0].ProductID
	err, thumbnailID := models.GetImageIDByProductID(db, firstProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get image_id for thumbnail"})
		return
	}
	order := &models.Order{
		UserID:        userID.(int64),
		PaymentStatus: "unpaid",
		OrderStatus:   "pending",
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		TotalAmount:   totalAmount,
		ThumbnailID:   int(thumbnailID),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	// gọi transaction
	if err := CreateOrderWithItems(db, order, items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// trả về kết quả
	c.JSON(http.StatusOK, gin.H{
		"message": "order created successfully",
	})

}
func GetOrdersByUserIDHandler(c *gin.Context, db *sql.DB) {
	// lấy userID từ context (đã qua middleware auth)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// gọi models để lấy orders
	orders, err := models.GetOrdersByUserID(db, userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// response chuẩn
	resp := models.OrdersOfUserResponse{
		Orders: orders,
	}

	c.JSON(http.StatusOK, resp)
}
func GetOrderDetailHandler(c *gin.Context, db *sql.DB) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}

	// lấy orderID từ param
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	// gọi models để lấy dữ liệu
	orderDetail, err := models.GetDetailOrder(db, orderID, userID.(int64), role.(string))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// trả về response JSON
	c.JSON(http.StatusOK, orderDetail)
}

// func for shipper
func ReceiveOrderByShipperHandler(c *gin.Context, db *sql.DB, hub *websocket.Hub) {
	// Lấy userID từ context
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDVal.(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	// Bind request JSON
	var req models.ReceiveOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	orderID := req.OrderID

	// Lấy order theo ID
	order, err := models.GetOrderByID(db, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get order"})
		return
	}
	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// Kiểm tra số lượng đơn shipper đang nhận
	num, err := models.CheckNumberOfOrdersShipper(db, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check number of orders"})
		return
	}
	if num >= MaxOrders {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you have reached the maximum number of orders"})
		return
	}

	// Cập nhật shipper cho order
	if err := models.UpdateShipperForOrder(db, orderID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order"})
		return
	}

	// Gửi thông báo qua WebSocket
	// msg := websocket.Message{
	// 	Type: "shipped_order",
	// 	Data: map[string]interface{}{
	// 		"order_id":     orderID,
	// 		"order_status": "shipped",
	// 		"shipper":      userID,
	// 	},
	// }
	// hub.SendToUser(order.UserID, &msg)

	// Trả về response thành công
	c.JSON(http.StatusOK, gin.H{"message": "order received successfully"})
}

func UpdateOrderShipper(c *gin.Context, db *sql.DB) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req models.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	check, err := models.CheckShipperOrder(db, userID.(int64), req.OrderID)
	if err != nil || check == false {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	err = models.UpdateStatusOrder(db, int64(req.OrderID), &req.PaymentStatus, &req.OrderStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't update this order"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "update successfully"})

}
func AcceptOrderAdmin(c *gin.Context, db *sql.DB) {
	orderIDstr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	OrderStatus := "processing"

	err = models.UpdateStatusOrder(db, orderID, nil, &OrderStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't update this order"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "update successfully"})

}

// func cancle order by user
func CancleOrderByUserHandler(c *gin.Context, db *sql.DB) {
	orderIDstr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid orderID"})
		return
	}
	userIDstr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDstr.(int64)
	bool, err := models.CheckOrderUser(db, userID, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check order"})
		return
	}
	if bool == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This order is not your"})
		return
	}
	order, err := models.GetOrderByID(db, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get status of order"})
		return
	}
	if order.OrderStatus != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can't cancel this order, because it's not pending"})
		return
	}
	orderStatus := "cancelled"
	err = models.UpdateStatusOrder(db, orderID, nil, &orderStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update this order"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cancelled order successfully"})

}

// func get order by shipper
func GetOrdersByAdminHandler(c *gin.Context, db *sql.DB) {
	// Lấy query param
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	orders, total, err := models.GetAllOrders(db, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// get order by shipper
func GetOrdersByShipperHandler(c *gin.Context, db *sql.DB) {
	// Lấy query param
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	orders, total, err := models.GetOrdersByShipper(db, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// get received order by shipper
func GetReceivedOrdersByShipperHandler(c *gin.Context, db *sql.DB) {
	// Lấy query param
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	shipperIDstr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	shipperID := shipperIDstr.(int64)
	orders, total, err := models.GetReceivedOrdersByShipper(db, shipperID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}
func GetNumberOfOrderAndRevenueHandler(c *gin.Context, db *sql.DB) {
	num, revenue, err := models.GetNumberAndRevenueOfOrders(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get infor order"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"number of order": num, "revenue": revenue})
}

func GetShipperInfoByOrderIDHandler(c *gin.Context, db *sql.DB) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}
	userIDstr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDstr.(int64)
	bool, err := models.CheckOrderUser(db, userID, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check order"})
		return
	}
	if bool == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This order is not your"})
		return
	}
	var res shipperInfoRes
	res.ID, res.Name, res.Phone, err = models.GetShipperInfoFromOrderID(db, orderID)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "No shipper assigned to this order"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get shipper info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"shipper": res})
}
