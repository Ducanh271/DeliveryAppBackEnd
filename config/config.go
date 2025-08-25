package config

import (
	"fmt"
	"os"
)

var (
	DBUser     = getEnv("DB_USER", "root")
	DBPassword = getEnv("DB_PASSWORD", "Zxc13sdw@")
	DBHost     = getEnv("DB_HOST", "127.0.0.1")
	DBPort     = getEnv("DB_PORT", "3306")
	DBName     = getEnv("DB_NAME", "DeliveryAppDB")
	JwtKey     = []byte("super_secret_key") // Để .env thì bảo mật hơn
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		DBUser, DBPassword, DBHost, DBPort, DBName)
}
