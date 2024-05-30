package users

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"one-way-ticket/db"
	"one-way-ticket/models"
	"os"
	"strconv"
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
	_, err = db.Dbx.Exec("TRUNCATE TABLE bookings, showtimes, movies, users RESTART IDENTITY CASCADE")
	if err != nil {
		panic(err)
	}
	code := m.Run()
	err = db.Dbx.Close()
	if err != nil {
		panic(err)
	}
	os.Exit(code)
}

func TestGetUsers(t *testing.T) {
	router := setupRouter()

	db.Dbx.MustExec("INSERT INTO users (username, password, email) VALUES ('testuser', 'password', 'test@example.com')")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var users []models.User
	err := json.Unmarshal(w.Body.Bytes(), &users)
	assert.NoError(t, err)
}

func TestGetUser(t *testing.T) {
	router := setupRouter()

	db.Dbx.MustExec("INSERT INTO users (username, password, email) VALUES ('testuser2', 'password', 'test2@example.com')")
	var userID int
	err := db.Dbx.Get(&userID, "SELECT user_id FROM users WHERE username='testuser2' AND email='test2@example.com'")
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+strconv.Itoa(userID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var user models.User
	err = json.Unmarshal(w.Body.Bytes(), &user)
	assert.NoError(t, err)
}

func TestCreateUser(t *testing.T) {
	router := setupRouter()

	userInput := models.UserInput{
		Username: "newuser",
		Password: "newpassword",
		Email:    "newuser@example.com",
	}
	jsonValue, _ := json.Marshal(userInput)
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

	db.Dbx.MustExec("INSERT INTO users (username, password, email) VALUES ('updateuser', 'password', 'update@example.com')")
	var userID int
	err := db.Dbx.Get(&userID, "SELECT user_id FROM users WHERE username='updateuser' AND email='update@example.com'")
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	userInput := models.UserInput{
		Username: "updateduser",
		Password: "updatedpassword",
		Email:    "updated@example.com",
	}
	jsonValue, _ := json.Marshal(userInput)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/"+strconv.Itoa(userID), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var user models.User
	err = json.Unmarshal(w.Body.Bytes(), &user)
	assert.NoError(t, err)
	assert.Equal(t, "updateduser", user.Username)
}

func TestDeleteUser(t *testing.T) {
	router := setupRouter()

	db.Dbx.MustExec("INSERT INTO users (username, password, email) VALUES ('deleteuser', 'password', 'delete@example.com')")
	var userID int
	err := db.Dbx.Get(&userID, "SELECT user_id FROM users WHERE username='deleteuser' AND email='delete@example.com'")
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/"+strconv.Itoa(userID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
