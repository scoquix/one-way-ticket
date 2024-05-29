package models

type Showtime struct {
	ShowtimeID int    `db:"showtime_id" json:"showtime_id"`
	MovieID    int    `db:"movie_id" json:"movie_id"`
	Showtime   string `db:"showtime" json:"showtime"`
	Hall       string `db:"hall" json:"hall"`
}

type ShowtimeInput struct {
	MovieID  int    `db:"movie_id" json:"movie_id" binding:"required"`
	Showtime string `db:"showtime" json:"showtime" binding:"required"`
	Hall     string `db:"hall" json:"hall" binding:"required"`
}
