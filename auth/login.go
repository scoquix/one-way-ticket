package auth

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
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

// Handler struct to handle login requests and interact with DynamoDB
type Handler struct {
	ddb dynamodbiface.DynamoDBAPI
}

// NewHandler creates a new Handler with the provided DynamoDB client
func NewHandler(ddb dynamodbiface.DynamoDBAPI) *Handler {
	return &Handler{ddb: ddb}
}

func (h *Handler) Login(c *gin.Context) {
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
	err = CreateSession(h.ddb, t, ttl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": t})
}

func (h *Handler) AuthenticateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte("secret"), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//verify the token
		if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		_, err = GetSessionForUser(h.ddb, tokenString)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Next()
	}
}
