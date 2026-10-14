package router

import (
	_ "github.com/MGT06/e-wallet_backend/docs"
	"github.com/MGT06/e-wallet_backend/internal/controller"
	"github.com/MGT06/e-wallet_backend/internal/middleware"
	"github.com/MGT06/e-wallet_backend/internal/repository"
	"github.com/MGT06/e-wallet_backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	r   *gin.Engine
	db  *pgxpool.Pool
	rdb *redis.Client
	m   *middleware.Middleware
	ar  *repository.AuthRepo
	as  *service.AuthService
	ac  *controller.AuthContorller
}

func NewRouter(router *gin.Engine, pool *pgxpool.Pool, redis *redis.Client, origins []string) *Router {
	// Auth
	ar := repository.NewAuthRepo()
	as := service.NewAuthService(ar, redis, pool)
	ah := controller.NewAuthContorller(as)

	return &Router{
		r:   router,
		db:  pool,
		rdb: redis,
		m:   middleware.InitMiddleWare(origins),
		ar:  repository.NewAuthRepo(),
		as:  as,
		ac:  ah,
	}
}

func (r *Router) Connect() {
	r.r.Use(r.m.Origin, r.m.Cors)
	r.r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.statusServer()
	r.AuthRouter()
}
