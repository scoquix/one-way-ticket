package auth

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"one-way-ticket/mocks"
	"strings"
	"testing"
	"time"
)

var mockSvc *mocks.MockDynamoDBClient

func init() {
	mockSvc = new(mocks.MockDynamoDBClient)
	mockSvc.On("GetItem", &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"token": {
				S: aws.String("nonexistent-token"),
			},
		},
	}).Return(&dynamodb.GetItemOutput{
		Item: nil,
	}, nil)

	mockSvc.On("PutItem", &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"token": {
				S: aws.String("nonexistent-token"),
			},
		},
	}).Return(&dynamodb.PutItemOutput{}, nil)
}

func TestLoginUnauthorizedUser(t *testing.T) {
	router := gin.Default()
	// Create a new Handler with the mock client
	handler := NewHandler(&mocks.MockDynamoDBClient{})
	router.POST("/login", handler.Login)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", strings.NewReader("username=John&password=123"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "{\"status\":\"unauthorized\"}")
}

//func TestLoginAdminUser(t *testing.T) {
//	router := gin.Default()
//	// Create a new Handler with the mock client
//	mockSvc := new(mocks.MockDynamoDBClient)
//	mockSvc.On("PutItem", &dynamodb.GetItemInput{
//		TableName: aws.String(tableName),
//		Key: map[string]*dynamodb.AttributeValue{
//			"token": {
//				S: aws.String(""),
//			},
//		},
//	}).Return(&dynamodb.PutItemOutput{}, nil)
//
//	handler := NewHandler(mockSvc)
//	router.POST("/login", handler.Login)
//
//	username := "admin"
//
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("POST", "/login", strings.NewReader("username="+username+"&password=password"))
//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//
//	router.ServeHTTP(w, req)
//
//	assert.Equal(t, http.StatusOK, w.Code)
//
//	var response map[string]string
//	err := json.Unmarshal(w.Body.Bytes(), &response)
//	assert.NoError(t, err)
//	token := response["token"]
//	assert.NotEmpty(t, token)
//
//	claims := &Claims{}
//	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
//		return []byte("secret"), nil
//	})
//
//	assert.NoError(t, err)
//	assert.True(t, parsedToken.Valid)
//	assert.Equal(t, username, claims.Username)
//	assert.WithinDuration(t, time.Now().Add(time.Minute*15), time.Unix(claims.ExpiresAt, 0), 5*time.Second)
//}

func generateTestToken(secret string, expirationTime time.Duration) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(expirationTime).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func TestAuthenticateMiddleware(t *testing.T) {
	router := gin.Default()
	mockSvc = new(mocks.MockDynamoDBClient)
	mockSvc.On("GetItem", &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"token": {
				S: aws.String("admin"),
			},
		},
	}).Return(&dynamodb.GetItemOutput{
		Item: map[string]*AttributeValue{},
	}, nil)
	handler := NewHandler(mockSvc)
	router.Use(handler.AuthenticateMiddleware())
	router.GET("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	t.Run("Valid Token", func(t *testing.T) {
		token, _ := generateTestToken("secret", time.Minute*5)
		req, _ := http.NewRequest("GET", "/users", nil)
		req.Header.Set("Authorization", token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"status": "success"}`, w.Body.String())
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users", nil)
		req.Header.Set("Authorization", "invalid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Expired Token", func(t *testing.T) {
		token, _ := generateTestToken("secret", -time.Minute*5)
		req, _ := http.NewRequest("GET", "/users", nil)
		req.Header.Set("Authorization", token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
