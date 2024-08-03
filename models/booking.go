package models

import (
	"Desktop/Projects/Hotel_Booking/config"
	// "Desktop/Projects/Hotel_Booking/controllers"

	// "Desktop/Projects/Hotel_Booking/controllers"
	"time"

	// "github.com/gin-gonic/gin"
	"fmt"
)

type Booking struct {
	Booking_ID     uint      `json:"booking_id"`
	Room_ID        int       `json:"room_id"`
	Guest_ID       string    `json:"guest_id"`
	From_Date      string    `json:"from_date"`
	To_Date        string    `json:"to_date"`
	Adults         int       `json:"adults"`
	Children       int       `json:"children"`
	Check_in_Time  string    `json:"check-in time"`
	Check_out_Time string    `json:"check-out time"`
	TotalPrice     float64   `json:"total_price"`
	Status         string    `json:"status" gorm:"default:'Confirmed'"`
	Booking_Date   time.Time `json:"booking_date"`
}

func BookRoom(booking *Booking, guestid string) (uint, error) {

	query := `INSERT INTO bookings (booking_id, room_id, guest_id, num_of_adults, num_of_children, checkin_date, checkout_date, checkin_time, checkout_time, price, status, booking_date)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// booking.Status = "Booked"
	result, err := config.DB.Exec(query, booking.Booking_ID, booking.Room_ID, guestid, booking.Adults, booking.Children, booking.From_Date, booking.To_Date, booking.Check_in_Time, booking.Check_out_Time, booking.TotalPrice, string("Booked"), time.Now())
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
	query := "UPDATE rooms SET availability = ? WHERE room_id = ?"
	_, err := config.DB.Exec(query, available, roomID)
	if err != nil {
		fmt.Println(err)
	}

	return err

}

func FetchRoomPrice(roomID int) (float64, error) {
	var price float64

	fmt.Println("Fetching price from roomID:", roomID)

	err := config.DB.QueryRow("SELECT price FROM rooms WHERE room_id = ?", roomID).Scan(&price)

	if err != nil {
		fmt.Println("Error fetching room price:", err)
		return 0, err
	}

	fmt.Println("Fetched price:", price)

	return price, nil
}

func GetBooking(booking_id int) (Booking, error) {
	var booking Booking

	query := "select * from bookings where booking_id = ?"
	row := config.DB.QueryRow(query, booking_id)
	if err := row.Scan(&booking.Booking_ID, &booking.Room_ID, &booking.Guest_ID, &booking.Adults, &booking.Children, &booking.To_Date, &booking.From_Date, &booking.Check_in_Time, &booking.Check_out_Time, &booking.TotalPrice, &booking.Status, &booking.Booking_Date); err != nil {
		fmt.Println("failed to scan to row")

		return booking, err
	}
	return booking, nil
}

// func UpdateRoomAvailability(roomID int, isAvailable bool) error {
//   query := `update rooms set available = ? where room_id = ?`

//   _, err := config.DB.Exec(query, )

//     return nil
// }
