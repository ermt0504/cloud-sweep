package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Logger returns a gin middleware for logging requests
func Logger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health", "/ready"},
	})
}

// CORS returns a gin middleware for handling CORS
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestID returns a gin middleware that injects a request ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// Timeout returns a gin middleware that sets a request timeout
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Note: For actual timeout handling, consider using context.WithTimeout
		// This is a simplified version
		c.Set("timeout", timeout)
		c.Next()
	}
}

// RateLimit returns a gin middleware for rate limiting
// Note: For production, use a proper rate limiter like golang.org/x/time/rate
func RateLimit(requestsPerSecond int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Placeholder for rate limiting logic
		// In production, implement proper rate limiting with Redis or in-memory store
		c.Next()
	}
}

// Auth returns a gin middleware for authentication
// Note: This is a placeholder - implement proper auth (JWT, OAuth, etc.)
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Placeholder for authentication logic
		// In production, validate JWT tokens, API keys, etc.

		// For now, just check for Authorization header presence
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
