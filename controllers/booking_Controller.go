package controllers

import (
	"Desktop/Projects/Hotel_Booking/config"
	"Desktop/Projects/Hotel_Booking/models"

	// "errors"
	"net/http"
	"time"

	"fmt"

	"strconv"

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

	booking.Room_ID, _ = strconv.Atoi(roomID)

	// Parse dates

	// fromDateString := c.Query("from_date")
	// toDateString := c.Query("to_date")

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

	// Example calculation
	duration := toDate.Sub(fromDate)
	booking.TotalPrice = float64(duration.Hours() / 24 * 1000)

	bookingID, err := models.BookRoom(&booking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update room availability
	err = models.UpdateRoomAvailability(booking.Room_ID, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room availability"})
		return
	}

	response := gin.H{
		"booking_id":     bookingID,
		"room_id":        roomID,
		"guest_id":       booking.Guest_ID,
		"adults":         booking.Adults,
		"children":       booking.Children,
		"check_in_date":  fromDate.Format("2006-01-02"),
		"check_out_date": toDate.Format("2006-01-02"),
		"check_in_time":  booking.Check_in_Time,
		"check_out_time": booking.Check_out_Time,
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
