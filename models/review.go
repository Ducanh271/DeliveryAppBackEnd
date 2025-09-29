package models

import (
	"database/sql"
	"time"
)

// struct
type Review struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	UserID    int64     `json:"user_id"`
	OrderID   int64     `json:"order_id"`
	Rate      int8      `json:"rate"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type ReviewImage struct {
	ImageID int    `json:"image_id"`
	Url     string `json:"url"`
}
type ReviewForRes struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	UserID    int64     `json:"user_id"`
	UserName  string    `json:"user_name"`
	OrderID   int64     `json:"order_id"`
	Rate      int8      `json:"rate"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
type ReviewResWithTotal struct {
	Reviews    []ReviewForRes
	TotalCount int
}

// func check user can review
func CanUserReview(db *sql.DB, userID, productID int) (bool, error) {
	query := `
        SELECT COUNT(*) 
        FROM orders o
        JOIN order_items oi ON o.id = oi.order_id
        WHERE o.user_id = ? 
          AND oi.product_id = ?
          AND o.payment_status = 'paid'
          AND o.order_status = 'delivered'
    `

	var count int
	err := db.QueryRow(query, userID, productID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// func create reviews
func CreateReviewImagesTx(tx *sql.Tx, reviewID int64, imageURL string, publicID string) (int64, error) {
	// 1. Insert ảnh vào bảng Images
	queryImg := `INSERT INTO Images (url, public_id) VALUES (?, ?)`
	result, err := tx.Exec(queryImg, imageURL)
	if err != nil {
		return 0, err
	}
	imageID, _ := result.LastInsertId()

	// 2. Insert mapping vào ProductImages
	queryMap := `INSERT INTO ReviewImages (review_id, image_id) VALUES (?, ?)`
	_, err = tx.Exec(queryMap, reviewID, imageID)
	if err != nil {
		return 0, err
	}
	return imageID, nil
}

func CreateReviewTx(tx *sql.Tx, review *Review) (int64, error) {
	query := `insert into Reviews (product_id, user_id, order_id, rate, content, created_at) values(?, ?, ?, ?, ?, ?)`
	result, err := tx.Exec(query, review.ProductID, review.UserID, review.OrderID, review.Rate, review.Content, review.CreatedAt)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetReviewByProductID(db *sql.DB, productID int64, limit int, offset int) (*ReviewResWithTotal, error) {
	query := `select r.id, r.product_id, r.user_id, u.name, r.order_id, r.rate, r.content, r.created_at from Reviews r left join users u on r.user_id = u.id where product_id = ? order by r.created_at desc limit ? offset ?`
	rows, err := db.Query(query, productID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ListReview []ReviewForRes
	for rows.Next() {
		var review ReviewForRes
		err := rows.Scan(&review.ID, &review.ProductID, &review.UserID, &review.UserName, &review.OrderID, &review.Rate, &review.Content, &review.CreatedAt)
		if err != nil {
			return nil, err
		}
		ListReview = append(ListReview, review)
	}
	// Query tổng số review
	countQuery := `SELECT COUNT(*) FROM Reviews WHERE product_id = ?`
	var total int
	if err := db.QueryRow(countQuery, productID).Scan(&total); err != nil {
		return nil, err
	}
	return &ReviewResWithTotal{
		Reviews:    ListReview,
		TotalCount: total,
	}, nil

}
func GetReviewImage(db *sql.DB, reviewID int64) ([]ReviewImage, error) {
	query := `select i.id, i.url from Images i join ReviewImages ri on i.id = ri.image_id where ri.review_id = ?`
	rows, err := db.Query(query, reviewID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ListImage []ReviewImage
	for rows.Next() {
		var reviewImage ReviewImage
		err := rows.Scan(&reviewImage.ImageID, &reviewImage.Url)
		if err != nil {
			return nil, err
		}
		ListImage = append(ListImage, reviewImage)
	}
	return ListImage, nil
}
