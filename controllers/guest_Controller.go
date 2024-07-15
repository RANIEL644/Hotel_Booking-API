package controllers

import (
	"net/http"

	"Desktop/Projects/Hotel_Booking/config"
	"Desktop/Projects/Hotel_Booking/models"
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterGuest(c *gin.Context) {

	var guest struct {
		Username    string `json:"username" validate:"required"`
		Email       string `json:"email" validate:"required,email"`
		PhoneNumber string `json:"phone_number" validate:"required"`
		Password    string `json:"password" validate:"required, min=8"`
	}

	err := c.ShouldBindJSON(&guest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(guest.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	query := `insert into guest (guest_id, guest_name, email, phone_number, password, updated_at, created_at) values (?,?,?,?,?,?,?)` ///removed "token"

	res, err := config.DB.Exec(query, uuid.New().String(), guest.Username, guest.Email, guest.PhoneNumber, hashedPassword, time.Now(), time.Now())
	fmt.Println(res)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Guest registered successfully"})
}

func LoginGuest(c *gin.Context) {
	var jwtKey = []byte("your_secret_key")

	var guest models.Guest
	var foundGuest models.Guest

	if err := c.ShouldBindJSON(&guest); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `SELECT guest_id, guest_name, email, phone_number, password, updated_at, created_at
              FROM Guest WHERE email = ?`

	fmt.Println("Query:", query)
	fmt.Println("Email:", guest.Email)

	err := config.DB.QueryRow(query, &guest.Email).Scan(&foundGuest.Guestid, &foundGuest.Guestname, &foundGuest.Email, &foundGuest.Phone_Number, &foundGuest.Password, &foundGuest.Updated_at, &foundGuest.Created_at)
	if err != nil {
		if err == sql.ErrNoRows {

			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(foundGuest.Password), []byte(guest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Subject:   foundGuest.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updateToken := `update guest set token = ? Where guest_id = ?`
	_, err = config.DB.Exec(updateToken, tokenString, foundGuest.Guestid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "Message": "Login Successfull"})
}
