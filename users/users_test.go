package users

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestT(t *testing.T) {
	db, err := sqlx.Connect("postgres", "user=test dbname=test sslmode=disable password=test host=localhost")
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Connected")
	}
}

func TestGetUsers(t *testing.T) {
	router := gin.Default()
	router.GET("/users/", GetUsers)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateUser(t *testing.T) {
	router := gin.Default()
	router.POST("/users/", CreateUser)

	w := httptest.NewRecorder()
	user := `{"username": "testuser", "password": "testpass"}`
	req, _ := http.NewRequest("POST", "/users/", strings.NewReader(user))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	var response map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "testuser", response["username"])
}
