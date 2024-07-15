package controllers

import (
	"net/http"
	"time"

	"Desktop/Projects/Hotel_Booking/config"
	"Desktop/Projects/Hotel_Booking/models"

	// "Desktop/Projects/Hotel_Booking/models"

	"github.com/golang-jwt/jwt"

	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context) {

	var user struct {
		Username    string `json:"username" validate:"required"`
		Email       string `json:"email" validate:"required,email"`
		PhoneNumber string `json:"phone_number" validate:"required"`
		Password    string `json:"password" validate:"required, min=8"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	query := `insert INTO users (user_id, user_name, email, phone_number,updated_at, created_at, password ) values (?,?,?,?,?,?,?)`

	_, err = config.DB.Exec(query, uuid.New().String(), user.Username, user.Email, user.PhoneNumber, time.Now(), time.Now(), string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})

}

var jwtKey = []byte("your_secret_key")

func LoginUser(c *gin.Context) {

	var user models.User
	var foundUser models.User

	if err := c.ShouldBindJSON(&user); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `SELECT user_id, user_name, email, phone_number, updated_at, created_at, password
              FROM users WHERE email = ?`

	err := config.DB.QueryRow(query, user.Email).Scan(&foundUser.User_id, &foundUser.Username, &foundUser.Email, &foundUser.Phone_Number, &foundUser.Updated_at, &foundUser.Created_at, &foundUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {

			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Subject:   foundUser.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString, "Message": "Login Successfull"})

}

// func CreateUser(db *sql.DB, user User) error {
// 	apiKey, err := utils.GenerateAPIKey()
// 	if err != nil {
// 		return err
// 	}

// 	user.APIKey = apiKey
// 	// Insert user into the database
// 	_, err = db.Exec(`INSERT INTO users (name, email, api_key) VALUES (?, ?, ?)`, user.Name, user.Email, user.APIKey)
// 	return err
// }
