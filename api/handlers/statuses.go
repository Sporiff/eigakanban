package handlers

import (
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"codeberg.org/sporiff/eigakanban/helpers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type StatusHandler struct {
	db *pgxpool.Pool
	q  *queries.Queries
}

func NewStatusHandler(db *pgxpool.Pool) *StatusHandler {
	return &StatusHandler{
		db: db,
		q:  queries.New(db),
	}
}

func (h *StatusHandler) AddStatus(c *gin.Context) {
	// AddStatusRequest represents the request body for creating a status
	// @Description A request body for adding a new status
	type AddStatusRequest struct {
		StatusLabel string `json:"status_label" example:"test" binding:"required"`
		UserUUID    string `json:"user" example:"77b62cff-0020-43d9-a90c-5d35bff89f7a" binding:"required"`
	}

	// TODO implement auth and get UUID from claims

	var req AddStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	helpers.HandleErrors(c, func() error {
		pgUuid, err := helpers.ValidateAndConvertUUID(req.UserUUID)
		if err != nil {
			return err
		}

		params := queries.AddStatusParams{
			StatusLabel: pgtype.Text{String: req.StatusLabel, Valid: true},
			UserUuid:    pgUuid,
		}

		status, err := h.q.AddStatus(c.Request.Context(), params)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, gin.H{"status": status})
		return nil
	})
}
