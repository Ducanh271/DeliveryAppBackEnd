package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID                int64      `json:"id"`
	Name              string     `json:"name"`
	Email             string     `json:"email"`
	Password          string     `json:"-"`
	Phone             string     `json:"phone"`
	Address           string     `json:"address"`
	Role              string     `json:"role"`
	CreatedAt         time.Time  `json:"created_at"`
	OTPCode           *string    `json:"-"`
	OTPExpiresAt      *time.Time `json:"-"`
	IsVerified        bool       `json:"is_verified"`
	ResetOTP          *string    `json:"-"`
	ResetOTPExpiresAt *time.Time `json:"-"`
}

func CheckEmailExists(db *sql.DB, email string) (bool, error) {
	var exist int
	query := "select 1 from users where email = ? limit 1"
	err := db.QueryRow(query, email).Scan(&exist)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil

}
func CreateUserTx(tx *sql.Tx, user *User) (int64, error) {
	query := "insert into users (name, email, password, phone, address, role, created_at) values (?,?, ?, ?, ?, ?,?)"
	user.CreatedAt = time.Now()
	result, err := tx.Exec(query, user.Name, user.Email, user.Password, user.Phone, user.Address, user.Role, user.CreatedAt)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	query := "select id, name, email, password, phone, address, role, created_at, otp_code, otp_expires_at, is_verified, reset_otp, reset_otp_expires_at from users where email = ?"
	row := db.QueryRow(query, email)
	var user User
	var createdAtstr string

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.Address, &user.Role, &createdAtstr, &user.OTPCode, &user.OTPExpiresAt, &user.IsVerified, &user.ResetOTP, &user.ResetOTPExpiresAt)
	if err != nil {
		return nil, err
	}
	user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtstr)
	return &user, nil
}
func GetUserByID(db *sql.DB, userID int64) (*User, error) {
	query := "select id, name, email, password, phone, address, role, created_at from users where id = ?"
	row := db.QueryRow(query, userID)
	var user User
	var createdAtstr string

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.Address, &user.Role, &createdAtstr)
	if err != nil {
		return nil, err
	}
	user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtstr)
	return &user, nil
}
func UpdateOTPTx(tx *sql.Tx, userEmail string, otp string, expiry time.Time) error {
	query := `update users set otp_code = ?, otp_expires_at = ? where email = ?`
	_, err := tx.Exec(query, otp, expiry, userEmail)
	return err
}
func VerifyUser(db *sql.DB, userEmail string) error {
	updateQuery := `UPDATE users SET is_verified = true WHERE email = ?`
	_, err := db.Exec(updateQuery, userEmail)
	return err
}
func ClearOTP(db *sql.DB, userID int64) error {
	_, err := db.Exec(`UPDATE users SET otp_code = NULL, otp_expires_at = NULL WHERE id = ?`, userID)
	return err
}

func SetResetOTP(db *sql.DB, email, otp string, expiry time.Time) error {
	_, err := db.Exec(`
		UPDATE users
		SET reset_otp = ?, reset_otp_expires_at = ?
		WHERE email = ?`, otp, expiry, email)
	return err
}

func UpdatePasswordByEmail(db *sql.DB, email, hashed string) error {
	_, err := db.Exec(`UPDATE users SET password = ? WHERE email = ?`, hashed, email)
	return err
}

func ClearResetOTP(db *sql.DB, userID int64) error {
	_, err := db.Exec(`UPDATE users SET reset_otp = NULL, reset_otp_expires_at = NULL WHERE id = ?`, userID)
	return err
}
