package middleware

import (
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"strings"
	"time"
)

type AuthMiddlewareHandler struct {
	db *pgxpool.Pool
	q  *queries.Queries
}

func NewAuthMiddlewareHandler(db *pgxpool.Pool) *AuthMiddlewareHandler {
	return &AuthMiddlewareHandler{
		db: db,
		q:  queries.New(db),
	}
}

func (h *AuthMiddlewareHandler) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the access token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		// Extract the token from the Bearer string
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		// Parse and validate the access token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil // TODO: Look into using the .env file for this
		})

		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
				return
			}
			// If the token is malformed or invalid, reject the request
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if !token.Valid {
			// If the token is invalid reject the request
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Check if the token is expired
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			exp := claims["exp"].(float64)
			if time.Now().Unix() > int64(exp) {
				// The token has expired. Get a new token using the refresh token
				h.handleExpiredToken(c)
				return
			}
		}

		// The token is valid
		c.Next()
	}
}

func (h *AuthMiddlewareHandler) handleExpiredToken(c *gin.Context) {
	// TODO: look up how to handle this with headers and cookies
	refreshToken := c.GetHeader("Refresh-Token")
	if refreshToken == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Refresh token missing"})
		return
	}

	// Validate the refresh token
	fetchedToken, err := h.q.GetRefreshTokenByToken(c.Request.Context(), refreshToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	if time.Now().After(fetchedToken.ExpiresAt.Time) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
		return
	}

	userId := pgtype.Int8{
		Int64: fetchedToken.UserID,
		Valid: true,
	}

	fetchedUser, err := h.q.GetUserById(c.Request.Context(), userId)
	if fetchedUser == (queries.GetUserByIdRow{}) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
	}

	// Refresh token is valid, issue a new access token
	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_uuid": fetchedUser.Uuid.String(),
		"exp":       time.Now().Add(time.Hour * 1).Unix(), // Keep access key's life short
	})

	newAccessTokenString, err := newAccessToken.SignedString([]byte("your-secret-key"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to generate new access token"})
		return
	}

	// Attach the new access token to the response headers
	c.Header("New-Access-Token", newAccessTokenString)
	c.Next()
}
