package router

import (
	"fmt"
	"os"

	_ "github.com/MGT06/e-wallet_backend/docs"
	"github.com/MGT06/e-wallet_backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router(router *gin.Engine, pool *pgxpool.Pool, redis *redis.Client) {
	fmt.Println(os.Getenv("ORIGIN"))

	router.Use(middleware.Cors)
	router.Static("public/image", "public/photo_profiles")
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	StatusServer(router)
	AuthRouter(router, pool, redis)
	UserRouter(router, pool, redis)
}
