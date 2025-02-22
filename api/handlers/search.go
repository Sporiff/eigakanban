package handlers

import (
	"codeberg.org/sporiff/eigakanban/helpers"
	"codeberg.org/sporiff/eigakanban/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SearchHandler struct {
	searchService *services.SearchService
}

func NewSearchHandler(searchService *services.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// SearchMovie searches for movies on TMDB
//
//	@Summary		Search for movies on TMDB
//	@Description	Search for movies on TMDB
//	@Security		BearerAuth
//	@Tags			search
//	@Accept			json
//	@Produce		json
//	@Param			uuid		path		string	true	"List UUID"
//	@Param			q			query		string	true	"Search query"
//	@Param			page		query		int		false	"Page"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	types.PaginatedListItemsResponse
//	@Failure		500			{object}	types.ErrorResponse
//	@Router			/search [get]
func (h *SearchHandler) SearchMovie(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query required"})
		return
	}

	pagination, err := helpers.ValidatePagination(c)
	if err != nil {
		helpers.HandleValidationError(c, err)
		return
	}

	results, err := h.searchService.SearchMovie(&pagination, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
