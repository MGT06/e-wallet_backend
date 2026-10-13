package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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
}
