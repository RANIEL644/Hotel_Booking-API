package controllers

import (
	"Desktop/Projects/Hotel_Booking/config"
	"Desktop/Projects/Hotel_Booking/models"
	Utils "Desktop/Projects/Hotel_Booking/utils"
	"log"

	// "errors"
	// "database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func BookRoom(c *gin.Context) {

	roomID := c.Param("room_id")
	if roomID == "" {

		c.JSON(http.StatusBadRequest, gin.H{"error": "room_id is required"})
		return
	}

	var booking models.Booking

	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}
	////////////////////////////
	guestID, err := Utils.ExtractGuestIDFromToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token", "M": err.Error()})
		return

	}
	booking.Guest_ID = string(guestID)
	fmt.Println(booking.Guest_ID)
	booking.Room_ID, _ = strconv.Atoi(roomID)

	//////////////////////////////////////////////
	if booking.From_Date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from_date is empty"})
		return
	}

	fromDate, err := time.Parse("2006-01-02", booking.From_Date)
	if err != nil {
		fmt.Println("Error parsing date:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from_date format, use yyyy-mm-dd"})
		return
	}

	toDate, err := time.Parse("2006-01-02", booking.To_Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid to_date format, use yyyy-mm-dd"})
		return
	}

	checkinTime, err := time.Parse("15:04:05", booking.Check_in_Time)
	if err != nil {
		fmt.Println("error parsing time", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "wrong time")
		return
	}

	checkoutTime, err := time.Parse("15:04:05", booking.Check_out_Time)
	if err != nil {
		fmt.Println("error parsing time", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "wrong time")

		return
	}

	roomPrice, err := models.FetchRoomPrice(booking.Room_ID)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch room price"})
		return
	}
	// price calculation
	duration := toDate.Sub(fromDate)
	booking.TotalPrice = float64(duration.Hours() / 24 * roomPrice)
	log.Println(booking.TotalPrice)

	// Update room availability
	err = models.UpdateRoomAvailability(booking.Room_ID, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room availability"})
		return
	}

	guestid, err := FetchGuestID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch GuestID"})
	}

	bookingID, err := models.BookRoom(&booking, guestid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"booking_id":     bookingID,
		"room_id":        roomID,
		"guest_id":       guestid,
		"adults":         booking.Adults,
		"children":       booking.Children,
		"check-in_date":  fromDate.Format("2006-01-02"),
		"check-out_date": toDate.Format("2006-01-02"),
		"check-in_time":  checkinTime,
		"check-out_time": checkoutTime,
		"total_price":    booking.TotalPrice,
		"status":         "Confirmed",
	}

	c.JSON(http.StatusOK, response)
}

func GetBookingByID(bookingID int) (*models.Booking, error) {
	var booking models.Booking
	err := config.DB.QueryRow("SELECT * FROM bookings WHERE booking_id = ?", bookingID).Scan(
		&booking.Booking_ID, &booking.Room_ID, &booking.Guest_ID, &booking.From_Date, &booking.To_Date,
		&booking.Adults, &booking.Children, &booking.Check_in_Time, &booking.Check_out_Time,
		&booking.TotalPrice, &booking.Status)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

// var userdetails struct {
// 	Username    string `json:"username" validate:"required"`
// 	Email       string `json:"email" validate:"required,email"`
// 	PhoneNumber string `json:"phone_number" validate:"required"`
// 	Password    string `json:"password" validate:"required, min=8"`
// }

// query := `select guest_id from guest where email = ? and guest_name = ?`
// var guestid string

// userdetails.Username = "joe"
// userdetails.Email = "joeb@gmail.com"

// err = config.DB.QueryRow(query, userdetails.Email, userdetails.Username).Scan(&guestid)
// if err != nil {
// 	if err == sql.ErrNoRows {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "No guest found with the provided email and username"})
// 	} else {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 	}
// 	return
// }
