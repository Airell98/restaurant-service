package helpers

import (
	"restaurant-service/pkg/errs"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetQueryParam(c *gin.Context,key string)(int, errs.MessageErr) {
	value := c.Query(key)

	data, err := strconv.Atoi(value)

	if err != nil {
		return 0, errs.NewBadRequest("invalid params")
	}

	return data, nil
}