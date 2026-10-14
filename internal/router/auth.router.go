package router

import (
	"github.com/MGT06/e-wallet_backend/internal/controller"
	"github.com/MGT06/e-wallet_backend/internal/repository"
	"github.com/MGT06/e-wallet_backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func AuthRouter(router *gin.Engine, pool *pgxpool.Pool, redis *redis.Client) {
	auth := router.Group("auth")

	ar := repository.NewAuthRepo(pool)
	as := service.NewAuthService(ar, redis)
	ac := controller.NewAuthContorller(as)

	auth.POST("register", ac.RegisterAH)
}
