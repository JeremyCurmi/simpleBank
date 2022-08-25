package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/go_projects_jer/simple_bank/pkg/utils"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := utils.TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
