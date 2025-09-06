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
	ID     int64
	URL    string
	IsMain bool
}
type ProductResponse struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       float64        `json:"price"`
	QtyInitial  int64          `json:"qty_initial"`
	QtySold     int64          `json:"qty_sold"`
	CreatedAt   time.Time      `json:"created_at"`
	Images      []ProductImage `json:"images"`
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

// get products
func GetProductsPaginated(db *sql.DB, page, limit int) ([]ProductResponse, int, error) {
	offset := (page - 1) * limit

	// Query lấy sản phẩm + ảnh
	query := `
        SELECT p.id, p.name, p.description, p.price, p.qty_initial, p.qty_sold, p.created_at,
               i.id, i.url, pi.is_main
        FROM Products p
        LEFT JOIN ProductImages pi ON p.id = pi.product_id
        LEFT JOIN Images i ON pi.image_id = i.id
        ORDER BY p.created_at DESC
        LIMIT ? OFFSET ?
    `
	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	productsMap := make(map[int64]*ProductResponse)
	for rows.Next() {
		var (
			p      ProductResponse
			imgID  sql.NullInt64
			imgURL sql.NullString
			isMain sql.NullBool
		)
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price,
			&p.QtyInitial, &p.QtySold, &p.CreatedAt,
			&imgID, &imgURL, &isMain,
		)
		if err != nil {
			return nil, 0, err
		}

		existing, ok := productsMap[p.ID]
		if !ok {
			existing = &p
			existing.Images = []ProductImage{}
			productsMap[p.ID] = existing
		}

		if imgID.Valid {
			existing.Images = append(existing.Images, ProductImage{
				ID:     imgID.Int64,
				URL:    imgURL.String,
				IsMain: isMain.Bool,
			})
		}
	}

	products := make([]ProductResponse, 0, len(productsMap))
	for _, v := range productsMap {
		products = append(products, *v)
	}

	// Query tổng số sản phẩm
	var total int
	err = db.QueryRow("SELECT COUNT(*) FROM Products").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
func GetProductByID(db *sql.DB, id int64) (*Product, error) {
	query := `
        SELECT id, name, description, price, qty_initial, qty_sold, created_at
        FROM Products
        WHERE id = ?
    `
	row := db.QueryRow(query, id)

	var p Product
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.QtyInitial, &p.QtySold, &p.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func GetImagesByProductID(db *sql.DB, productID int64) ([]ProductImage, error) {
	query := `
        SELECT i.id, i.url, pi.is_main
        FROM Images i
        INNER JOIN ProductImages pi ON i.id = pi.image_id
        WHERE pi.product_id = ?
    `
	rows, err := db.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []ProductImage
	for rows.Next() {
		var img ProductImage
		err := rows.Scan(&img.ID, &img.URL, &img.IsMain)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}
