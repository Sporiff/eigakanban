package helpers

import (
	"codeberg.org/sporiff/eigakanban/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func ValidatePagination(ctx *gin.Context) (*types.Pagination, error) {
	page, err := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 32)
	if err != nil {
		errorMessage := fmt.Sprintf("invalid page parameter %d", page)
		return nil, types.NewAPIError(http.StatusBadRequest, errorMessage)
	}

	pageSize, err := strconv.ParseInt(ctx.DefaultQuery("page_size", "50"), 10, 32)
	if err != nil {
		errorMessage := fmt.Sprintf("invalid page_size parameter %d", pageSize)
		return nil, types.NewAPIError(http.StatusBadRequest, errorMessage)
	}

	result := types.Pagination{
		Page:     int32(page - 1),
		PageSize: int32(pageSize),
	}

	return &result, nil
}
