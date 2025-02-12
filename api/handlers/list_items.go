package handlers

import (
	"codeberg.org/sporiff/eigakanban/helpers"
	"codeberg.org/sporiff/eigakanban/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ListItemsHandler struct {
	listItemsService *services.ListItemsService
}

func NewListItemsHandler(listItemsService *services.ListItemsService) *ListItemsHandler {
	return &ListItemsHandler{
		listItemsService: listItemsService,
	}
}

// GetAllListItems returns a paginated list of all list items
//
//	@Summary		Get all list items
//	@Description	Get all list items in a paginated list
//	@Security		BearerAuth
//	@Tags			list_items
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int	false	"Page"
//	@Param			page_size	query		int	false	"Page size"
//	@Success		200			{object}	types.PaginatedListItemsResponse
//	@Failure		500			{object}	types.ErrorResponse
//	@Router			/list_items [get]
func (h *ListItemsHandler) GetAllListItems(c *gin.Context) {
	pagination, err := helpers.ValidatePagination(c)
	if err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	listItems, updatedPagination, err := h.listItemsService.GetAllListItems(c.Request.Context(), &pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pagination": updatedPagination, "list_items": listItems})
}

// GetListItemsForList returns all items in a list
//
//	@Summary		Get all items in a list
//	@Description	Get all items in a list as a paginated list
//	@Security		BearerAuth
//	@Tags			lists
//	@Accept			json
//	@Produce		json
//	@Param			uuid		path		string	true	"List UUID"
//	@Param			page		query		int		false	"Page"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	types.PaginatedListItemsResponse
//	@Failure		500			{object}	types.ErrorResponse
//	@Router			/lists/{uuid}/items [get]
func (h *ListItemsHandler) GetListItemsForList(c *gin.Context) {
	listUuid := c.Param("uuid")

	pagination, err := helpers.ValidatePagination(c)
	if err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	listItems, updatedPagination, err := h.listItemsService.GetListItemsForList(c.Request.Context(), &pagination, listUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pagination": updatedPagination, "list_items": listItems})
}
