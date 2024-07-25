package main

import (
	"Desktop/Projects/Hotel_Booking/config"
	"Desktop/Projects/Hotel_Booking/models"

	"Desktop/Projects/Hotel_Booking/routes"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load() //load environment variables
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	router := gin.Default() ///initialize Gin router

	routes.InitializeRoutes(router) //initialize routes

	config.InitDB()
	defer config.DB.Close()

	filters := map[string]string{
		"min_price": "500",
		"max_price": "1000",
		"page":      "1",
		"size":      "10",
	}

	rooms, err := models.GetRooms(config.DB, filters)
	if err != nil {
		log.Fatal("Error getting rooms:", err)
	}

	for _, room := range rooms {
		fmt.Printf("%+v\n", room)
	}

	routes.ProtectedRoutes(router)

	router.Run(":8080")
}
