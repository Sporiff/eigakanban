package handlers

import (
	"codeberg.org/sporiff/eigakanban/helpers"
	"codeberg.org/sporiff/eigakanban/services"
	"codeberg.org/sporiff/eigakanban/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ItemsHandler struct {
	itemsService *services.ItemsService
}

func NewItemsHandler(itemsService *services.ItemsService) *ItemsHandler {
	return &ItemsHandler{
		itemsService: itemsService,
	}
}

// GetAllItems returns a paginated list of all items
//
//	@Summary		Get all items
//	@Description	Get all items in a paginated list
//	@Security		BearerAuth
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int	false	"Page"
//	@Param			page_size	query		int	false	"Page size"
//	@Success		200			{object}	types.PaginatedItemsResponse
//	@Failure		500			{object}	types.ErrorResponse
//	@Router			/items [get]
func (h *ItemsHandler) GetAllItems(c *gin.Context) {
	pagination, err := helpers.ValidatePagination(c)
	if err != nil {
		helpers.HandleAPIError(c, err)
		return
	}

	result, err := h.itemsService.GetAllItems(c.Request.Context(), pagination)
	if err != nil {
		helpers.HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetItemByUuid returns an item by UUID
//
//	@Summary		Get item by UUID
//	@Description	Get an item by UUID
//	@Security		BearerAuth
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path		string	true	"Item UUID"
//	@Success		200		{object}	types.ItemsResponse
//	@Failure		400		{object}	types.ErrorResponse
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/items/{uuid} [get]
func (h *ItemsHandler) GetItemByUuid(c *gin.Context) {
	itemUuid := c.Param("uuid")
	item, err := h.itemsService.GetItemByUuid(c.Request.Context(), itemUuid)
	if err != nil {
		helpers.HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": item})
}

// AddItem adds a new item to the system
//
//	@Summary		Add a new item
//	@Description	Add a new item
//	@Tags			items
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		types.AddItemRequest		true	"Item details"
//	@Success		200		{object}	types.ItemsResponse			"Item added successfully"
//	@Failure		400		{object}	types.MissingFieldResponse	"Missing mandatory fields"
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/items [post]
func (h *ItemsHandler) AddItem(c *gin.Context) {
	var req types.AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	item, err := h.itemsService.AddItem(c.Request.Context(), req)
	if err != nil {
		helpers.HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"item": item})
}

// UpdateItem updates item details
//
//	@Summary		Update item details
//	@Description	Update item details by UUID
//	@Security		BearerAuth
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path		string					true	"Item UUID"
//	@Param			body	body		types.UpdateItemRequest	true	"Item details to update"
//	@Success		200		{object}	types.ItemsResponse
//	@Failure		400		{object}	types.ErrorResponse
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/items/{uuid} [patch]
func (h *ItemsHandler) UpdateItem(c *gin.Context) {
	var req types.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	itemUuid := c.Param("uuid")
	item, err := h.itemsService.UpdateItem(c.Request.Context(), itemUuid, req.ItemTitle)
	if err != nil {
		helpers.HandleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": item})
}

// DeleteItem deletes an item from the database by UUID
//
//	@Summary		Delete item
//	@Description	Delete an item by UUID
//	@Security		BearerAuth
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path		string						true	"Board UUID"
//	@Success		200		{object}	types.ItemDeletedResponse	"Board deleted successfully"
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/items/{uuid} [delete]
func (h *ItemsHandler) DeleteItem(c *gin.Context) {
	itemUuid := c.Param("uuid")
	err := h.itemsService.DeleteItem(c.Request.Context(), itemUuid)
	if err != nil {
		helpers.HandleAPIError(c, err)
		return
	}

	responseString := fmt.Sprintf("Item deleted: %s", itemUuid)

	c.JSON(http.StatusOK, gin.H{"success": responseString})
}
