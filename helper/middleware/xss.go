package middleware

import (
	"github.com/gin-gonic/gin"
)

var (
	RateLimitConfig = map[string]int{}
)

// X-XSS-Protection: 1; mode-block can help to prevent XSS attacks
// X-Frame-Options: deny can help to prevent clickjacking attacks
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("X-Frame-Options", "deny")
		c.Next()
	}
}

