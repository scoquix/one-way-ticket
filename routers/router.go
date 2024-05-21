package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"one-way-ticket/auth"
	"one-way-ticket/users"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/health", func(c *gin.Context) {
		// Get name from query parameter (optional)
		name := c.Query("name")
		if name != "" {
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Hello, %s!", name)})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
		}
	})
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
