package showtimes

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
	r.GET("/showtimes", GetShowtimes)
	r.GET("/showtimes/:id", GetShowtime)
	r.POST("/showtimes", CreateShowtime)
	r.PUT("/showtimes/:id", UpdateShowtime)
	r.DELETE("/showtimes/:id", DeleteShowtime)
	return r
}

func TestMain(m *testing.M) {
	err := db.Connect()
	if err != nil {
		panic(err)
	}
	_, err = db.Dbx.Exec("TRUNCATE TABLE bookings, showtimes, movies, users RESTART IDENTITY CASCADE")
	if err != nil {
		panic(err)
	}
	_, err = db.Dbx.Exec("INSERT INTO movies (title, duration, genre) VALUES ('Sample Movie', 120, 'Action')")
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

func TestGetShowtimes(t *testing.T) {
	router := setupRouter()

	db.Dbx.MustExec("INSERT INTO showtimes (movie_id, showtime, hall) VALUES (1, '2024-05-30 12:00:00', 'Hall 1')")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/showtimes", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var showtimes []models.Showtime
	err := json.Unmarshal(w.Body.Bytes(), &showtimes)
	assert.NoError(t, err)
	assert.NotEmpty(t, showtimes)
}

func TestGetShowtime(t *testing.T) {
	router := setupRouter()

	db.Dbx.MustExec("INSERT INTO showtimes (movie_id, showtime, hall) VALUES (1, '2024-05-30 12:00:00', 'Hall 1')")
	var showtimeID int
	err := db.Dbx.Get(&showtimeID, "SELECT showtime_id FROM showtimes WHERE movie_id=1 AND showtime='2024-05-30 12:00:00' AND hall='Hall 1'")
	if err != nil {
		t.Fatalf("Failed to get showtime ID: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/showtimes/"+strconv.Itoa(showtimeID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var showtime models.Showtime
	err = json.Unmarshal(w.Body.Bytes(), &showtime)
	assert.NoError(t, err)
	assert.Equal(t, 1, showtime.MovieID)
	assert.Equal(t, "Hall 1", showtime.Hall)
}

func TestCreateShowtime(t *testing.T) {
	router := setupRouter()

	showtimeInput := models.ShowtimeInput{
		MovieID:  1,
		Showtime: "2024-05-30 16:00",
		Hall:     "Hall 1",
	}
	jsonValue, _ := json.Marshal(showtimeInput)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/showtimes", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var showtime models.Showtime
	err := json.Unmarshal(w.Body.Bytes(), &showtime)
	assert.NoError(t, err)
	assert.Equal(t, 1, showtime.MovieID)
	assert.Equal(t, "2024-05-30 16:00", showtime.Showtime)
	assert.Equal(t, "Hall 1", showtime.Hall)
}

func TestCreateShowtimeOverlap(t *testing.T) {
	router := setupRouter()

	db.Dbx.MustExec("INSERT INTO showtimes (movie_id, showtime, hall) VALUES (1, '2024-05-30 12:00:00', 'Hall 1')")

	showtimeInput := models.ShowtimeInput{
		MovieID:  1,
		Showtime: "2024-05-30 13:00",
		Hall:     "Hall 1",
	}
	jsonValue, _ := json.Marshal(showtimeInput)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/showtimes", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), OverlappingShowtimeError)
}

func TestUpdateShowtime(t *testing.T) {
	router := setupRouter()

	db.Dbx.MustExec("INSERT INTO showtimes (movie_id, showtime, hall) VALUES (1, '2024-05-30 12:00:00', 'Hall 1')")
	var showtimeID int
	err := db.Dbx.Get(&showtimeID, "SELECT showtime_id FROM showtimes WHERE movie_id=1 AND showtime='2024-05-30 12:00:00' AND hall='Hall 1'")
	if err != nil {
		t.Fatalf("Failed to get showtime ID: %v", err)
	}

	showtimeInput := models.ShowtimeInput{
		MovieID:  1,
		Showtime: "2024-05-30 20:00",
		Hall:     "Hall 1",
	}
	jsonValue, _ := json.Marshal(showtimeInput)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/showtimes/"+strconv.Itoa(showtimeID), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var showtime models.Showtime
	err = json.Unmarshal(w.Body.Bytes(), &showtime)
	assert.NoError(t, err)
	assert.Equal(t, 1, showtime.MovieID)
	assert.Equal(t, "2024-05-30 20:00", showtime.Showtime)
	assert.Equal(t, "Hall 1", showtime.Hall)
}

func TestDeleteShowtime(t *testing.T) {
	router := setupRouter()

	db.Dbx.MustExec("INSERT INTO showtimes (movie_id, showtime, hall) VALUES (1, '2024-05-30 12:00:00', 'Hall 1')")
	var showtimeID int
	err := db.Dbx.Get(&showtimeID, "SELECT showtime_id FROM showtimes WHERE movie_id=1 AND showtime='2024-05-30 12:00:00' AND hall='Hall 1'")
	if err != nil {
		t.Fatalf("Failed to get showtime ID: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/showtimes/"+strconv.Itoa(showtimeID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
