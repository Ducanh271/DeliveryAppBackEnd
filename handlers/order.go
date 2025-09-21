package handlers

import (
	"database/sql"
	"example.com/delivery-app/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

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

	// lấy orderID từ param
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	// gọi models để lấy dữ liệu
	orderDetail, err := models.GetDetailOrder(db, orderID, userID.(int64))
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
	err = models.UpdateStatusOrder(db, int(req.OrderID), &req.PaymentStatus, &req.OrderStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't update this order"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "uodate successfully"})

}
