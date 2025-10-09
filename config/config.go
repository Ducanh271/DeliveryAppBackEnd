package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EmailConfig struct {
	From     string
	Password string
	Host     string
	Port     string
}

var (
	CloudinaryURL string
	Email         EmailConfig
	JWTSecret     string
	DBUser        string
	DBPass        string
	DBHost        string
	DBPort        string
	DBName        string
)

func LoadConfig() {
	// Load .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("⚠️  Không tìm thấy file .env, dùng biến môi trường có sẵn")
	}
	Email = EmailConfig{
		From:     os.Getenv("EMAIL_FROM"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
	}
	//cloudinary url
	CloudinaryURL = os.Getenv("CLOUDINARY_URL")
	fmt.Println(CloudinaryURL)
	// Gán giá trị
	JWTSecret = os.Getenv("JWT_SECRET")
	DBUser = os.Getenv("DB_USER")
	DBPass = os.Getenv("DB_PASS")
	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBName = os.Getenv("DB_NAME")

	if JWTSecret == "" {
		log.Fatal("❌ JWT_SECRET chưa được set trong .env")
	}
}
