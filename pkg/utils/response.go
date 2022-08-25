package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	UserCreatedMessage = "user registered successfully"
)

func SendOKResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func SendErrorResponse(c *gin.Context, errCode int, errMessage string) {
	c.AbortWithStatusJSON(errCode, map[string]string{"error": errMessage})
}

