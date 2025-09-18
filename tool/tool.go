package main

import (
	"context"
	"database/sql"
	"example.com/delivery-app/config"
	"example.com/delivery-app/database"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xuri/excelize/v2"
	"log"
	"strconv"
	"time"
)

// Structs
type ProductImg struct {
	Url string
}

type Product struct {
	Name        string
	Description string
	Price       float64
	QtyInitial  int64
	QtySold     int64
	CreatedAt   time.Time
	ProductImgs []ProductImg
}

// Đọc từ Excel
func ReadProductsFromExcel(filename string) ([]Product, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows("Products")
	if err != nil {
		return nil, err
	}

	var products []Product

	// bỏ qua header (dòng 0)
	for i := 1; i < len(rows); i++ {
		row := rows[i]

		if len(row) < 6 {
			continue // skip nếu thiếu cột
		}

		price, _ := strconv.ParseFloat(row[2], 64)
		qtyInitial, _ := strconv.Atoi(row[3])
		qtySold, _ := strconv.Atoi(row[4])
		createdAt := time.Now() // tùy format trong file

		imgs := []ProductImg{}
		for j := 6; j < len(row); j++ {
			if row[j] != "" {
				imgs = append(imgs, ProductImg{Url: row[j]})
			}
		}

		product := Product{
			Name:        row[0],
			Description: row[1],
			Price:       price,
			QtyInitial:  int64(qtyInitial),
			QtySold:     int64(qtySold),
			CreatedAt:   createdAt,
			ProductImgs: imgs,
		}

		products = append(products, product)
	}

	return products, nil
}

// Insert với 1 transaction duy nhất
func InsertProducts(db *sql.DB, products []Product) error {
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}

	for _, p := range products {
		// Insert Product
		productQuery := `INSERT INTO Products (name, description, price, qty_initial, qty_sold, created_at) 
		                 VALUES (?, ?, ?, ?, ?, ?)`
		result, err := tx.ExecContext(ctx, productQuery, p.Name, p.Description, p.Price, p.QtyInitial, p.QtySold, p.CreatedAt)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("insert product failed: %w", err)
		}

		productID, err := result.LastInsertId()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("get product id failed: %w", err)
		}

		// Insert Images + ProductImages
		for j, img := range p.ProductImgs {
			imageQuery := `INSERT INTO Images (url) VALUES (?)`
			resImg, err := tx.ExecContext(ctx, imageQuery, img.Url)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("insert image failed: %w", err)
			}

			imageID, err := resImg.LastInsertId()
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("get image id failed: %w", err)
			}

			isMain := 0
			if j == 0 {
				isMain = 1
			}
			linkQuery := `INSERT INTO ProductImages (product_id, image_id, is_main) VALUES (?, ?, ?)`
			_, err = tx.ExecContext(ctx, linkQuery, productID, imageID, isMain)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("insert product_image failed: %w", err)
			}
		}
	}

	// commit nếu tất cả đều ok
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx failed: %w", err)
	}

	return nil
}

func main() {
	config.LoadConfig()
	database.InitDB()
	defer database.DB.Close()
	if err := database.CreateDefaultAdmin(database.DB); err != nil {
		log.Fatal("Error seeding admin:", err)
	}

	// đọc excel
	products, err := ReadProductsFromExcel("products_fixed.xlsx")
	if err != nil {
		panic(err)
	}

	// insert vào DB
	if err := InsertProducts(database.DB, products); err != nil {
		panic(err)
	}

	fmt.Println("Imported products successfully with single transaction!")
}
