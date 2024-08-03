package controllers

import (
	"Desktop/Projects/Hotel_Booking/config"
	"Desktop/Projects/Hotel_Booking/models"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetRooms(c *gin.Context) {
	filters := map[string]string{
		"availability": c.Query("availability"),
		"min_price":    c.Query("min_price"),
		"max_price":    c.Query("max_price"),
		"page":         c.Query("page"),
		"size":         c.Query("size"),
		"amenities":    c.Query("amenities"),
	}
	room, err := models.GetRooms(config.DB, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, room)
}

func GetRoom(c *gin.Context) {

	id := c.Param("room_id")

	roomID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": fmt.Sprintf("Invalid id: %s", id)})
		return
	}

	room, err := models.GetRoom(config.DB, roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	c.JSON(http.StatusOK, room)

}

func AddRoom(c *gin.Context) {
	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		log.Printf("Binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM room_type WHERE room_type_id = ?)", room.Room_Type).Scan(&exists)
	if err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room_type_id does not exist"})
		return
	}

	err = models.AddRoom(room)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room added successfully"})
}

func DelRoom(c *gin.Context) {

	roomIDstr := c.Param("room_id")
	roomID, err := strconv.Atoi(roomIDstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	err = models.DeleteRoomByID(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room deleted successfully"})

}

func EditRoom(c *gin.Context) {
	// Get the room_id from the URL
	roomID, err := strconv.Atoi(c.Param("room_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	// Bind the JSON body to the Room struct
	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the Room_ID to the one from the URL
	room.Room_ID = roomID

	// Get the DB instance
	db := config.DB

	// Update the room
	if err := models.UpdateRoom(db, room); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room updated successfully"})
}
