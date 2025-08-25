package database

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
)

func CreateDefaultAdmin(db *sql.DB) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		_, err = db.Exec(`
            INSERT INTO users (name, email, password, role, phone, address)
            VALUES (?, ?, ?, ?, ?, ?)`,
			"Admin", "admin@example.com", string(hashedPassword), "admin", "0000000000", "Admin Address")
		return err
	}
	return nil
}
