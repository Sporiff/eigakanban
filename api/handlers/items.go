package handlers

import (
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"codeberg.org/sporiff/eigakanban/helpers"
	"codeberg.org/sporiff/eigakanban/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type ItemsHandler struct {
	db *pgxpool.Pool
	q  *queries.Queries
}

func NewItemsHandler(db *pgxpool.Pool) *ItemsHandler {
	return &ItemsHandler{
		db: db,
		q:  queries.New(db),
	}
}

type ItemsResponse struct {
	UUID  string `json:"uuid" example:"00000000-0000-0000-0000-000000000000"`
	Title string `json:"title" example:"Item title"`
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
//	@Success		200			{object}	handlers.GetAllItems.PaginatedItemsResponse
//	@Failure		500			{object}	types.ErrorResponse
//	@Router			/items [get]
func (h *ItemsHandler) GetAllItems(c *gin.Context) {
	// PaginatedItemsResponse represents a response containing a list of items
	//	@Description	a response containing a list of items and a pagination object
	type PaginatedItemsResponse struct {
		Pagination types.Pagination `json:"pagination"`
		Items      []ItemsResponse  `json:"items"`
	}

	page, pageSize, err := helpers.ValidatePagination(c)
	if err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	total, err := h.q.GetItemsCount(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items, err := h.q.GetAllItems(c.Request.Context(), queries.GetAllItemsParams{
		Page:     int32(page - 1),
		PageSize: int32(pageSize),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	paginationValues := types.Pagination{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Total:    total,
	}

	c.JSON(http.StatusOK, gin.H{"pagination": paginationValues, "items": items})
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
//	@Success		200		{object}	ItemsResponse
//	@Failure		400		{object}	types.ErrorResponse
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/items/{uuid} [get]
func (h *ItemsHandler) GetItemByUuid(c *gin.Context) {
	itemUuid := c.Param("uuid")
	pgUuid, err := helpers.ValidateAndConvertUUID(itemUuid)
	if err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	item, err := h.q.GetItemByUuid(c.Request.Context(), pgUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
//	@Param			body	body		handlers.AddItem.AddItemRequest	true	"Item details"
//	@Success		200		{object}	ItemsResponse					"Item added successfully"
//	@Failure		400		{object}	handlers.MissingFieldResponse	"Missing mandatory fields"
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/items [post]
func (h *ItemsHandler) AddItem(c *gin.Context) {
	// AddItemRequest represents the request body for adding a new item
	//	@Description	A request body for adding a new item
	type AddItemRequest struct {
		ItemTitle string `json:"title" example:"Item title"`
	}

	var req AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	item, err := h.q.AddItem(c.Request.Context(), req.ItemTitle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
//	@Param			uuid	path		string									true	"Item UUID"
//	@Param			body	body		handlers.UpdateItem.UpdateItemRequest	true	"Item details to update"
//	@Success		200		{object}	ItemsResponse
//	@Failure		400		{object}	types.ErrorResponse
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/items/{uuid} [patch]
func (h *ItemsHandler) UpdateItem(c *gin.Context) {
	type UpdateItemRequest struct {
		ItemTitle string `json:"title" example:"Item title"`
	}

	var req UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	itemUuid := c.Param("uuid")
	pgUuid, err := helpers.ValidateAndConvertUUID(itemUuid)
	if err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	var params = queries.UpdateItemParams{
		ItemTitle: req.ItemTitle,
	}

	if req.ItemTitle != "" {
		params.ItemTitle = req.ItemTitle
	} else {
		currentItem, err := h.q.GetItemByUuid(c.Request.Context(), pgUuid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		params.ItemTitle = currentItem.Title
	}

	item, err := h.q.UpdateItem(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
//	@Param			uuid	path		string									true	"Board UUID"
//	@Success		200		{object}	handlers.DeleteItem.ItemDeletedResponse	"Board deleted successfully"
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/items/{uuid} [delete]
func (h *ItemsHandler) DeleteItem(c *gin.Context) {
	// ItemDeletedResponse represents a success message for an item deletion
	//	@Description	A success message confirming the item was deleted
	type ItemDeletedResponse struct {
		Message string `json:"success" example:"Item deleted: 77b62cff-0020-43d9-a90c-5d35bff89f7a"`
	}

	itemUuid := c.Param("uuid")
	pgUuid, err := helpers.ValidateAndConvertUUID(itemUuid)
	if err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	err = h.q.DeleteItem(c.Request.Context(), pgUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseString := ItemDeletedResponse{
		Message: fmt.Sprintf("Item deleted: %s", pgUuid),
	}

	c.JSON(http.StatusOK, gin.H{"result": responseString.Message})
}
