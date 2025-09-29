package main

import (
	"example.com/delivery-app/config"
	"example.com/delivery-app/database"
	"example.com/delivery-app/routes"
	"github.com/cloudinary/cloudinary-go/v2"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Init DB
	config.LoadConfig()
	database.InitDB()
	defer database.DB.Close()
	if err := database.CreateDefaultAdmin(database.DB); err != nil {
		log.Fatal("Error seeding admin:", err)
	}

	// Tạo Gin engine
	r := gin.Default()
	// create cloudinary
	cld, err := cloudinary.NewFromURL(config.CloudinaryURL)
	if err != nil {
		log.Fatal("Failed to connect to Cloudinary")
	}
	// Middleware CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Setup routes (truyền DB vào nếu cần)
	routes.SetupRoutes(r, database.DB, cld)

	// Run server
	r.Run(":8080")
}
