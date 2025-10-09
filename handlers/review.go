package handlers

import (
	"context"
	"database/sql"
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

type CreateReviewRequest struct {
	ProductID int64  `json:"product_id" binding:"required"`
	OrderID   int64  `json:"order_id" binding:"required"`
	Rate      int8   `json:"rate"`
	Content   string `json:"content"`
}

type ReviewInfor struct {
	UserID    int64     `json:"user_id"`
	UserName  string    `json:"user_name"`
	Rate      int8      `json:"rate"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Images    []ReviewImageRes
}
type ReviewImageRes struct {
	ImageID int64  `json:"image_id"`
	Url     string `json:"url"`
}

// func
// upload image to cloudinary
func uploadToCloudinaryForReview(file multipart.File, fileHeader *multipart.FileHeader) (string, string, error) {
	cld, err := cloudinary.NewFromURL(config.CloudinaryURL)
	if err != nil {
		return "", "", err
	}
	ctx := context.Background()
	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: fileHeader.Filename,
		Folder:   "review",
	})
	if err != nil {
		return "", "", err
	}
	return uploadResult.SecureURL, uploadResult.PublicID, nil
}
func CreateNewReviewHandler(c *gin.Context, db *sql.DB) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var req CreateReviewRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//
	// form, err := c.MultipartForm()
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Can't read the file image"})
	// 	return
	// }
	// files := form.File["images"]
	// if len(files) == 0 {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Review just have > 0 image"})
	// 	return
	// }

	// Bắt đầu transaction
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't start transaction"})
		return
	}
	defer tx.Rollback() // rollback nếu chưa commit

	// create review
	review := models.Review{
		ProductID: req.ProductID,
		UserID:    userID.(int64),
		OrderID:   req.OrderID,
		Rate:      req.Rate,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}
	if review.Rate > 5 || review.Rate < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rate is between 1 and 5"})
		return
	}
	canReview, err := models.CanUserReview(db, int(review.UserID), int(review.ProductID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't check buying of user"})
		return
	}
	if canReview == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you not buy this product, review cc"})
		return
	}
	review.ID, err = models.CreateReviewTx(tx, &review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// // Upload ảnh
	// var urls []string
	// for _, fileHeader := range files {
	// 	openFile, err := fileHeader.Open()
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't open image"})
	// 		return
	// 	}
	//
	// 	url, publicID, err := uploadToCloudinary(openFile, fileHeader)
	// 	openFile.Close()
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't upload file image"})
	// 		return
	// 	}
	// 	urls = append(urls, url)
	//
	// 	_, err = models.CreateReviewImagesTx(tx, review.ID, url, publicID)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't save url image"})
	// 		return
	// 	}
	// }
	//
	// Commit transaction
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Created new review successfully",
		"review":  review,
		// "images":  urls,
	})
}

// func get review by product id
func GetReviewsByProductIDHandler(c *gin.Context, db *sql.DB) {
	idStr := c.Param("id")
	productID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert product id"})
		return
	}
	// get page and limit from query(?page=1&limit=10)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	reviews, err := models.GetReviewByProductID(db, productID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	var reviewInfos []ReviewInfor
	for _, r := range reviews.Reviews {
		imgs, err := models.GetReviewImage(db, r.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get review image"})
			return
		}
		var imageRes []ReviewImageRes
		for _, img := range imgs {
			imageRes = append(imageRes, ReviewImageRes{
				ImageID: int64(img.ImageID),
				Url:     img.Url,
			})
		}

		reviewInfos = append(reviewInfos, ReviewInfor{
			UserID:    r.UserID,
			UserName:  r.UserName,
			Rate:      r.Rate,
			Content:   r.Content,
			CreatedAt: r.CreatedAt,
			Images:    imageRes,
		})

	}
	totalCount := reviews.TotalCount
	totalPage := (totalCount + limit - 1) / limit
	c.JSON(http.StatusOK, gin.H{
		"page":       page,
		"limit":      limit,
		"totalCount": totalCount,
		"totalPage":  totalPage,
		"reviews":    reviewInfos,
	})

}
