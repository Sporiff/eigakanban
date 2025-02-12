package middleware

import (
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type SuperUserMiddlewareHandler struct {
	db *pgxpool.Pool
	q  *queries.Queries
}

func NewSuperUserMiddlewareHandler(db *pgxpool.Pool) *SuperUserMiddlewareHandler {
	return &SuperUserMiddlewareHandler{
		db: db,
		q:  queries.New(db),
	}
}

func (h *SuperUserMiddlewareHandler) SuperUserStatusRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		superuser, exists := c.Get("superuser")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "you are not authorized to access this resource"})
		}

		if !superuser.(bool) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "you are not authorized to access this resource"})
		}

		c.Next()
	}
}
