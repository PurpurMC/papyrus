package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/purpurmc/papyrus/types"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Return404(c *gin.Context) {
	c.JSON(404, types.ErrorResponse{
		Error: "Not Found",
		Code: 404,
	})
}
