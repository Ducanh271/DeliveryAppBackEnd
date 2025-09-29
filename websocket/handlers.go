package websocket

import (
	"example.com/delivery-app/middleware"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"

	"example.com/delivery-app/config"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWs(hub *Hub, c *gin.Context) {
	// 1️⃣ Lấy token từ query param
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "please authenticate first"})
		return
	}

	// 🧩 Debug: log token và secret
	fmt.Println("🔹 [ServeWs] Incoming WebSocket connection")
	fmt.Println("🔹 [ServeWs] Token:", tokenStr)
	fmt.Println("🔹 [ServeWs] JWT Secret:", config.JWTSecret)

	// 2️⃣ Parse và xác thực JWT
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return middleware.JwtKey, nil
	})

	// 🧩 Debug: log lỗi parse nếu có
	if err != nil {
		fmt.Println("❌ [ServeWs] JWT parse error:", err)
	}
	if !token.Valid {
		fmt.Println("❌ [ServeWs] Token is invalid")
	}

	// 🧩 Debug: in thông tin exp, now
	if exp, ok := claims["exp"].(float64); ok {
		fmt.Println("🔹 [ServeWs] Token exp:", int64(exp))
		fmt.Println("🔹 [ServeWs] Current time:", time.Now().Unix())
		if int64(exp) < time.Now().Unix() {
			fmt.Println("⚠️  [ServeWs] Token has expired")
		}
	} else {
		fmt.Println("⚠️  [ServeWs] Token has no valid exp claim")
	}

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	// 3️⃣ Lấy userID từ claims
	userIDFloat, ok := claims["userID"].(float64)
	if !ok {
		fmt.Println("❌ [ServeWs] Token payload invalid (missing userID)")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token payload"})
		return
	}
	userID := int64(userIDFloat)

	// (Tuỳ chọn) Lấy role nếu bạn cần phân quyền shipper/admin
	role, _ := claims["role"].(string)

	// 🧩 Debug: in userID và role
	fmt.Printf("✅ [ServeWs] Token valid → userID=%d, role=%s\n", userID, role)

	// 4️⃣ Upgrade lên WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("❌ [ServeWs] WebSocket upgrade error:", err)
		return
	}

	// 5️⃣ Tạo client mới
	client := &Client{
		ID:   userID,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	hub.Register <- client

	fmt.Printf("🎉 [ServeWs] WebSocket connected: userID=%d, role=%s\n", userID, role)

	// 6️⃣ Start read/write goroutines
	go client.WritePump()
	go client.ReadPump(hub)
}

