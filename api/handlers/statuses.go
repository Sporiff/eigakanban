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

func (h *StatusesHandler) AddStatus(c *gin.Context) {
	userUuid := c.Param("uuid")

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
