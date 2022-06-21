package middleware

import (
	"fmt"
	"golang-restaurant-management/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("no authorization header provided")})
			c.Abort()
			return
		}
		clamis, err := helper.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("email", clamis.Email)
		c.Set("first_name", clamis.First_name)
		c.Set("last_name", clamis.Last_name)
		c.Set("uid", clamis.Uid)
		c.Next()
	}
}
