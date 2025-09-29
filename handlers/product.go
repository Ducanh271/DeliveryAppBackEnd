package handlers

import (
	"context"
	"database/sql"
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"example.com/delivery-app/config"
	"example.com/delivery-app/models"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

type NewProductRequest struct {
	Name        string  `form:"name" binding:"required"`
	Description string  `form:"description" binding:"required"`
	Price       float64 `form:"price" binding:"required,gt=0"`
	QtyInitial  int64   `form:"qty_initial" binding:"required,gte=0"`
	QtySold     int64   `form:"qty_sold" binding:"gte=0"`
}

// upload image to cloudinary
func uploadToCloudinary(file multipart.File, fileHeader *multipart.FileHeader) (string, string, error) {
	cld, err := cloudinary.NewFromURL(config.CloudinaryURL)
	if err != nil {
		return "", "", err
	}
	ctx := context.Background()
	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: fileHeader.Filename,
		Folder:   "product",
	})
	if err != nil {
		return "", "", err
	}
	return uploadResult.SecureURL, uploadResult.PublicID, nil
}
func CreateNewProductHandler(c *gin.Context, db *sql.DB) {
	var req NewProductRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can't read the file image"})
		return
	}
	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phải có ít nhất 1 ảnh"})
		return
	}

	// Bắt đầu transaction
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't start transaction"})
		return
	}
	defer tx.Rollback() // rollback nếu chưa commit

	// Tạo sản phẩm
	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		QtyInitial:  req.QtyInitial,
		QtySold:     req.QtySold,
		CreatedAt:   time.Now(),
	}

	product.ID, err = models.CreateProductTx(tx, &product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't create product"})
		return
	}

	// Upload ảnh
	var urls []string
	isFirst := true
	for _, fileHeader := range files {
		openFile, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't open image"})
			return
		}

		url, publicID, err := uploadToCloudinary(openFile, fileHeader)
		openFile.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't upload file image"})
			return
		}
		urls = append(urls, url)

		_, err = models.AddProductImageTx(tx, product.ID, url, publicID, isFirst)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		isFirst = false
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Created new product successfully",
		"product": product,
		"images":  urls,
	})
}

// func get products
func GetProductsHandler(c *gin.Context, db *sql.DB) {
	// Lấy query param
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	products, total, err := models.GetProductsPaginated(db, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}
func GetProductByIDHandler(c *gin.Context, db *sql.DB) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := models.GetProductByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	images, err := models.GetImagesByProductID(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load product images"})
		return
	}
	avgRate, count, err := models.GetRatingByProductID(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get rate"})
		return
	}
	productRes := models.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		QtyInitial:  product.QtyInitial,
		QtySold:     product.QtySold,
		CreatedAt:   product.CreatedAt,
		Images:      images,
		AvgRate:     avgRate,
		ReviewCount: count,
	}

	c.JSON(http.StatusOK, gin.H{
		"product": productRes,
	})
}
func DeleteProductHandler(c *gin.Context, db *sql.DB, cld *cloudinary.Cloudinary) {
	//  Lấy productID từ URL
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	//  Bắt đầu transaction
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}

	//  Xóa ảnh trên Cloudinary
	err = models.DeleteProductImages(cld, tx, productID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete images", "details": err.Error()})
		return
	}

	//  Xóa bản ghi trong bảng Products (các bảng con có cascade delete)
	_, err = models.DeleteProductByID(tx, productID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product", "details": err.Error()})
		return
	}

	//  Commit transaction
	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	//  Trả về response thành công
	c.JSON(http.StatusOK, gin.H{
		"message":    "product deleted successfully",
		"product_id": productID,
	})
}
