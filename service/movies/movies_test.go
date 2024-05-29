package movies

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"one-way-ticket/db"
	"one-way-ticket/models"
	"strconv"
	"testing"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/movies", GetMovies)
	r.GET("/movies/:id", GetMovie)
	r.POST("/movies", CreateMovie)
	r.PUT("/movies/:id", UpdateMovie)
	r.DELETE("/movies/:id", DeleteMovie)
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
	m.Run()
}

func TestGetMovies(t *testing.T) {
	router := setupRouter()

	db.Dbx.MustExec("INSERT INTO movies (title, duration, genre) VALUES ('test', '123', 'comedy')")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/movies", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var users []models.User
	err := json.Unmarshal(w.Body.Bytes(), &users)
	assert.NoError(t, err)
}

func TestGetMovie(t *testing.T) {
	router := setupRouter()

	db.Dbx.MustExec("INSERT INTO movies (title, duration, genre) VALUES ('Inception', 148, 'Sci-Fi')")
	var movieID int
	err := db.Dbx.Get(&movieID, "SELECT movie_id FROM movies WHERE title='Inception' AND duration=148 AND genre='Sci-Fi'")
	if err != nil {
		t.Fatalf("Failed to get booking ID: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/movies/"+strconv.Itoa(movieID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var movie models.Movie
	err = json.Unmarshal(w.Body.Bytes(), &movie)
	assert.NoError(t, err)
}

func TestCreateMovie(t *testing.T) {
	router := setupRouter()

	movieInput := models.MovieInput{
		Title:    "Interstellar",
		Duration: 169,
		Genre:    "Sci-Fi",
	}
	jsonValue, _ := json.Marshal(movieInput)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/movies", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var movie models.Movie
	err := json.Unmarshal(w.Body.Bytes(), &movie)
	assert.NoError(t, err)
	assert.Equal(t, "Interstellar", movie.Title)
}

func TestUpdateMovie(t *testing.T) {
	router := setupRouter()

	db.Dbx.MustExec("INSERT INTO movies (title, duration, genre) VALUES ('Inception', 148, 'Sci-Fi')")
	var movieID int
	err := db.Dbx.Get(&movieID, "SELECT movie_id FROM movies WHERE title='Inception' AND duration=148 AND genre='Sci-Fi'")
	if err != nil {
		t.Fatalf("Failed to get booking ID: %v", err)
	}

	movieInput := models.MovieInput{
		Title:    "Inception Updated",
		Duration: 150,
		Genre:    "Sci-Fi",
	}
	jsonValue, _ := json.Marshal(movieInput)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/movies/"+strconv.Itoa(movieID), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var movie models.Movie
	err = json.Unmarshal(w.Body.Bytes(), &movie)
	assert.NoError(t, err)
	assert.Equal(t, "Inception Updated", movie.Title)
}

func TestDeleteMovie(t *testing.T) {
	router := setupRouter()

	db.Dbx.MustExec("INSERT INTO movies (title, duration, genre) VALUES ('Inception', 148, 'Sci-Fi')")

	var movieID int
	err := db.Dbx.Get(&movieID, "SELECT movie_id FROM movies WHERE title='Inception' AND duration=148 AND genre='Sci-Fi'")
	if err != nil {
		t.Fatalf("Failed to get booking ID: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/movies/"+strconv.Itoa(movieID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
