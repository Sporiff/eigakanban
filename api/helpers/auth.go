package helpers

import (
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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
func GenerateAccessToken(user queries.GetExistingUserRow) (string, error) {
	// Generate a JWT token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_uuid": user.Uuid.String(),
		"exp":       time.Now().Add(time.Hour * 1).Unix(),
	})

	// TODO: Look into getting users to set their own key
	accessTokenString, err := accessToken.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", errors.New("error generating access token")
	}

	return accessTokenString, nil
}

// GenerateRefreshToken creates a new refresh token for logged-in users
func GenerateRefreshToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// ValidateUserUuidFromClaims validates the user UUID from the claims is present and of the correct type
func ValidateUserUuidFromClaims(c *gin.Context) (string, error) {
	// If the claim is missing, return an error
	userUuid, exists := c.Get("user_uuid")
	if !exists {
		return "", errors.New("missing user uuid")
	}

	// Verify that the user_uuid is a pgtype.UUID value
	switch v := userUuid.(type) {
	case string:
		return v, nil
	default:
		return "", errors.New("invalid user uuid")
	}
}
