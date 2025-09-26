package models

import (
	"database/sql"
	"log"
	"time"
)

type RefreshToken struct {
	ID        int       `json:"id"`
	UserID    int64     `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func SaveRefreshToken(db *sql.DB, userID int64, token string, expiresAt time.Time) error {
	query := "insert into refresh_tokens (user_id, token, expires_at) values (?, ?, ?)"
	_, err := db.Exec(query, userID, token, expiresAt)
	return err
}

func GetRefreshTokenByToken(db *sql.DB, token string) (*RefreshToken, error) {
	var rt RefreshToken
	query := "select id, user_id, token, expires_at, created_at from refresh_tokens where token = ?"
	err := db.QueryRow(query, token).Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func GetRefreshTokensByUserID(db *sql.DB, userID int64) ([]RefreshToken, error) {
	rows, err := db.Query(
		"SELECT token, expires_at FROM refresh_tokens WHERE user_id = ?",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []RefreshToken
	for rows.Next() {
		var token RefreshToken
		token.UserID = userID
		if err := rows.Scan(&token.Token, &token.ExpiresAt); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}
func CheckRefreshToken(db *sql.DB, token string) (bool, error) {
	query := `select 1 from refresh_tokens where token = ?`
	var exists int
	err := db.QueryRow(query, token).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
func DeleteRefreshToken(db *sql.DB, token string) error {
	query := `delete from refresh_tokens where token = ?`
	_, err := db.Exec(query, token)
	return err
}
func UpdateRefreshToken(db *sql.DB, old_token string, new_token string) error {
	query := `update refresh_tokens set token = ? where token = ?`
	_, err := db.Exec(query, new_token, old_token)
	return err
}
func DeleteUserRefreshToken(db *sql.DB, userID int64) error {
	query := `delete from refresh_tokens where user_id = ?`
	_, err := db.Exec(query, userID)
	return err
}
func StartTokenCleanUp(db *sql.DB) {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			_, err := db.Exec("delete from refresh_tokens where expires_at < NOW()")
			if err != nil {
				log.Println("Error cleaning expires tokens: ", err)
			} else {
				log.Println("Expires refresh tokens cleaned successfully")
			}
		}
	}()
}
