package controllers

import (
	"Desktop/Projects/Hotel_Booking/config"
	"Desktop/Projects/Hotel_Booking/models"

	// "database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type GuestController struct {
	JWTKey []byte
}

func NewGuestController() *GuestController {
	return &GuestController{
		JWTKey: []byte("your_secret_key"),
	}
}

//	type guest struct {
//		Username    string `json:"username" validate:"required"`
//		Email       string `json:"email" validate:"required,email"`
//		PhoneNumber string `json:"phone_number" validate:"required"`
//		Password    string `json:"password" validate:"required, min=8"`
//	}

func (gc *GuestController) RegisterGuest(c *gin.Context) {
	var GUEST models.Guest
	if err := c.ShouldBindJSON(&GUEST); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if !isEmailValid(GUEST.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	valid, err := isValidPhoneNumber(GUEST.Phone_Number)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error validating phone number"})
		return
	}

	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number format"})
		return
	}

	// alphanumericID := generateAlphanumericID()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(GUEST.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	guestid := generateAlphanumericID()
	GUEST.Guestid = guestid
	GUEST.Password = string(hashedPassword)
	GUEST.Created_at = time.Now()
	GUEST.Updated_at = time.Now()

	if err := models.RegisterGuest(config.DB, GUEST, c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"type": "success", "message": gin.H{"message": "Guest registered successfully", "guest_id": GUEST.Guestid, "username": GUEST.Guestname, "email": GUEST.Email}})
}

func (gc *GuestController) LoginGuest(c *gin.Context) {
	var guest models.Guest
	if err := c.ShouldBindJSON(&guest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	foundGuest, err := models.GetGuestByEmail(config.DB, guest.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundGuest.Password), []byte(guest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	claims := &models.CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Subject:   foundGuest.Email,
		},
		GuestID: foundGuest.Guestid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := models.UpdateToken(config.DB, foundGuest.Guestid, tokenString); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "Message": "Login Successful"})
}

func isValidPhoneNumber(phoneNumber string) (bool, error) {
	regex := `^[789]\d{9}$`
	match, err := regexp.MatchString(regex, phoneNumber)
	if err != nil {
		return false, fmt.Errorf("error matching regex: %v", err)
	}

	return match, nil
}

func isEmailValid(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailRegex).MatchString(email)
}

func generateAlphanumericID() string {
	const (
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		length  = 8
	)
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
