package controller

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func getBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	fmt.Println("verifyBearerToken() Authorization: ", authHeader)
	tokenSlices := strings.Split(authHeader, "Bearer ")
	fmt.Println("verifyBearerTOken() token: ", tokenSlices[1])

	return tokenSlices[1]
}
