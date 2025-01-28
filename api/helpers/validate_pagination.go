package helpers

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func ValidatePagination(c *gin.Context) (int64, int64, error) {
	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	if err != nil {
		return 0, 0, err
	}

	pageSize, err := strconv.ParseInt(c.DefaultQuery("page_size", "50"), 10, 32)
	if err != nil {
		return 0, 0, err
	}

	return page, pageSize, nil
}
