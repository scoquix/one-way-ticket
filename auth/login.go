package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
	"time"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// perform authentication here
	if username != "admin" && password != "password" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	// TODO - check if user exists in users database

	// set TTL for session
	ttl := time.Now().Add(10 * time.Minute).Unix()

	// set claims
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: ttl,
		},
	}

	// generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// generate encoded token and send it as a response
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	log.Println(t, reflect.TypeOf(t))
	err = CreateSession(t, ttl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": t})
}
