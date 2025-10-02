package models

import (
	"context"
	"database/sql"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
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
	AvgRate     float64        `json:"avg_rate"`
	ReviewCount int            `json:"review_count"`
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
func AddProductImageTx(tx *sql.Tx, productID int64, imageURL string, publicID string, isMain bool) (int64, error) {
	// 1. Insert ảnh vào bảng Images
	queryImg := `INSERT INTO Images (url, public_id) VALUES (?, ?)`
	result, err := tx.Exec(queryImg, imageURL, productID)
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

// get rating by product id
func GetRatingByProductID(db *sql.DB, productID int64) (float64, int, error) {
	var avgRate sql.NullFloat64
	var count int
	query := `select avg(rate), count(*) from Reviews where product_id = ?`
	err := db.QueryRow(query, productID).Scan(&avgRate, &count)
	if err != nil {
		return 0, 0, err
	}
	if !avgRate.Valid {
		return 0, count, nil
	}
	return avgRate.Float64, count, nil
}

// get products
func GetProductsPaginated(db *sql.DB, page, limit int) ([]ProductResponse, int, error) {
	offset := (page - 1) * limit

	// Query lấy sản phẩm + ảnh
	query := `
    SELECT p.id, p.name, p.description, p.price, p.qty_initial, p.qty_sold, p.created_at,
           i.id, i.url, pi.is_main
    FROM (
        SELECT * FROM Products
        ORDER BY id DESC
        LIMIT ? OFFSET ?
    ) p
    LEFT JOIN ProductImages pi ON p.id = pi.product_id
    LEFT JOIN Images i ON pi.image_id = i.id
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
		p.AvgRate, p.ReviewCount, err = GetRatingByProductID(db, p.ID)
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

// func get price of product by productID
func GetPriceProduct(db *sql.DB, productID int64) float64 {
	query := "select price from Products where id = ?"
	var price float64
	err := db.QueryRow(query, productID).Scan(&price)
	if err == sql.ErrNoRows {
		return 0
	}
	if err != nil {
		return 0
	}
	return price
}
func GetImageIDByProductID(db *sql.DB, productID int64) (error, int64) {
	query := "select image_id from ProductImages where product_id = ? and is_main = true"
	var imageID int64
	err := db.QueryRow(query, productID).Scan(&imageID)
	if err == sql.ErrNoRows {
		return nil, 0
	}
	if err != nil {
		return err, 0
	}
	return nil, imageID
}

func DeleteProductByID(tx *sql.Tx, productID int64) (int64, error) {
	query := "DELETE FROM Products WHERE id = ?"
	result, err := tx.Exec(query, productID)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	if rowsAffected == 0 {
		return 0, sql.ErrNoRows // Không tìm thấy sản phẩm để xóa
	}

	return rowsAffected, nil
}

func DeleteProductImages(cld *cloudinary.Cloudinary, tx *sql.Tx, productID int64) error {
	rows, err := tx.Query(`
        SELECT i.public_id 
        FROM Images i 
        JOIN ProductImages pi ON i.id = pi.image_id 
        WHERE pi.product_id = ?`, productID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var publicID string
		if err := rows.Scan(&publicID); err != nil {
			return err
		}

		_, err := cld.Upload.Destroy(context.Background(), uploader.DestroyParams{
			PublicID:   publicID,
			Invalidate: api.Bool(true),
		})
		if err != nil {
			return err
		}
	}

	// Kiểm tra lỗi duyệt rows
	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
func SearchProductsPaginated(db *sql.DB, keyword string, page, limit int) ([]ProductResponse, int, error) {
	offset := (page - 1) * limit
	searchTerm := "%" + keyword + "%"

	query := `
    SELECT p.id, p.name, p.description, p.price, p.qty_initial, p.qty_sold, p.created_at,
           i.id, i.url, pi.is_main
    FROM (
        SELECT * FROM Products
        WHERE name LIKE ? OR description LIKE ?
        ORDER BY id DESC
        LIMIT ? OFFSET ?
    ) p
    LEFT JOIN ProductImages pi ON p.id = pi.product_id
    LEFT JOIN Images i ON pi.image_id = i.id
	`

	rows, err := db.Query(query, searchTerm, searchTerm, limit, offset)
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

		//  Giữ nguyên cách tính rating
		p.AvgRate, p.ReviewCount, err = GetRatingByProductID(db, p.ID)
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

	//  Gom map -> slice
	products := make([]ProductResponse, 0, len(productsMap))
	for _, v := range productsMap {
		products = append(products, *v)
	}

	//  Lấy tổng số sản phẩm phù hợp điều kiện tìm kiếm
	var total int
	countQuery := `SELECT COUNT(*) FROM Products WHERE name LIKE ? OR description LIKE ?`
	err = db.QueryRow(countQuery, searchTerm, searchTerm).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
