package handlers

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"net/smtp"
	"time"

	"example.com/delivery-app/config"
	"example.com/delivery-app/middleware"
	"example.com/delivery-app/models"

	"fmt"

	"crypto/rand"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"math/big"
)

func generateOTP() (string, error) {
	// Tạo số ngẫu nhiên từ 0 đến 999999 (6 chữ số)
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %v", err)
	}
	// Định dạng số thành chuỗi 6 chữ số, thêm số 0 nếu cần
	return fmt.Sprintf("%06d", n), nil
}

//	func generateOTP() string {
//		return fmt.Sprintf("%06d", rand.Intn(1000000))
//	}
func sendEmail(to, otp string) error {
	from := config.Email.From
	password := config.Email.Password
	smtpHost := config.Email.Host
	smtpPort := config.Email.Port
	msg := []byte("Subject: OTP Verifycation\n\nYour OTP code is: " + otp)
	auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
}

// func create random string for request token
func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// func create access token and refresh token
func createTokens(user models.User) (error, string, string) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"role":   user.Role,
		"exp":    time.Now().Add(15 * time.Minute).Unix(),
	})
	accessTokenStr, err := accessToken.SignedString([]byte(middleware.JwtKey))
	if err != nil {
		return errors.New("Can't create access token"), "", ""
	}
	refreshTokenStr, err := generateRefreshToken()
	if err != nil {
		return errors.New("Can't create refresh token"), "", ""
	}
	return nil, accessTokenStr, refreshTokenStr
}

type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
type CreateShipperRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

func CreateShipper(c *gin.Context, db *sql.DB) {
	var req CreateShipperRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var exist bool
	exist, err := models.CheckEmailExists(db, req.Email)
	if exist == true {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already in use"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check email"})
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	req.Password = string(hashedPassword)
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
		Address:  req.Address,
		Role:     "shipper",
	}

	insertedID, err := models.CreateUser(db, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "id": insertedID})

}

func SignupHandler(c *gin.Context, db *sql.DB) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var exist bool
	exist, err := models.CheckEmailExists(db, req.Email)
	if exist == true {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already in use"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check email"})
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	req.Password = string(hashedPassword)
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
		Address:  req.Address,
		Role:     "customer",
	}

	insertedID, err := models.CreateUser(db, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}
	// create otp
	otp, err := generateOTP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not gen opt"})
		return
	}
	expiry := time.Now().Add(10 * time.Minute)
	// update otp to db
	if err := models.UpdateOTP(db, user.Email, otp, expiry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not set OTP"})
		return
	}
	// send email
	if err := sendEmail(user.Email, otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully, please verify your email with OTP", "id": insertedID})
}

func VerifyOTPHandler(c *gin.Context, db *sql.DB) {
	var req struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserByEmail(db, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found" + user.Email})
		return
	}

	if *user.OTPCode != req.OTP || time.Now().After(*user.OTPExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	// Cập nhật trạng thái xác thực
	if err := models.VerifyUser(db, user.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not verify user"})
		return
	}
	_ = models.ClearOTP(db, user.ID)

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}
func LoginHandler(c *gin.Context, db *sql.DB) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	exist, err := models.CheckEmailExists(db, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not check email"})
		return
	}
	if exist == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is not exist"})
		return
	}
	user, err := models.GetUserByEmail(db, req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		log.Println(err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Println(req.Password)
		log.Println(user.Password)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		log.Println(err)
		return
	}
	if !user.IsVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please verify your email before login"})
		return
	}
	err, accessTokenStr, refreshTokenStr := createTokens(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	err = models.SaveRefreshToken(db, user.ID, refreshTokenStr, time.Now().Add(7*24*time.Hour))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessTokenStr, "refresh_token": refreshTokenStr})
}
func RefreshTokenHandler(c *gin.Context, db *sql.DB) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	refreshToken, err := models.GetRefreshTokenByToken(db, req.RefreshToken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}
	if refreshToken.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token is expired"})
		return
	}
	user, err := models.GetUserByID(db, refreshToken.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server can't get user info from request"})
		return
	}
	err, accessTokenStr, newRefreshTokenStr := createTokens(*user)
	err = models.UpdateRefreshToken(db, req.RefreshToken, newRefreshTokenStr)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server can't update refresh token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": accessTokenStr, "refresh_token": newRefreshTokenStr})

}
func ForgetPasswordHandler(c *gin.Context, db *sql.DB) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := models.GetUserByEmail(db, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found"})
		return
	}
	otp, err := generateOTP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not gen OTP"})
		return
	}
	expiry := time.Now().Add(10 * time.Minute)
	if err := models.SetResetOTP(db, user.Email, otp, expiry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not set OTP"})
		return
	}
	sendEmail(user.Email, otp)
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to your email"})
}

func VerifyOTPForResetHandler(c *gin.Context, db *sql.DB) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
		OTP   string `json:"otp" binding "required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}
	user, err := models.GetUserByEmail(db, req.Email)
	if err != nil || user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is invalid"})
		return
	}
	if user.ResetOTP == nil || user.ResetOTPExpiresAt == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP not found"})
		return
	}
	if *user.ResetOTP != req.OTP || time.Now().After(*user.ResetOTPExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired otp"})
		return
	}
	claims := jwt.MapClaims{
		"email":   user.Email,
		"purpose": "reset_password",
		"exp":     time.Now().Add(5 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	// (Optional) Xoá OTP sau khi verify để không reuse
	_ = models.ClearResetOTP(db, user.ID)

	c.JSON(http.StatusOK, gin.H{
		"reset_token": tokenString,
		"expires_in":  300, // 5 phút
	})

}
func ResetPasswordHandler(c *gin.Context, db *sql.DB) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return middleware.JwtKey, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Check purpose
	if p, ok := claims["purpose"].(string); !ok || p != "reset_password" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token purpose"})
		return
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	if err := models.UpdatePasswordByEmail(db, email, string(hashed)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

func ProfileHandler(c *gin.Context, db *sql.DB) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := models.GetUserByID(db, userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	})
}
func LogoutHandler(c *gin.Context, db *sql.DB) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
	}

	err := models.DeleteRefreshToken(db, req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete refresh token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
