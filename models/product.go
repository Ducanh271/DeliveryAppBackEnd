package models

import (
	"database/sql"
	"time"
)

type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	QtyInitial  int64     `json:"qty_initial"`
	QtySold     int64     `json:"qty_sold"`
	CreatedAt   time.Time `json:"created_at"`
}
type ProductImage struct {
	ProductId int
	ImageID   int
	IsMain    bool
}

// create new products
func CreateProductTx(tx *sql.Tx, p *Product) (int64, error) {
	query := `
        INSERT INTO Products (name, description, price, qty_initial, qty_sold, created_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `
	p.CreatedAt = time.Now()
	result, err := tx.Exec(query, p.Name, p.Description, p.Price, p.QtyInitial, p.QtySold, p.CreatedAt)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// add image
func AddProductImageTx(tx *sql.Tx, productID int64, imageURL string, isMain bool) (int64, error) {
	// 1. Insert ảnh vào bảng Images
	queryImg := `INSERT INTO Images (url) VALUES (?)`
	result, err := tx.Exec(queryImg, imageURL)
	if err != nil {
		return 0, err
	}
	imageID, _ := result.LastInsertId()

	// 2. Insert mapping vào ProductImages
	queryMap := `INSERT INTO ProductImages (product_id, image_id, is_main) VALUES (?, ?, ?)`
	_, err = tx.Exec(queryMap, productID, imageID, isMain)
	if err != nil {
		return 0, err
	}

	return imageID, nil
}
