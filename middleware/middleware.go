package middleware

import (
	"Desktop/Projects/Hotel_Booking/config"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func APIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
			c.Abort()
			return
		}

		var userID string
		err := config.DB.QueryRow("select user_id from users where api_key = ?").Scan(&userID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			c.Abort()
			return
		}

		c.Set("UserID", userID)
		c.Next()
	}

}

// import (
// 	"net/http"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// )

// func APIKeyAuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		apiKey := c.GetHeader("X-API-KEY")
// 		if apiKey == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
// 			c.Abort()
// 			return
// 		}

// 		c.Next()
// 	}
// }

// func validateAPIKey(apiKey string) bool {
// 	// This function should check if the API key exists in the database
// 	// For demonstration purposes, let's assume it always returns true
// 	// Implement your logic to validate the API key against your database

// 	// Example:
// 	// var exists bool
// 	// err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE api_key = ?)`, apiKey).Scan(&exists)
// 	// if err != nil {
// 	//     return false
// 	// }
// 	// return exists

// 	return true // Replace with actual validation
// }
