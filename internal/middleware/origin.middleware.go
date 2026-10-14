package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

type ClientInfo struct {
	IP        string
	UserAgent string
	Origin    string
}

func Origin(c *gin.Context) {
	client, _ := GetClientInfo(c)

	fmt.Printf("\n\n")
	log.Println("Request Accepted from : ")
	fmt.Println("    IP         : ", client.IP)
	fmt.Println("    Origin     : ", client.Origin)
	fmt.Printf("    User Agent : %s\n\n", client.UserAgent)

	origins := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	for _, o := range strings.Split(origins, ",") {
		trimmed := strings.Trim(o, " \t\n\r\"'")
		if trimmed != "" {
			allowedOrigins = append(allowedOrigins, trimmed)
		}
	}

	if client.Origin == "" {
		c.Next()
		// c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		// 	"status":  false,
		// 	"message": "request must have origin allowed",
		// })
		return
	}

	if !slices.Contains(allowedOrigins, client.Origin) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  false,
			"message": fmt.Sprintf("request %s not allowed", client.Origin),
		})
		return
	}
	c.Next()
}

func GetClientInfo(c *gin.Context) (*ClientInfo, error) {
	ip := c.ClientIP()
	if ip == "" {
		return nil, errors.New("unable to determine client IP")
	}

	userAgent := strings.TrimSpace(c.GetHeader("User-Agent"))
	if userAgent == "" {
		userAgent = "Unknown Device"
	}

	origin := c.GetHeader("Origin")

	return &ClientInfo{
		IP:        ip,
		UserAgent: userAgent,
		Origin:    origin,
	}, nil
}
