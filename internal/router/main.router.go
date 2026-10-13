package router

import (
	_ "github.com/MGT06/e-wallet_backend/docs"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	router *gin.Engine
	pool   *pgxpool.Pool
	redis  *redis.Client
}

func NewRouter(router *gin.Engine, pool *pgxpool.Pool, redis *redis.Client) *Router {
	return &Router{
		router: router,
		pool:   pool,
		redis:  redis,
	}
}

func (r *Router) Connect() {
	r.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.statusServer()
}
