package bookings

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
	InvalidBookingID          = "Invalid booking ID"
	SeatNumberOutOfRangeError = "Seat number must be between 1 and 100"
	OverlappingSeatError      = "Seat number is already booked for this showtime"
)

func GetBookings(c *gin.Context) {
	var bookings []models.Booking
	err := db.Dbx.Select(&bookings, "SELECT * FROM bookings")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

func GetBooking(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": InvalidBookingID})
		return
	}

	var booking models.Booking
	err = db.Dbx.Get(&booking, "SELECT * FROM bookings WHERE booking_id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, booking)
}

func CreateBooking(c *gin.Context) {
	var bookingInput models.BookingInput
	if err := c.ShouldBindJSON(&bookingInput); err != nil {
		log.Error("Error binding JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if bookingInput.SeatNumber < 1 || bookingInput.SeatNumber > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": SeatNumberOutOfRangeError})
		return
	}

	var count int
	err := db.Dbx.Get(&count, "SELECT COUNT(*) FROM bookings WHERE showtime_id=$1 AND seat_number=$2", bookingInput.ShowtimeID, bookingInput.SeatNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": OverlappingSeatError})
		return
	}

	query := `INSERT INTO bookings (user_id, showtime_id, seat_number) VALUES (:user_id, :showtime_id, :seat_number) RETURNING booking_id`
	rows, err := db.Dbx.NamedQuery(query, &bookingInput)
	if err != nil {
		log.Error("Error inserting booking: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var booking models.Booking
	if rows.Next() {
		err = rows.Scan(&booking.BookingID)
		if err != nil {
			log.Error("Error retrieving new booking ID: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	booking.UserID = bookingInput.UserID
	booking.ShowtimeID = bookingInput.ShowtimeID
	booking.SeatNumber = bookingInput.SeatNumber

	log.Info("Booking created successfully with ID:", booking.BookingID)
	c.JSON(http.StatusCreated, booking)
}

func UpdateBooking(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": InvalidBookingID})
		return
	}

	var bookingInput models.BookingInput
	if err := c.BindJSON(&bookingInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if bookingInput.SeatNumber < 1 || bookingInput.SeatNumber > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": SeatNumberOutOfRangeError})
		return
	}

	var count int
	err = db.Dbx.Get(&count, "SELECT COUNT(*) FROM bookings WHERE showtime_id=$1 AND seat_number=$2 AND booking_id<>$3", bookingInput.ShowtimeID, bookingInput.SeatNumber, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": OverlappingSeatError})
		return
	}

	booking := models.Booking{
		BookingID:  id,
		UserID:     bookingInput.UserID,
		ShowtimeID: bookingInput.ShowtimeID,
		SeatNumber: bookingInput.SeatNumber,
	}

	_, err = db.Dbx.NamedExec("UPDATE bookings SET user_id=:user_id, showtime_id=:showtime_id, seat_number=:seat_number WHERE booking_id=:booking_id", &booking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, booking)
}

func DeleteBooking(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": InvalidBookingID})
		return
	}

	_, err = db.Dbx.Exec("DELETE FROM bookings WHERE booking_id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}
