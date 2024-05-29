package users

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"one-way-ticket/db"
	"one-way-ticket/models"
	"strconv"
)

var log = logrus.New()

const (
	InvalidUserId = "Invalid user ID"
)

func GetUsers(c *gin.Context) {
	var users []models.User
	err := db.Dbx.Select(&users, "SELECT * FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": InvalidUserId})
		return
	}

	var user models.User
	err = db.Dbx.Get(&user, "SELECT * FROM users WHERE user_id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	var userInput models.UserInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		log.Error("Error binding JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO Users (username, password, email) VALUES (:username, :password, :email) RETURNING user_id`
	rows, err := db.Dbx.NamedQuery(query, &userInput)
	if err != nil {
		log.Error("Error inserting user: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Username: userInput.Username,
		Password: userInput.Password,
		Email:    userInput.Email,
	}

	if rows.Next() {
		err = rows.Scan(&user.ID)
		if err != nil {
			log.Error("Error retrieving new user ID: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	log.Info("User created successfully with ID:", user.ID)
	c.JSON(http.StatusCreated, user)
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": InvalidUserId})
		return
	}

	var userInput models.UserInput
	if err := c.BindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		ID:       uint(id),
		Username: userInput.Username,
		Password: userInput.Password,
		Email:    userInput.Email,
	}

	_, err = db.Dbx.NamedExec("UPDATE users SET username=:username, password=:password, email=:email WHERE user_id=:user_id", &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": InvalidUserId})
		return
	}

	_, err = db.Dbx.Exec("DELETE FROM users WHERE user_id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}
