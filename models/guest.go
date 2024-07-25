package models

import (
	// "Desktop/Projects/Hotel_Booking/config"
	// "database/sql"
	// "net/http"
	"time"
	// "github.com/gin-gonic/gin"
)

type Guest struct {
	// ID            int       `json:"id"`
	Guestid       string    `json:"guestid,omitempty"`
	Guestname     string    `json:"username,omitempty"`
	Email         string    `json:"email,omitempty"`
	Password      string    `json:"password"`
	Phone_Number  string    `json:"phone_number,omitempty"`
	Token         string    `json:"token,omitempty"`
	User_Type     string    `json:"user_type,omitempty"`
	Refresh_token string    `json:"refresh_token,omitempty"`
	Created_at    time.Time `json:"created_at,omitempty"`
	Updated_at    time.Time `json:"updated_at,omitempty"`
	User_id       string    `json:"user_id,omitempty"`
}

// func FetchGuestDetails(c *gin.Context) (string, error) {

// 	query := `select guest_id from guest where email = ? and guest_name = ?`
// 	var guestid string

// 	var GUEST Guest

// 	err := config.DB.QueryRow(query, GUEST.Email, GUEST.Guestname).Scan(&guestid)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "No guest found with the provided email and username"})
// 			return "", err
// 		} else {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return "", err
// 		}

// 	}
// 	return guestid, nil
// }
