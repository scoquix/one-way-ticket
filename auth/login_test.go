package auth

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"one-way-ticket/dynamo"
	"one-way-ticket/mocks"
	"one-way-ticket/models"
	"strings"
	"testing"
	"time"
)

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

func TestLoginAdminUser(t *testing.T) {
	// Create a new Handler with the mock client
	mockSvc := new(mocks.MockDynamoDBClient)
	mockSvc.On("PutItem", mock.MatchedBy(func(input *dynamodb.PutItemInput) bool {
		return input.TableName != nil && *input.TableName == dynamo.TableName &&
			input.Item["token"].S != nil && len(*input.Item["token"].S) > 0
	})).Return(&dynamodb.PutItemOutput{}, nil)

	// Create a new Handler with the mock client
	handler := NewHandler(mockSvc)
	router := gin.Default()
	router.POST("/login", handler.Login)

	username := "admin"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", strings.NewReader("username="+username+"&password=password"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	token := response["token"]
	assert.NotEmpty(t, token)

	claims := &models.Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)
	assert.Equal(t, username, claims.Username)
	assert.WithinDuration(t, time.Now().Add(time.Minute*15), time.Unix(claims.ExpiresAt, 0), 5*time.Second)
}
