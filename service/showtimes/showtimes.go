package showtimes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"one-way-ticket/db"
	"one-way-ticket/models"
)

var log = logrus.New()

const (
	InvalidShowtimeID        = "Invalid showtime ID"
	OverlappingShowtimeError = "Showtime overlaps with an existing showtime in the same hall"
)

func parseShowtime(showtimeStr string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04", showtimeStr)
}

func showtimeOverlap(movieID int, showtime time.Time, hall string) (bool, error) {
	var existingShowtimes []models.Showtime
	query := `SELECT * FROM showtimes WHERE hall = $1 AND showtime BETWEEN $2 AND $3`
	start := showtime.Add(-time.Hour * 3)
	end := showtime.Add(time.Hour * 3)

	err := db.Dbx.Select(&existingShowtimes, query, hall, start, end)
	if err != nil {
		return false, err
	}

	if len(existingShowtimes) == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func GetShowtimes(c *gin.Context) {
	var showtimes []models.Showtime
	err := db.Dbx.Select(&showtimes, "SELECT * FROM showtimes")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, showtimes)
}

func GetShowtime(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": InvalidShowtimeID})
		return
	}

	var showtime models.Showtime
	err = db.Dbx.Get(&showtime, "SELECT * FROM showtimes WHERE showtime_id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, showtime)
}

func CreateShowtime(c *gin.Context) {
	var showtimeInput models.ShowtimeInput
	if err := c.ShouldBindJSON(&showtimeInput); err != nil {
		log.Error("Error binding JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	showtimeTime, err := parseShowtime(showtimeInput.Showtime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid showtime format"})
		return
	}

	overlap, err := showtimeOverlap(showtimeInput.MovieID, showtimeTime, showtimeInput.Hall)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if overlap {
		c.JSON(http.StatusBadRequest, gin.H{"error": OverlappingShowtimeError})
		return
	}

	query := `INSERT INTO showtimes (movie_id, showtime, hall) VALUES (:movie_id, :showtime, :hall) RETURNING showtime_id`
	rows, err := db.Dbx.NamedQuery(query, &showtimeInput)
	if err != nil {
		log.Error("Error inserting showtime: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var showtime models.Showtime
	if rows.Next() {
		err = rows.Scan(&showtime.ShowtimeID)
		if err != nil {
			log.Error("Error retrieving new showtime ID: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	showtime.MovieID = showtimeInput.MovieID
	showtime.Showtime = showtimeInput.Showtime
	showtime.Hall = showtimeInput.Hall

	log.Info("Showtime created successfully with ID:", showtime.ShowtimeID)
	c.JSON(http.StatusCreated, showtime)
}

func UpdateShowtime(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": InvalidShowtimeID})
		return
	}

	var showtimeInput models.ShowtimeInput
	if err := c.BindJSON(&showtimeInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	showtimeTime, err := parseShowtime(showtimeInput.Showtime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid showtime format"})
		return
	}

	overlap, err := showtimeOverlap(showtimeInput.MovieID, showtimeTime, showtimeInput.Hall)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if overlap {
		c.JSON(http.StatusBadRequest, gin.H{"error": OverlappingShowtimeError})
		return
	}

	showtime := models.Showtime{
		ShowtimeID: id,
		MovieID:    showtimeInput.MovieID,
		Showtime:   showtimeInput.Showtime,
		Hall:       showtimeInput.Hall,
	}

	_, err = db.Dbx.NamedExec("UPDATE showtimes SET movie_id=:movie_id, showtime=:showtime, hall=:hall WHERE showtime_id=:showtime_id", &showtime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, showtime)
}

func DeleteShowtime(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": InvalidShowtimeID})
		return
	}

	_, err = db.Dbx.Exec("DELETE FROM showtimes WHERE showtime_id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}
