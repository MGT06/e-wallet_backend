package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/MGT06/e-wallet_backend/internal/dto"
	msgerr "github.com/MGT06/e-wallet_backend/internal/message"
	"github.com/MGT06/e-wallet_backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type Payload struct {
	ID        int64
	ExpiresIn int64
}

func AuthMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.GetHeader("Authorization")
		if bearer == "" {
			log.Printf("[AuthMiddleware] GetHeader  Authorization : from IP: %s\n", c.ClientIP())
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ResponseError{
				Success: false,
				Message: "please login first",
			})
			return
		}

		result := strings.Split(bearer, " ")
		if len(result) != 2 {
			log.Printf("[AuthMiddleware] Error: malformed Authorization header length (%d parts)\n", len(result))
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ResponseError{
				Success: false,
				Message: "Invalid authorization header format",
			})
			return
		}

		if result[0] != "Bearer" {
			log.Printf("[AuthMiddleware] Error: invalid scheme '%s', expected 'Bearer'\n", result[0])
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ResponseError{
				Success: false,
				Message: "Authorization scheme must be Bearer",
			})
			return
		}

		var token pkg.JWTClaims
		err := token.DecodeToken(result[1])
		if err != nil {
			log.Printf("[AuthMiddleware] DecodeToken : token validation failed: %v\n", err)
			if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenInvalidIssuer) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ResponseError{
					Success: false,
					Message: "Token is invalid or expired",
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ResponseError{
				Success: false,
				Message: "a system error has accurred",
			})
			return
		}

		key := fmt.Sprintf("ewallet:token:%d", token.ID)

		_, err = rdb.Get(c, key).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				log.Printf("[AuthMiddleware] Redis Get key %s : %s", key, err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ResponseError{
					Success: false,
					Message: msgerr.ErrTokenExpired.Error(),
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ResponseError{
				Success: false,
				Message: "a system error has accurred",
			})
			return
		}

		c.Set("ID", token.ID)
		c.Next()
	}
}
