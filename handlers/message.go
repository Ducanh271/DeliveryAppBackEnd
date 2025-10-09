package handlers

import (
	"database/sql"
	"example.com/delivery-app/models"
	// "example.com/delivery-app/websocket"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetMessageHandler(c *gin.Context, db *sql.DB) {
	oIDstring := c.Param("id")
	orderID, err := strconv.ParseInt(oIDstring, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can't convert to orderID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	isBelong, err := models.CheckOrderUser(db, userID.(int64), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check if order belongs to customer"})
		return
	}
	if !isBelong {
		isBelong, _ := models.CheckShipperOrder(db, userID.(int64), orderID)
		if !isBelong {
			c.JSON(http.StatusBadRequest, gin.H{"error": "This order is not belong to you"})
			return
		}

	}

	// Lấy query params
	limitStr := c.DefaultQuery("limit", "20") // mặc định 20 tin nhắn
	limit, _ := strconv.Atoi(limitStr)
	beforeStr := c.Query("before") // id hoặc timestamp của tin nhắn cuối cùng client đang có

	messages, err := models.GetMessagesByOrder(db, orderID, limit, beforeStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages of this order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}
