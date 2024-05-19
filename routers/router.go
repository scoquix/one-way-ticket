package routers

import (
	"github.com/gin-gonic/gin"
	"one-way-ticket/auth"
	"one-way-ticket/users"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/login", auth.Login)

	userRoutes := r.Group("/users")
	{
		userRoutes.GET("/", users.GetUsers)
		userRoutes.GET("/:id", users.GetUser)
		userRoutes.POST("/", users.CreateUser)
		userRoutes.PUT("/:id", users.UpdateUser)
		userRoutes.DELETE("/:id", auth.AuthenticateMiddleware(), users.DeleteUser)
	}

	return r
}
