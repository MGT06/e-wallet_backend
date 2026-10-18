package router

import (
	"os"

	"github.com/MGT06/e-wallet_backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router(router *gin.Engine, pool *pgxpool.Pool, redis *redis.Client, resty *resty.Client) {
	router.Use(middleware.Cors)
	router.Static(os.Getenv("PROFILE_URL"), "public/photo_profiles")
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	StatusServer(router)
	AuthRouter(router, pool, redis)
	UserRouter(router, pool, redis)
	TransactionRouter(router, pool, redis, resty)
	WebhookRouter(router, pool)
}
