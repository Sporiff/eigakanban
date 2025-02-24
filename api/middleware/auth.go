package middleware

import (
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
		token, err := h.extractAuthToken(c)
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
			exp := claims["expiry_date"].(float64)
			if time.Now().Unix() > int64(exp) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
				return
			}
		}

		var userUuid string
		var superUser bool

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			claimsUuid, ok := claims["user_uuid"]
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user UUID in token"})
				return
			}
			if claimsUuid.(string) == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user UUID in token"})
				return
			}

			claimsSuperUser, ok := claims["superuser"]
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}

			superUser = claimsSuperUser.(bool)
			userUuid = claimsUuid.(string)
		}

		// Store the user UUID in the context
		if token.Valid {
			c.Set("user_uuid", userUuid)
			c.Set("superuser", superUser)
		}

		// The token is valid
		c.Next()
	}
}

func (h *AuthMiddlewareHandler) extractAuthToken(c *gin.Context) (*jwt.Token, error) {
	// Extract the access token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization header missing")
	}

	// Extract the token from the Bearer string
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return nil, errors.New("invalid token format")
	}

	// Parse and validate the access token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil // TODO: Look into using the .env file for this
	})
	if err != nil {
		return nil, errors.New("invalid token")
	}

	return token, nil
}
