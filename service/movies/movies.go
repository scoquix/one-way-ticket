package movies

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"one-way-ticket/db"
	"one-way-ticket/models"
)

var log = logrus.New()

const (
	InvalidMovieId = "Invalid movie ID"
)

func GetMovies(c *gin.Context) {
	var movies []models.Movie
	err := db.Dbx.Select(&movies, "SELECT * FROM movies")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movies)
}

func GetMovie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": InvalidMovieId})
		return
	}

	var movie models.Movie
	err = db.Dbx.Get(&movie, "SELECT * FROM movies WHERE movie_id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movie)
}

func CreateMovie(c *gin.Context) {
	var movieInput models.MovieInput
	if err := c.ShouldBindJSON(&movieInput); err != nil {
		log.Error("Error binding JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO movies (title, duration, genre) VALUES (:title, :duration, :genre) RETURNING movie_id`
	rows, err := db.Dbx.NamedQuery(query, &movieInput)
	if err != nil {
		log.Error("Error inserting movie: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	movie := models.Movie{
		Title:    movieInput.Title,
		Duration: movieInput.Duration,
		Genre:    movieInput.Genre,
	}

	if rows.Next() {
		err = rows.Scan(&movie.MovieID)
		if err != nil {
			log.Error("Error retrieving new movie ID: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	log.Info("Movie created successfully with ID:", movie.MovieID)
	c.JSON(http.StatusCreated, movie)
}

func UpdateMovie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": InvalidMovieId})
		return
	}

	var movieInput models.MovieInput
	if err := c.BindJSON(&movieInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movie := models.Movie{
		MovieID:  id,
		Title:    movieInput.Title,
		Duration: movieInput.Duration,
		Genre:    movieInput.Genre,
	}

	_, err = db.Dbx.NamedExec("UPDATE movies SET title=:title, duration=:duration, genre=:genre WHERE movie_id=:movie_id", &movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movie)
}

func DeleteMovie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": InvalidMovieId})
		return
	}

	_, err = db.Dbx.Exec("DELETE FROM movies WHERE movie_id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}
