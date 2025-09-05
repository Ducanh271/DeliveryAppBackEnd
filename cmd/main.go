package main

import (
	"example.com/delivery-app/config"
	"example.com/delivery-app/database"
	"example.com/delivery-app/routes"
	"log"

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

	// Setup routes (truyền DB vào nếu cần)
	routes.SetupRoutes(r, database.DB)

	// Run server
	r.Run(":8080")
}
