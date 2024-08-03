package models

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Guest struct {
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

type CustomClaims struct {
	jwt.StandardClaims
	GuestID string `json:"guestid"`
}

func RegisterGuest(db *sql.DB, GUEST Guest, c *gin.Context) error {

	query := `insert into guest (guest_id, uuid, guest_name, email, phone_number, password, updated_at, created_at) values (?,?,?,?,?,?,?,?)`
	_, err := db.Exec(query, GUEST.Guestid, uuid.New().String(), GUEST.Guestname, GUEST.Email, GUEST.Phone_Number, GUEST.Password, GUEST.Updated_at, GUEST.Created_at)
	return err
}

func GetGuestByEmail(db *sql.DB, email string) (Guest, error) {
	var guest Guest
	query := `SELECT guest_id, guest_name, email, phone_number, password, updated_at, created_at FROM guest WHERE email = ?`
	err := db.QueryRow(query, email).Scan(&guest.Guestid, &guest.Guestname, &guest.Email, &guest.Phone_Number, &guest.Password, &guest.Updated_at, &guest.Created_at)
	return guest, err
}

func FetchGuestID(db *sql.DB, GUEST Guest, c *gin.Context) (string, error) {

	query := `select guest_id from guest where email = ? and guest_name = ?`
	var guestid string

	err := db.QueryRow(query, GUEST.Email, GUEST.Guestname).Scan(&guestid)
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

func UpdateToken(db *sql.DB, guestid string, tokenstring string) error {
	updateToken := `update guest set token = ? Where guest_id = ?`
	_, err := db.Exec(updateToken, tokenstring, guestid)
	return err
}

// func IDFetch(c *gin.Context) (string, error) {
// 	guestID, exists := c.Get("guestid")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Guest ID not found"})
// 	}
// 	return guestID.(string), nil

// }
