package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"one-way-ticket/db"
	"one-way-ticket/models"
	"testing"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/users", GetUsers)
	r.GET("/users/:id", GetUser)
	r.POST("/users", CreateUser)
	r.PUT("/users/:id", UpdateUser)
	r.DELETE("/users/:id", DeleteUser)
	return r
}

func TestMain(m *testing.M) {
	err := db.Connect()
	if err != nil {
		return
	}
	// Clear the users table before running tests
	db.Dbx.MustExec("TRUNCATE TABLE users RESTART IDENTITY")
	m.Run()
}

func TestGetUsers(t *testing.T) {
	router := setupRouter()

	// Insert test data
	db.Dbx.MustExec("INSERT INTO users (username, password, email) VALUES ('testuser', 'password', 'test@example.com')")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	fmt.Println(req)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var users []models.User
	err := json.Unmarshal(w.Body.Bytes(), &users)
	assert.NoError(t, err)
	assert.NotEmpty(t, users)
}

func TestGetUser(t *testing.T) {
	router := setupRouter()

	// Insert test data
	db.Dbx.MustExec("INSERT INTO users (username, password, email) VALUES ('testuser', 'password', 'test@example.com')")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var user models.User
	err := json.Unmarshal(w.Body.Bytes(), &user)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
}

func TestCreateUser(t *testing.T) {
	router := setupRouter()

	userInput := models.UserInput{
		Username: "newuser",
		Password: "newpassword",
		Email:    "newuser@example.com",
	}
	jsonValue, _ := json.Marshal(userInput)
	fmt.Println(jsonValue)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var user models.User
	err := json.Unmarshal(w.Body.Bytes(), &user)
	assert.NoError(t, err)
	assert.Equal(t, "newuser", user.Username)
}

func TestUpdateUser(t *testing.T) {
	router := setupRouter()

	// Insert test data
	db.Dbx.MustExec("INSERT INTO users (username, password, email) VALUES ('updateuser', 'password', 'update@example.com')")

	userInput := models.UserInput{
		Username: "updateduser",
		Password: "updatedpassword",
		Email:    "updated@example.com",
	}
	jsonValue, _ := json.Marshal(userInput)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var user models.User
	err := json.Unmarshal(w.Body.Bytes(), &user)
	assert.NoError(t, err)
	assert.Equal(t, "updateduser", user.Username)
}

func TestDeleteUser(t *testing.T) {
	router := setupRouter()

	// Insert test data
	db.Dbx.MustExec("INSERT INTO users (username, password, email) VALUES ('deleteuser', 'password', 'delete@example.com')")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
