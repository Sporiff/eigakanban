package helpers

import (
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"codeberg.org/sporiff/eigakanban/types"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// HashPassword creates a hashed password from a provided password string
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a hashed password to a database value
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateAccessToken generates an access token for the user
func GenerateAccessToken(user interface{}) (string, string, error) {
	var uuid pgtype.UUID
	var superuser bool

	// Use type assertions to handle both types
	switch u := user.(type) {
	case queries.GetExistingUserRow:
		uuid = u.Uuid
		superuser = u.Superuser
	case queries.GetUserByUuidRow:
		uuid = u.Uuid
		superuser = u.Superuser
	default:
		return "", "", errors.New("unsupported user type")
	}

	expiryDate := time.Now().UTC().Add(1 * time.Hour)

	// Generate a JWT token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_uuid":   uuid.String(),
		"superuser":   superuser,
		"expiry_date": expiryDate.Unix(),
	})

	// TODO: Look into getting users to set their own key
	accessTokenString, err := accessToken.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", "", errors.New("error generating access token")
	}

	return accessTokenString, expiryDate.String(), nil
}

// GenerateRefreshToken creates a new refresh token for logged-in users
func GenerateRefreshToken(length int) (*string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	refreshToken := hex.EncodeToString(b)
	return &refreshToken, nil
}

// ValidateUserUuidFromClaims validates the user UUID from the claims is present and of the correct type
func ValidateUserUuidFromClaims(c *gin.Context) (*string, error) {
	// If the claim is missing, return an error
	userUuid, exists := c.Get("user_uuid")
	if !exists {
		return nil, types.NewAPIError(http.StatusBadRequest, "missing user uuid")
	}

	// Verify that the user_uuid is a string value
	switch v := userUuid.(type) {
	case string:
		return &v, nil
	default:
		return nil, types.NewAPIError(http.StatusBadRequest, "invalid user uuid")
	}
}

// GetRefreshToken retrieves the refresh token from the browser cookies.
// Falls back to using the Refresh-Token header for API clients
func GetRefreshToken(c *gin.Context) (*string, error) {
	var refreshToken string

	cookie, err := c.Cookie("refresh_token")
	if err == nil {
		refreshToken = cookie
	} else {
		refreshToken = c.GetHeader("Refresh-Token")
		if refreshToken == "" {
			return nil, types.NewAPIError(http.StatusBadRequest, "missing refresh token")
		}
	}

	return &refreshToken, nil
}
