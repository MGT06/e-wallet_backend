package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) Cors(c *gin.Context) {
	origin := c.GetHeader("Origin")
	log.Printf("Origin : %s\n\n", origin)
	c.Header("Access-Control-Allow-Origin", origin)

	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Header("Access-Control-Allow-Methods", "DELETE, OPTIONS, PATCH")
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}
