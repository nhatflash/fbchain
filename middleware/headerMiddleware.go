package middleware

import (
	"slices"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FilterConfigurer(allowedOrigins ...string) gin.HandlerFunc {
	return func (c *gin.Context) {
		reqOrigin := c.GetHeader("Origin")

		isAllowed := slices.Contains(allowedOrigins, reqOrigin)
		if isAllowed {
			c.Header("Access-Control-Allow-Origin", reqOrigin)
			c.Header("Vary", "Origin")
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		c.Header("Content-Security-Policy", 
		"default-src 'self'; " + 
		"connect-src *; " +
		"img-src 'self' data:; " +  
		"script-src 'self'; " + 
		"frame-ancestors 'none'")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

		c.SetSameSite(http.SameSiteStrictMode)
		c.Next()
	}
}