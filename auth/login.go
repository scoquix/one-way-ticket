package auth

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"one-way-ticket/dynamo"
	"one-way-ticket/models"
	"time"
)

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

	// set TTL for session
	ttl := time.Now().Add(15 * time.Minute).Unix()

	// set claims
	claims := &models.Claims{
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
		return
	}

	err = dynamo.CreateSession(h.ddb, t, ttl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": t})
}
