package middleware

import "github.com/gin-gonic/gin"

func RateLimit(limit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if RateLimitConfig[ip] > limit {
			c.JSON(429, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}
		RateLimitConfig[ip]++
		c.Next()
	}
}
