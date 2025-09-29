package main

import (
	"fmt"
	"os"
)

func main() {
	JWTSecret := os.Getenv("JWT_SECRET")
	fmt.Println("JWT_SECRET:", JWTSecret)
}
