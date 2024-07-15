// package middleware

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