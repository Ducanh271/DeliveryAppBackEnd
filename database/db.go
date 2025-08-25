package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	dsn := "root:Zxc13sdw@@tcp(127.0.0.1:3306)/DeliveryAppDB?parseTime=true"

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
