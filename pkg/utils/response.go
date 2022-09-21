package utils

import (
	"github.com/gin-gonic/gin"
)

var (
	UserCreatedMessage  = "user registered successfully"
	UserLoggedInMessage = "user logged in successfully"
)

var (
	UserCreateDuplicateMessage = "duplicate username, use another username"
	UserNotFoundMessage        = "user not found"
)

func SendOKResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}

func SendErrorResponse(c *gin.Context, errCode int, errMessage string) {
	c.AbortWithStatusJSON(errCode, map[string]string{"error": errMessage})
}

func MessageResponse(message interface{}) map[string]interface{} {
	return map[string]interface{}{"message": message}
}

func TokenResponse(message, token string) map[string]interface{} {
	return map[string]interface{}{"message": message, "token": token}
}
