package routes

import (
	"Desktop/Projects/Hotel_Booking/controllers"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {

	guestController := &controllers.GuestController{}
	router.POST("guest/register", guestController.RegisterGuest) //done
	router.POST("guest/login", guestController.LoginGuest)       //done
	router.POST("users/register", controllers.RegisterUser)      //done
	router.POST("users/login", controllers.LoginUser)            //done *need to implement x-api key
	router.GET("/rooms", controllers.GetRooms)
	router.GET("rooms/:room_id", controllers.GetRoom)
	router.POST("bookings/booking_id", controllers.GetBookingByID)

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
		protected.GET("/guest/login/bookings", controllers.GetGuestBookings)

		// Add more protected routes here
	}
}
