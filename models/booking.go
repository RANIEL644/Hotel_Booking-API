package models

import (
	"Desktop/Projects/Hotel_Booking/config"
	"time"

	// "github.com/gin-gonic/gin"
	"fmt"
)

type Booking struct {
	Booking_ID     uint      `json:"booking_id"`
	Room_ID        int       `json:"room_id"`
	Guest_ID       int       `json:"guest_id"`
	From_Date      string    `json:"from_date"`
	To_Date        string    `json:"to_date"`
	Adults         int       `json:"adults"`
	Children       int       `json:"children"`
	Check_in_Time  string    `json:"checkin_time"`
	Check_out_Time string    `json:"checkout_time"`
	TotalPrice     float64   `json:"total_price"`
	Status         string    `json:"status" gorm:"default:'Confirmed'"`
	Booking_Date   time.Time `json:"booking_date"`
}

func BookRoom(booking *Booking) (uint, error) {
	fmt.Println("hellow")

	query := `INSERT INTO bookings (booking_id, room_id, guest_id, num_of_adults, num_of_children, checkin_date, checkout_date, checkin_time, checkout_time, price, status, booking_date)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := config.DB.Exec(query, booking.Booking_ID, booking.Room_ID, booking.Guest_ID, booking.Adults, booking.Children, booking.From_Date, booking.To_Date, booking.Check_in_Time, booking.Check_out_Time, booking.TotalPrice, booking.Status, time.Now())
	if err != nil {
		return 0, err
	}

	bookingID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint(bookingID), nil
}

func UpdateRoomAvailability(roomID int, available bool) error {
	query := "UPDATE rooms SET available = ? WHERE room_id = ?"
	_, err := config.DB.Exec(query, available, roomID)
	return err
}

// func UpdateRoomAvailability(roomID int, isAvailable bool) error {
//   query := `update rooms set available = ? where room_id = ?`

//   _, err := config.DB.Exec(query, )

//     return nil
// }
