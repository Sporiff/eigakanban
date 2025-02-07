package handlers

import (
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"codeberg.org/sporiff/eigakanban/helpers"
	"codeberg.org/sporiff/eigakanban/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type BoardsHandler struct {
	db *pgxpool.Pool
	q  *queries.Queries
}

func NewBoardsHandler(db *pgxpool.Pool) *BoardsHandler {
	return &BoardsHandler{
		db: db,
		q:  queries.New(db),
	}
}

type BoardsResponse struct {
	UUID        string `json:"uuid" example:"00000000-0000-0000-0000-000000000000"`
	Name        string `json:"name" example:"My queue"`
	Description string `json:"description" example:"A short description"`
}

// GetAllBoards returns a paginated list of all boards
//
//	@Summary		Get all boards
//	@Description	Get all boards in a paginated list
//	@Security		BearerAuth
//	@Tags			boards
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int	false	"Page"
//	@Param			page_size	query		int	false	"Page size"
//	@Success		200			{object}	handlers.GetAllBoards.PaginatedBoardsResponse
//	@Failure		500			{object}	types.ErrorResponse
//	@Router			/boards [get]
func (h *BoardsHandler) GetAllBoards(c *gin.Context) {
	// PaginatedBoardsResponse represents a response containing a list of boards
	//	@Description	a response containing a list of boards and a pagination object
	type PaginatedBoardsResponse struct {
		Pagination types.Pagination `json:"pagination"`
		Boards     []BoardsResponse `json:"boards"`
	}

	page, pageSize, err := helpers.ValidatePagination(c)
	if err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	total, err := h.q.GetBoardsCount(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	boards, err := h.q.GetAllBoards(c.Request.Context(), queries.GetAllBoardsParams{
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

	c.JSON(http.StatusOK, gin.H{"pagination": paginationValues, "boards": boards})
}

// GetBoardsForUser returns a paginated list of all boards
//
//	@Summary		Get all boards for a user
//	@Description	Get all boards for a user in a paginated list
//	@Security		BearerAuth
//	@Tags			boards
//	@Accept			json
//	@Produce		json
//	@Param			uuid		path		string	true "User UUID"
//	@Param			page		query		int		false	"Page"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	handlers.GetAllBoards.PaginatedBoardsResponse
//	@Failure		500			{object}	types.ErrorResponse
//	@Router			/user/{uuid}/boards [get]
func (h *BoardsHandler) GetBoardsForUser(c *gin.Context) {
	userUuid := c.Param("uuid")

	page, pageSize, err := helpers.ValidatePagination(c)
	if err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	pgUuid, err := helpers.ValidateAndConvertUUID(userUuid)
	if err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	total, err := h.q.GetBoardsCountForUser(c.Request.Context(), pgUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	boards, err := h.q.GetBoardsForUser(c.Request.Context(), queries.GetBoardsForUserParams{
		UserUuid: pgUuid,
		Page:     int32(page - 1),
		PageSize: int32(pageSize),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	paginationValues := types.Pagination{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Total:    total,
	}

	c.JSON(http.StatusOK, gin.H{"pagination": paginationValues, "boards": boards})
}

// GetBoardByUuid returns a board by UUID
//
//	@Summary		Get board by UUID
//	@Description	Get a board by UUID
//	@Security		BearerAuth
//	@Tags			boards
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path		string	true	"Board UUID"
//	@Success		200		{object}	BoardsResponse
//	@Failure		400		{object}	types.ErrorResponse
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/boards/{uuid} [get]
func (h *BoardsHandler) GetBoardByUuid(c *gin.Context) {
	boardUuid := c.Param("uuid")
	pgUuid, err := helpers.ValidateAndConvertUUID(boardUuid)
	if err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	board, err := h.q.GetBoardByUuid(c.Request.Context(), pgUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"board": board})
}

// AddBoard adds a new board to the system
//
//	@Summary		Add a new board
//	@Description	Add a new board for a user
//	@Tags			boards
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		handlers.AddBoard.AddBoardRequest	true	"Board details"
//	@Success		200		{object}	BoardsResponse						"Board added successfully"
//	@Failure		400		{object}	handlers.MissingFieldResponse		"Missing mandatory fields"
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/boards [post]
func (h *BoardsHandler) AddBoard(c *gin.Context) {
	// AddBoardRequest represents the request body for adding a board
	//	@Description	A request body for adding a new board
	type AddBoardRequest struct {
		BoardName        string `json:"name" example:"My Queue" binding:"required"`
		BoardDescription string `json:"description" example:"A short description"`
		UserUuid         string `json:"user_uuid" example:"00ca71c5-7c8a-4470-ab47-f962d33c1303" binding:"required"`
	}

	var req AddBoardRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	pgUuid, err := helpers.ValidateAndConvertUUID(req.UserUuid)
	if err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	params := queries.CountBoardsByNameAndOwnerParams{
		Name:     req.BoardName,
		UserUuid: pgUuid,
	}

	existingBoardCount, err := h.q.CountBoardsByNameAndOwner(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existingBoardCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Board with this name already exists for this user"})
		return
	}

	boardParams := queries.AddBoardParams{
		BoardName:        req.BoardName,
		BoardDescription: pgtype.Text{String: req.BoardDescription, Valid: true},
		UserUuid:         pgUuid,
	}

	board, err := h.q.AddBoard(c.Request.Context(), boardParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"board": board})
}

// UpdateBoard updates board details
//
//	@Summary		Update board details
//	@Description	Update board details by UUID
//	@Security		BearerAuth
//	@Tags			boards
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path		string									true	"Board UUID"
//	@Param			body	body		handlers.UpdateBoard.UpdateBoardRequest	true	"Board details to update"
//	@Success		200		{object}	BoardsResponse
//	@Failure		400		{object}	types.ErrorResponse
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/boards/{uuid} [patch]
func (h *BoardsHandler) UpdateBoard(c *gin.Context) {
	// UpdateBoardRequest represents the request body for updating a board
	//	@Description	a request body for updating a board
	type UpdateBoardRequest struct {
		NewName        *string `json:"name" example:"My Board"`
		NewDescription *string `json:"description" example:"A short description"`
	}

	var req UpdateBoardRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	boardUuid := c.Param("uuid")
	pgUuid, err := helpers.ValidateAndConvertUUID(boardUuid)
	if err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	var params = queries.UpdateBoardParams{
		BoardUuid: pgUuid,
	}

	// If no board name is provided, we need to fetch the current board name so it doesn't get overwritten
	if req.NewName != nil {
		params.BoardName = *req.NewName
	} else {
		currentBoard, err := h.q.GetBoardByUuid(c.Request.Context(), pgUuid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch board details"})
			return
		}
		params.BoardName = currentBoard.Name
	}

	// Ensure that any missing fields don't overwrite the current field
	if req.NewDescription != nil {
		params.BoardDescription = pgtype.Text{String: *req.NewDescription, Valid: true}
	} else {
		params.BoardDescription = pgtype.Text{Valid: false}
	}

	// Update the board details
	board, err := h.q.UpdateBoard(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"board": board})
}

// DeleteBoard deletes a board from the database by UUID
//
//	@Summary		Delete board
//	@Description	Delete a board by UUID
//	@Security		BearerAuth
//	@Tags			boards
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path		string										true	"Board UUID"
//	@Success		200		{object}	handlers.DeleteBoard.BoardDeletedResponse	"Board deleted successfully"
//	@Failure		500		{object}	types.ErrorResponse
//	@Router			/boards/{uuid} [delete]
func (h *BoardsHandler) DeleteBoard(c *gin.Context) {
	// BoardDeletedResponse represents a success message for a board deletion
	//	@Description	A success message confirming the board was deleted
	type BoardDeletedResponse struct {
		Message string `json:"success" example:"Board deleted: 77b62cff-0020-43d9-a90c-5d35bff89f7a"`
	}

	boardUuid := c.Param("uuid")
	pgUuid, err := helpers.ValidateAndConvertUUID(boardUuid)
	if err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	err = h.q.DeleteBoard(c.Request.Context(), pgUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseString := BoardDeletedResponse{
		Message: fmt.Sprintf("Board deleted: %s", pgUuid),
	}

	c.JSON(http.StatusOK, gin.H{"result": responseString.Message})
}
