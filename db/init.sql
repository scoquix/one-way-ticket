CREATE TABLE users (
                       user_id SERIAL PRIMARY KEY,
                       username VARCHAR(50) NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       email VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE movies (
                        movie_id SERIAL PRIMARY KEY,
                        title VARCHAR(100) NOT NULL,
                        duration INT NOT NULL,
                        genre VARCHAR(50) NOT NULL
);

CREATE TABLE showtimes (
                           showtime_id SERIAL PRIMARY KEY,
                           movie_id INT NOT NULL,
                           showtime TIMESTAMP NOT NULL,
                           hall VARCHAR(50) NOT NULL,
                           FOREIGN KEY (movie_id) REFERENCES movies(movie_id)
);

CREATE TABLE Bookings (
                          booking_id SERIAL PRIMARY KEY,
                          user_id INT NOT NULL,
                          showtime_id INT NOT NULL,
                          seat_number INT NOT NULL CHECK (seat_number > 0 AND seat_number <= 100),
                          FOREIGN KEY (user_id) REFERENCES Users(user_id),
                          FOREIGN KEY (showtime_id) REFERENCES showtimes(showtime_id),
                          UNIQUE (showtime_id, seat_number)
);
