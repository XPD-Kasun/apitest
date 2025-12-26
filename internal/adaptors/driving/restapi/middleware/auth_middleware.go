package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func isValid(auth string) (bool, string) {
	if auth == "" || !strings.HasPrefix(auth, "bearer") {
		return false, ""
	}
	words := strings.Split(auth, " ")
	if len(words) < 2 {
		return false, ""
	}
	return true, words[1]
}

func BearerAuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		auth := strings.ToLower(c.GetHeader("authorization"))

		if valid, token := isValid(auth); !valid {
			c.String(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return

		} else {
			c.Set("jwtToken", token)
		}
		c.Next()
	}
}
