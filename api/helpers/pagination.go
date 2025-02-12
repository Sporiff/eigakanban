package helpers

import (
	"codeberg.org/sporiff/eigakanban/types"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func ValidatePagination(c *gin.Context) (types.Pagination, error) {
	var pagination = types.Pagination{}

	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	if err != nil {
		errorMessage := fmt.Sprintf("invalid page parameter %d", page)
		return pagination, errors.New(errorMessage)
	}

	pageSize, err := strconv.ParseInt(c.DefaultQuery("page_size", "50"), 10, 32)
	if err != nil {
		errorMessage := fmt.Sprintf("invalid page_size parameter %d", pageSize)
		return pagination, errors.New(errorMessage)
	}

	pagination.Page = int32(page - 1)
	pagination.PageSize = int32(pageSize)

	return pagination, nil
}
