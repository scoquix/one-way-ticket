package models

type Booking struct {
	BookingID  int `db:"booking_id" json:"booking_id"`
	UserID     int `db:"user_id" json:"user_id"`
	ShowtimeID int `db:"showtime_id" json:"showtime_id"`
	SeatNumber int `db:"seat_number" json:"seat_number"`
}

type BookingInput struct {
	UserID     int `db:"user_id" json:"user_id" binding:"required"`
	ShowtimeID int `db:"showtime_id" json:"showtime_id" binding:"required"`
	SeatNumber int `db:"seat_number" json:"seat_number" binding:"required"`
}
