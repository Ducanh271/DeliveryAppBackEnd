package database

import (
	"database/sql"
	"example.com/delivery-app/config"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.DBUser,
		config.DBPass,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Database connection failed:", err)
	}

	// Kiểm tra kết nối
	if err := DB.Ping(); err != nil {
		log.Fatal("❌ Database not reachable:", err)
	}

	fmt.Println("✅ Connected to MySQL!")
}
