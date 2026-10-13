package router

import (
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
	r.statusServer()
	r.router.GET("/docs/v1/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
