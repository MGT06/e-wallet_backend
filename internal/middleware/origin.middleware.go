package middleware

import (
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/MGT06/e-wallet_backend/internal/dto"
	"github.com/gin-gonic/gin"
)

func Origin(c *gin.Context) {
	origin := c.GetHeader("Origin")
	origins := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	for _, o := range strings.Split(origins, ",") {
		trimmed := strings.Trim(o, " \t\n\r\"'")
		if trimmed != "" {
			allowedOrigins = append(allowedOrigins, trimmed)
		}
	}

	if origin == "" {
		c.Next()
		return
	}

	if !slices.Contains(allowedOrigins, origin) {
		c.AbortWithStatusJSON(http.StatusForbidden, dto.ResponseError{
			Success: false,
			Message: fmt.Sprintf("request %s not allowed", origin),
		})
		return
	}
	c.Next()
}
