package controllers

import (
	"Desktop/Projects/Hotel_Booking/config"
	"Desktop/Projects/Hotel_Booking/models"
	Utils "Desktop/Projects/Hotel_Booking/utils"
	"fmt"
	"net/http"
	"time"

	// "Desktop/Projects/Hotel_Booking/models"

	"github.com/golang-jwt/jwt"

	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	JWTKey []byte
}

func NewUserController() *UserController {
	return &UserController{
		JWTKey: []byte("your_secret_key"),
	}
}

func (uc *UserController) RegisterUser(c *gin.Context) {
	var USER models.User
	// var user struct {
	// 	Username    string `json:"username" validate:"required"`
	// 	Email       string `json:"email" validate:"required,email"`
	// 	PhoneNumber string `json:"phone_number" validate:"required"`
	// 	Password    string `json:"password" validate:"required, min=8"`
	// }

	if err := c.ShouldBindJSON(&USER); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isEmailValid(USER.Email) {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
	}

	valid, err := isValidPhoneNumber((USER.Phone_Number))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error validating phone number"})
		return
	}

	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number format"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(USER.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	userid := generateAlphanumericID()
	USER.User_id = userid
	USER.Password = string(hashedPassword)
	USER.Created_at = time.Now()
	USER.Updated_at = time.Now()
	USER.API_Key, err = Utils.GenerateAPIKey()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to generate API-Key"})
	}

	query := `insert INTO users (user_id, user_name, email, phone_number,updated_at, created_at, password, api_key ) values (?,?,?,?,?,?,?,?)`

	_, err = config.DB.Exec(query, uuid.New().String(), USER.Username, USER.Email, USER.Phone_Number, time.Now(), time.Now(), string(hashedPassword), USER.API_Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"type": "success", "message": gin.H{"message": "User registered successfully", "guest_id": USER.User_id, "username": USER.Username, "email": USER.Email}})
}

// //////////////////////////////////////////////////////////
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

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// claims := &jwt.StandardClaims{}

		guestID, token, err := Utils.ExtractGuestIDFromTokenn(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 	return jwtKey, nil
		// })

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// c.Set("userEmail", claims.Subject)
		c.Set("guestid", guestID)
		c.Next()

		fmt.Println("guestid", guestID)
	}
}
