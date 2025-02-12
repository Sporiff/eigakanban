package handlers

import (
	"codeberg.org/sporiff/eigakanban/helpers"
	"codeberg.org/sporiff/eigakanban/services"
	"codeberg.org/sporiff/eigakanban/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type StatusesHandler struct {
	statusesService *services.StatusesService
}

func NewStatusesHandler(statusesService *services.StatusesService) *StatusesHandler {
	return &StatusesHandler{
		statusesService: statusesService,
	}
}

// AddStatus adds a status to the database
//
//	@Summary		Add a new status
//	@Description	Add a new status
//	@Tags			statuses
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		types.AddStatusRequest		true	"Status details"
//	@Success		200		{object}	types.StatusesResponse		"Status added successfully"
//	@Failure		400		{object}	types.MissingFieldResponse	"Missing mandatory fields"
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/statuses [post]
func (h *StatusesHandler) AddStatus(c *gin.Context) {
	userUuid, err := helpers.ValidateUserUuidFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req types.AddStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	status, err := h.statusesService.AddStatus(c.Request.Context(), req, userUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}

// GetStatusesForUser fetches all statuses belonging to the authenticated user
//
//	@Summary		Fetch all statuses
//	@Description	Fetch all statuses as a paginated list
//	@Tags			statuses
//	@Security		BearerAuth
//	@Accept			json
//	@Param			page		query		int								false	"Page"
//	@Param			page_size	query		int								false	"Page size"
//	@Success		200			{object}	types.PaginatedStatusesResponse	"Statuses"
//	@Failure		500			{object}	types.ErrorResponse
//	@Router			/statuses [get]
func (h *StatusesHandler) GetStatusesForUser(c *gin.Context) {
	userUuid, err := helpers.ValidateUserUuidFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you are not authorized to access this resource"})
		return
	}

	pagination, err := helpers.ValidatePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	statuses, updatedPagination, err := h.statusesService.GetStatusesForUser(c.Request.Context(), userUuid, &pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pagination": updatedPagination, "statuses": statuses})
}
