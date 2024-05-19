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
	userRoutes.Use(auth.AuthenticateMiddleware())
	{
		userRoutes.GET("/", users.GetUsers)
		userRoutes.GET("/:id", users.GetUser)
		userRoutes.POST("/", users.CreateUser)
		userRoutes.PUT("/:id", users.UpdateUser)
		userRoutes.DELETE("/:id", users.DeleteUser)
	}

	return r
}
