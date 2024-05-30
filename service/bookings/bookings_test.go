package bookings

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
	"os"
	"strconv"
	"testing"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/bookings", GetBookings)
	r.GET("/bookings/:id", GetBooking)
	r.POST("/bookings", CreateBooking)
	r.PUT("/bookings/:id", UpdateBooking)
	r.DELETE("/bookings/:id", DeleteBooking)
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
	_, err = db.Dbx.Exec("INSERT INTO users (username, password, email) VALUES ('testuser', 'password', 'test@example.com')")
	if err != nil {
		panic(err)
	}
	_, err = db.Dbx.Exec("INSERT INTO showtimes (movie_id, showtime, hall) VALUES (1, '2024-05-30 12:00:00', 'Hall 1')")
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

func TestGetBookings(t *testing.T) {
	router := setupRouter()

	// Insert test data
	db.Dbx.MustExec("INSERT INTO bookings (user_id, showtime_id, seat_number) VALUES (1, 1, 1)")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/bookings", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var bookings []models.Booking
	err := json.Unmarshal(w.Body.Bytes(), &bookings)
	assert.NoError(t, err)
	assert.NotEmpty(t, bookings)
}

func TestGetBooking(t *testing.T) {
	router := setupRouter()

	// Insert test data
	db.Dbx.MustExec("INSERT INTO bookings (user_id, showtime_id, seat_number) VALUES (1, 1, 2)")

	var bookingID int
	err := db.Dbx.Get(&bookingID, "SELECT booking_id FROM bookings WHERE user_id=1 AND showtime_id=1 AND seat_number=2")
	if err != nil {
		t.Fatalf("Failed to get booking ID: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/bookings/"+strconv.Itoa(bookingID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var booking models.Booking
	err = json.Unmarshal(w.Body.Bytes(), &booking)
	assert.NoError(t, err)
	assert.Equal(t, 1, booking.UserID)
	assert.Equal(t, 1, booking.ShowtimeID)
	assert.Equal(t, 2, booking.SeatNumber)
}

func TestCreateBooking(t *testing.T) {
	router := setupRouter()

	bookingInput := models.BookingInput{
		UserID:     1,
		ShowtimeID: 1,
		SeatNumber: 3,
	}
	jsonValue, _ := json.Marshal(bookingInput)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var booking models.Booking
	err := json.Unmarshal(w.Body.Bytes(), &booking)
	assert.NoError(t, err)
	assert.Equal(t, 1, booking.UserID)
	assert.Equal(t, 1, booking.ShowtimeID)
	assert.Equal(t, 3, booking.SeatNumber)
}

func TestCreateBookingOverlap(t *testing.T) {
	router := setupRouter()

	db.Dbx.MustExec("INSERT INTO bookings (user_id, showtime_id, seat_number) VALUES (1, 1, 4)")

	bookingInput := models.BookingInput{
		UserID:     1,
		ShowtimeID: 1,
		SeatNumber: 4,
	}
	jsonValue, _ := json.Marshal(bookingInput)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), OverlappingSeatError)
}

func TestUpdateBooking(t *testing.T) {
	router := setupRouter()

	db.Dbx.MustExec("INSERT INTO bookings (user_id, showtime_id, seat_number) VALUES (1, 1, 10)")

	var bookingID int
	err := db.Dbx.Get(&bookingID, "SELECT booking_id FROM bookings WHERE user_id=1 AND showtime_id=1 AND seat_number=10")
	if err != nil {
		t.Fatalf("Failed to get booking ID: %v", err)
	}

	fmt.Println(bookingID)

	bookingInput := models.BookingInput{
		UserID:     1,
		ShowtimeID: 1,
		SeatNumber: 6,
	}
	jsonValue, _ := json.Marshal(bookingInput)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/bookings/"+strconv.Itoa(bookingID), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var booking models.Booking
	err = json.Unmarshal(w.Body.Bytes(), &booking)
	assert.NoError(t, err)
	assert.Equal(t, 1, booking.UserID)
	assert.Equal(t, 1, booking.ShowtimeID)
	assert.Equal(t, 6, booking.SeatNumber)
}

func TestDeleteBooking(t *testing.T) {
	router := setupRouter()

	db.Dbx.MustExec("INSERT INTO bookings (user_id, showtime_id, seat_number) VALUES (1, 1, 9)")

	var bookingID int
	err := db.Dbx.Get(&bookingID, "SELECT booking_id FROM bookings WHERE user_id=1 AND showtime_id=1 AND seat_number=9")
	if err != nil {
		t.Fatalf("Failed to get booking ID: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/bookings/"+strconv.Itoa(bookingID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
