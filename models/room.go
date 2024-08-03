package models

import (
	"Desktop/Projects/Hotel_Booking/config"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"log"
	_ "net/http"
	"strconv"
	"strings"

	_ "github.com/gin-gonic/gin"
)

type Room struct {
	Room_ID          int      `json:"room_id,omitempty"`
	Room_Type        *string  `json:"room_type_id,omitempty"`
	Room_Description *string  `json:"room_description,omitempty"`
	Amenities        []string `json:"room_amenities,omitempty"`
	Price            *float64 `json:"price,omitempty"`
	Available        *int     `json:"available,omitempty"`
}

type Amenities []string

func (a *Amenities) Scan(value interface{}) error {
	byteArray, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(byteArray, a)
}

func (a Amenities) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func GetRooms(db *sql.DB, filters map[string]string) ([]Room, error) {
	query := `SELECT r.room_id, rt.type_name AS room_type, r.room_description AS description, r.price, 
              COALESCE(JSON_ARRAYAGG(a.amenity_name), '[]') AS amenities, r.availability
              FROM rooms r
              JOIN room_type rt ON r.room_type_id = rt.room_type_id
              LEFT JOIN room_amenity ra ON r.room_id = ra.room_id
              LEFT JOIN amenity a ON ra.amenity_id = a.amenity_id
              WHERE 1=1`

	var args []interface{}

	if filters["availability"] != "" {
		query += " AND availability = ?"
		args = append(args, filters["availability"])
	}

	if filters["min_price"] != "" { // Handle filters
		query += " AND r.price >= ?"
		args = append(args, filters["min_price"])
	}

	if filters["max_price"] != "" {
		query += " AND r.price <= ?"
		args = append(args, filters["max_price"])
	}

	if filters["amenities"] != "" {
		amenities := strings.Split(filters["amenities"], ",")
		for _, amenity := range amenities {
			query += " AND FIND_IN_SET(?, GROUP_CONCAT(a.amenity_name))"
			args = append(args, strings.TrimSpace(amenity))
		}
	}

	pageStr := filters["page"]
	if pageStr == "" {
		pageStr = "1"
	}
	sizeStr := filters["size"]
	if sizeStr == "" {
		sizeStr = "10"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(sizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	query += " GROUP BY r.room_id, rt.type_name, r.room_description, r.price, r.availability"

	query += " LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	query += ";"

	log.Printf("Executing query: %s with args: %v", query, args)

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var rooms []Room
	for rows.Next() {
		var room Room
		var amenities string

		if err := rows.Scan(&room.Room_ID, &room.Room_Type, &room.Room_Description, &room.Price, &amenities, &room.Available); err != nil {
			return nil, err
		}

		err := json.Unmarshal([]byte(amenities), &room.Amenities)
		if err != nil {
			return nil, err
		}

		rooms = append(rooms, room)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

func AddRoom(room Room) error {

	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}

	var roomTypeID int64
	err = tx.QueryRow("SELECT room_type_id FROM room_type WHERE room_type_id = ?", room.Room_Type).Scan(&roomTypeID)
	if err == sql.ErrNoRows {
		tx.Rollback()
		return errors.New("rooom_type_id does not exist")

	} else if err != nil {
		tx.Rollback()
		return err
	}

	// insert the details

	result, err := tx.Exec("INSERT INTO rooms (room_type_id, price, room_description) VALUES (?, ?, ?)",
		room.Room_Type, room.Price, room.Room_Description)

	if err != nil {
		tx.Rollback()
		return err
	}

	/// gives the last inserted autoincremented ID//
	roomID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, amenity := range room.Amenities {
		var amenityID int64
		err := tx.QueryRow("select amenity_id FROM amenity WHERE amenity_name = ?", amenity).Scan(&amenityID)
		if err == sql.ErrNoRows {
			result, err := tx.Exec("Insert INTO amenity (amenity_name) VALUES (?)", amenity)
			if err != nil {
				tx.Rollback()
				return err
			}
			amenityID, err = result.LastInsertId()
			if err != nil {
				tx.Rollback()
				return err
			}

		} else if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.Exec("INSERT INTO room_amenity (room_id, amenity_id) VALUES (?, ?)", roomID, amenityID)
		if err != nil {
			tx.Rollback()
			return err
		}

	}

	return tx.Commit()

}

func DeleteRoomByID(roomID int) error {

	// Begin a transaction
	tx, err := config.DB.Begin()
	if err != nil {
		log.Printf("Transaction error: %v", err)
		return err
	}

	// Delete room-amenity
	_, err = tx.Exec("DELETE FROM room_amenity WHERE room_id = ?", roomID)
	if err != nil {
		tx.Rollback()
		log.Printf("Delete from room_amenity error: %v", err)
		return err

	}

	log.Println("room_id deleted from room_amenity successfully")

	// Delete room
	_, err = tx.Exec("DELETE FROM rooms WHERE room_id = ?", roomID)
	if err != nil {
		tx.Rollback()
		log.Printf("Delete from room error: %v", err)
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Commit error: %v", err)
		return err
	}

	return nil

}

func UpdateRoom(db *sql.DB, room Room) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if room.Available != nil {
		query := `UPDATE rooms SET availability = ? WHERE room_id = ?`
		_, err := tx.Exec(query, *room.Available, room.Room_ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if room.Room_Type != nil {
		query := `UPDATE rooms SET room_type_id = ? WHERE room_id = ?`
		_, err := tx.Exec(query, *room.Room_Type, room.Room_ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if room.Room_Description != nil {
		query := `UPDATE rooms SET room_description = ? WHERE room_id = ?`
		_, err := tx.Exec(query, *room.Room_Description, room.Room_ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if room.Price != nil {
		query := `UPDATE rooms SET price = ? WHERE room_id = ?`
		_, err := tx.Exec(query, *room.Price, room.Room_ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(room.Amenities) > 0 {

		query := `DELETE FROM room_amenity WHERE room_id = ?`
		_, err := tx.Exec(query, room.Room_ID)
		if err != nil {
			tx.Rollback()
			return err
		}

		query = `INSERT INTO room_amenity (room_id, amenity_id) VALUES (?, (SELECT amenity_id FROM amenity WHERE amenity_name = ?))`
		for _, amenity := range room.Amenities {
			_, err = tx.Exec(query, room.Room_ID, amenity)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func GetRoom(db *sql.DB, roomID int) (*Room, error) {

	query := `select r.room_id, rt.type_name, r.room_description as description, r.price,
coalesce(JSON_ARRAYAGG(a.amenity_name), '[]') as amenities, r.availability from rooms r
join room_type rt on r.room_type_id = rt.room_type_id
join room_amenity ra on r.room_id = ra.room_id
join amenity a on ra.amenity_id = a.amenity_id where r.room_id = ? group by r.room_id, rt.type_name, r.room_description, r.price, r.availability;`

	row := config.DB.QueryRow(query, roomID)

	var room Room
	var amenities string

	err := row.Scan(&room.Room_ID, &room.Room_Type, &room.Room_Description, &room.Price, &amenities, &room.Available)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Room not found
		}
		return nil, err
	}

	// Split the amenities string into a slice of strings
	var amenitiesSlice []string
	if err := json.Unmarshal([]byte(amenities), &amenitiesSlice); err != nil {
		return nil, err
	}
	room.Amenities = amenitiesSlice

	return &room, nil
}
