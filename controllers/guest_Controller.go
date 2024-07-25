package controllers

import (
	"Desktop/Projects/Hotel_Booking/config"
	"Desktop/Projects/Hotel_Booking/models"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length  = 8
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func generateAlphanumericID() string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

type guest struct {
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required, min=8"`
}

var GUEST guest

func RegisterGuest(c *gin.Context) {

	err := c.ShouldBindJSON(&GUEST)

	log.Print("GUEST email:", GUEST.Email)       ///////////////
	log.Print("GUEST username:", GUEST.Username) ///!!!!!!!!!!!////////

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if !isEmailValid(GUEST.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	valid, err := isValidPhoneNumber(GUEST.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error validating phone number"})
		return
	}

	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number format"})
		return
	}

	alphanumericID := generateAlphanumericID()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(GUEST.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query := `insert into guest (guest_id, uuid, guest_name, email, phone_number, password, updated_at, created_at) values (?,?,?,?,?,?,?,?)`
	res, err := config.DB.Exec(query, alphanumericID, uuid.New().String(), GUEST.Username, GUEST.Email, GUEST.PhoneNumber, hashedPassword, time.Now(), time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println(res)

	c.JSON(http.StatusOK, gin.H{"message": "Guest registered successfully", "guest_id": alphanumericID, "username": GUEST.Username, "email": GUEST.Email})
}

func isValidPhoneNumber(phoneNumber string) (bool, error) {
	// Regular expression for basic phone number validation
	// This pattern allows for international phone numbers starting with a plus sign (+)
	// It checks for 6 to 14 digits after the plus sign
	regex := `^[789]\d{9}$`
	match, err := regexp.MatchString(regex, phoneNumber)
	if err != nil {
		return false, fmt.Errorf("error matching regex: %v", err)
	}

	return match, nil
}

// log.Print("GUEST email:", GUEST.Email)  ///////////////
// log.Print("GUEST email:", GUEST.Username) ///!!!!!!!!!!!////////
func FetchGuestID(c *gin.Context) (string, error) {

	query := `select guest_id from guest where email = ? and guest_name = ?`
	var guestid string

	log.Println("Guest email2:", GUEST.Email)
	log.Println("GIest name2:", GUEST.Username)

	err := config.DB.QueryRow(query, GUEST.Email, GUEST.Username).Scan(&guestid)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No guest found with the provided email and username"})
			return "", err
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return "", err
		}

	}
	return guestid, nil
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
	type CustomClaims struct {
		jwt.StandardClaims
		GuestID string `json:"guestid"`
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
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

	updateToken := `update guest set token = ? Where guest_id = ?`
	_, err = config.DB.Exec(updateToken, tokenString, foundGuest.Guestid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "Message": "Login Successful"})
}

func isEmailValid(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailRegex).MatchString(email)
}

// package controllers

// import (
// 	"net/http"
// 	"regexp"

// 	"Desktop/Projects/Hotel_Booking/config"
// 	"Desktop/Projects/Hotel_Booking/models"
// 	"database/sql"
// 	"fmt"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt"
// 	"github.com/google/uuid"
// 	"golang.org/x/crypto/bcrypt"
// )

// func RegisterGuest(c *gin.Context) {

// 	var guest struct {
// 		Username    string `json:"username" validate:"required"`
// 		Email       string `json:"email" validate:"required,email"`
// 		PhoneNumber string `json:"phone_number" validate:"required"`
// 		Password    string `json:"password" validate:"required, min=8"`
// 	}

// 	err := c.ShouldBindJSON(&guest)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
// 		return
// 	}

// 	// phoneNumber := "+917083296844" // Example phone number to validate
// 	valid, err := isValidPhoneNumber(guest.PhoneNumber)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error validating phone number"})
// 		return
// 	}
// 	// fmt.Println("Phone number is valid:", valid)

// 	if !valid {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number format"})
// 	} else {
// 		fmt.Println("Phone number is valid")
// 	}

// 	if !isEmailValid(guest.Email) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
// 		return
// 	}

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(guest.Password), bcrypt.DefaultCost)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return

// 	}
// 	query := `insert into guest (guest_id, guest_name, email, phone_number, password, updated_at, created_at) values (?,?,?,?,?,?,?)` ///removed "token"

// 	res, err := config.DB.Exec(query, uuid.New().String(), guest.Username, guest.Email, guest.PhoneNumber, hashedPassword, time.Now(), time.Now())
// 	fmt.Println(res)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	lastInsertID, err := res.LastInsertId()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Guest registered successfully", "guest_id": lastInsertID, "username": guest.Username, "email": guest.Email})
// }

// func LoginGuest(c *gin.Context) {
// 	var jwtKey = []byte("your_secret_key")

// 	var guest models.Guest
// 	var foundGuest models.Guest

// 	if err := c.ShouldBindJSON(&guest); err != nil {

// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	query := `SELECT guest_id, guest_name, email, phone_number, password, updated_at, created_at
//               FROM Guest WHERE email = ?`

// 	fmt.Println("Query:", query)
// 	fmt.Println("Email:", guest.Email)

// 	err := config.DB.QueryRow(query, &guest.Email).Scan(&foundGuest.Guestid, &foundGuest.Guestname, &foundGuest.Email, &foundGuest.Phone_Number, &foundGuest.Password, &foundGuest.Updated_at, &foundGuest.Created_at)
// 	if err != nil {
// 		if err == sql.ErrNoRows {

// 			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

// 		} else {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		}
// 		return
// 	}
// 	if err := bcrypt.CompareHashAndPassword([]byte(foundGuest.Password), []byte(guest.Password)); err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
// 		return
// 	}

// 	expirationTime := time.Now().Add(24 * time.Hour)
// 	claims := &jwt.StandardClaims{
// 		ExpiresAt: expirationTime.Unix(),
// 		Subject:   foundGuest.Email,
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString(jwtKey)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	updateToken := `update guest set token = ? Where guest_id = ?`
// 	_, err = config.DB.Exec(updateToken, tokenString, foundGuest.Guestid)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"token": tokenString, "Message": "Login Successfull"})
// }

// func isEmailValid(email string) bool {
// 	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
// 	return regexp.MustCompile(emailRegex).MatchString(email)
// }

// func isValidPhoneNumber(phoneNumber string) (bool, error) {
// 	// Regular expression for basic phone number validation
// 	// This pattern allows for international phone numbers starting with a plus sign (+)
// 	// It checks for 6 to 14 digits after the plus sign
// 	regex := `^\+\d{6,14}$`
// 	match, err := regexp.MatchString(regex, phoneNumber)
// 	if err != nil {
// 		return false, fmt.Errorf("error matching regex: %v", err)
// 	}

// 	return match, nil
// }
