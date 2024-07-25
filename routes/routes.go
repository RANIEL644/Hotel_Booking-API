package routes

import (
	"Desktop/Projects/Hotel_Booking/controllers"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {
	router.POST("guest/register", controllers.RegisterGuest) //done
	router.POST("guest/login", controllers.LoginGuest)       //done

	router.POST("users/register", controllers.RegisterUser) //done
	router.POST("users/login", controllers.LoginUser)       //done *need to implement x-api key
	router.GET("/rooms", controllers.GetRooms)              /// view all rooms //done
	router.GET("rooms/:room_id", controllers.GetRoom)       // select a particular room
	// router.POST("/rooms", controllers.AddRoom)               //done
	// router.DELETE("rooms/:room_id", controllers.DelRoom)      //done
	// router.PATCH("rooms/:room_id", controllers.EditRoom)      //edit a room  *******
	// router.POST("/rooms/:room_id/book", controllers.BookRoom) //done
	router.POST("booking")
	router.GET("booking")
	router.DELETE("booking")

}

//protected routes

func ProtectedRoutes(router *gin.Engine) {
	protected := router.Group("/")
	protected.Use(controllers.AuthMiddleware())
	{
		protected.DELETE("/rooms/:room_id", controllers.DelRoom)
		protected.PATCH("rooms/:room_id", controllers.EditRoom)
		protected.POST("/rooms", controllers.AddRoom)
		protected.POST("/rooms/:room_id/book", controllers.BookRoom)

		// Add more protected routes here
	}
}
