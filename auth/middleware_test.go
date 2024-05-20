package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func generateTestToken(secret string, expirationTime time.Duration) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(expirationTime).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func TestAuthenticateMiddleware(t *testing.T) {
	router := gin.Default()
	router.Use(AuthenticateMiddleware())
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
