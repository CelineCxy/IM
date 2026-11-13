package middleware

import (
	"crypto/rand"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Printf("request_id=%s panic=%v\n%s",
					c.GetString("request_id"),
					recovered,
					debug.Stack(),
				)
				c.Abort()
				if c.Writer.Written() {
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "internal server error",
					"data":    nil,
				})
			}
		}()
		c.Next()
	}
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		request_id := rand.Text()
		c.Set("request_id", request_id)
		c.Header("X-Request-ID", request_id)
		c.Next()
	}
}

func TimeDuration() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		c.Set("duration", time.Since(start))
	}
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		duration := c.GetDuration("duration")
		log.Printf(
			"request_id=%s method=%s path=%s status=%d duration=%v",
			c.GetString("request_id"),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
		)
	}
}
