package util

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	fmt.Println("getBearerToken() Authorization: ", authHeader)

	return authHeader
}
