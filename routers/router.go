package routers

import (
	"github.com/gin-gonic/gin"
	"one-way-ticket/auth"
	"one-way-ticket/dynamo"
	"one-way-ticket/service/bookings"
	"one-way-ticket/service/movies"
	"one-way-ticket/service/showtimes"
	"one-way-ticket/service/users"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	handler := auth.NewHandler(dynamo.NewDynamoClient())
	r.POST("/login", handler.Login)

	userRoutes := r.Group("/users")
	userRoutes.Use(handler.AuthenticateMiddleware())
	{
		userRoutes.GET("/", users.GetUsers)
		userRoutes.GET("/:id", users.GetUser)
		userRoutes.POST("/", users.CreateUser)
		userRoutes.PUT("/:id", users.UpdateUser)
		userRoutes.DELETE("/:id", users.DeleteUser)
	}

	moviesRoutes := r.Group("/movies")
	moviesRoutes.Use(handler.AuthenticateMiddleware())
	{
		moviesRoutes.GET("/", movies.GetMovies)
		moviesRoutes.GET("/:id", movies.GetMovie)
		moviesRoutes.POST("/", movies.CreateMovie)
		moviesRoutes.PUT("/:id", movies.UpdateMovie)
		moviesRoutes.DELETE("/:id", movies.DeleteMovie)
	}

	showTimesRoutes := r.Group("/showtimes")
	showTimesRoutes.Use(handler.AuthenticateMiddleware())
	{
		showTimesRoutes.GET("/", showtimes.GetShowtimes)
		showTimesRoutes.GET("/:id", showtimes.GetShowtime)
		showTimesRoutes.POST("/", showtimes.CreateShowtime)
		showTimesRoutes.PUT("/:id", showtimes.UpdateShowtime)
		showTimesRoutes.DELETE("/:id", showtimes.DeleteShowtime)
	}

	bookingsRoutes := r.Group("/bookings")
	bookingsRoutes.Use(handler.AuthenticateMiddleware())
	{
		bookingsRoutes.GET("/", bookings.GetBookings)
		bookingsRoutes.GET("/:id", bookings.GetBooking)
		bookingsRoutes.POST("/", bookings.CreateBooking)
		bookingsRoutes.PUT("/:id", bookings.UpdateBooking)
		bookingsRoutes.DELETE("/:id", bookings.DeleteBooking)
	}

	return r
}
