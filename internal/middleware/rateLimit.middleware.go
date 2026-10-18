package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MGT06/e-wallet_backend/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimitByIP(rdb *redis.Client, limit int64, time time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("ewallet:rl:%s:%s", c.FullPath(), c.ClientIP())
		n, err := rdb.Incr(c.Request.Context(), key).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ResponseError{
				Success: false,
				Message: "a system error has accurred",
			})
			return
		}
		if n == 1 {
			rdb.Expire(c.Request.Context(), key, time)
		}
		if n > limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, dto.ResponseError{
				Success: false,
				Message: "to many request",
			})
			return
		}
		c.Next()
	}
}
