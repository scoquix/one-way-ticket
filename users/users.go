package users

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"one-way-ticket/models"
	"strconv"
)

func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, models.Users)
}

func GetUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for _, user := range models.Users {
		if user.ID == uint(id) {
			c.JSON(http.StatusOK, user)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"status": "user not found"})
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err == nil {
		user.ID = uint(len(models.Users) + 1)
		models.Users = append(models.Users, user)
		c.JSON(http.StatusCreated, user)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
	}
}

func UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var updateUser models.User
	if err := c.BindJSON(&updateUser); err == nil {
		for i, user := range models.Users {
			if user.ID == uint(id) {
				models.Users[i] = updateUser
				c.JSON(http.StatusOK, updateUser)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"status": "user not found"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
	}
}

func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, user := range models.Users {
		if user.ID == uint(id) {
			models.Users = append(models.Users[:i], models.Users[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"status": "user deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"status": "user not found"})
}
