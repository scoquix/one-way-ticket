package models

type Movie struct {
	MovieID  int    `db:"movie_id" json:"movie_id"`
	Title    string `db:"title" json:"title"`
	Duration int    `db:"duration" json:"duration"`
	Genre    string `db:"genre" json:"genre"`
}

type MovieInput struct {
	Title    string `json:"title" binding:"required"`
	Duration int    `json:"duration" binding:"required"`
	Genre    string `json:"genre" binding:"required"`
}
