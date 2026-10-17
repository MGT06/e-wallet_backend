package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cors(c *gin.Context) {
	allowed := []string{"http://localhost:5173", "http://192.168.0.186:5001", "http://192.168.0.186:5001", "http://localhost:5174"}
	for _, o := range allowed {
		if c.GetHeader("Origin") == o {
			c.Header("Access-Control-Allow-Origin", o)
		}
	}

	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Header("Access-Control-Allow-Methods", "POST, DELETE, OPTIONS, PATCH")
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}
