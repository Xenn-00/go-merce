package middleware

import (
	"net/http"

	token "github.com/Xenn-00/go-merce/tokens"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		ClientToken := c.Request.Header.Get("token")
		if ClientToken == "" {
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"error": "no authorization header",
			})
			return
		}
		claims, err := token.ValidateToken(ClientToken)
		if err != "" {
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"error": err,
			})
			return
		}

		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}
