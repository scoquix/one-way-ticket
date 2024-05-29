package routers

import (
	"github.com/gin-gonic/gin"
	"one-way-ticket/auth"
	"one-way-ticket/movies"
	"one-way-ticket/users"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/login", auth.Login)

	userRoutes := r.Group("/users")
	userRoutes.Use(auth.AuthenticateMiddleware())
	{
		userRoutes.GET("/", users.GetUsers)
		userRoutes.GET("/:id", users.GetUser)
		userRoutes.POST("/", users.CreateUser)
		userRoutes.PUT("/:id", users.UpdateUser)
		userRoutes.DELETE("/:id", users.DeleteUser)
	}

	moviesRoutes := r.Group("/movies")
	moviesRoutes.Use(auth.AuthenticateMiddleware())
	{
		moviesRoutes.GET("/", movies.GetMovies)
		moviesRoutes.GET("/:id", movies.GetMovie)
		moviesRoutes.POST("/", movies.CreateMovie)
		moviesRoutes.PUT("/:id", movies.UpdateMovie)
		moviesRoutes.DELETE("/:id", movies.DeleteMovie)
	}

	return r
}
