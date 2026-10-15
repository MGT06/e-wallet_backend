package router

import (
	"github.com/MGT06/e-wallet_backend/internal/controller"
	"github.com/MGT06/e-wallet_backend/internal/middleware"
	"github.com/MGT06/e-wallet_backend/internal/repository"
	"github.com/MGT06/e-wallet_backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func UserRouter(router *gin.Engine, pool *pgxpool.Pool, redis *redis.Client) {
	user := router.Group("user")

	ur := repository.NewUSerRepo(pool)
	us := service.NewUserService(ur)
	uc := controller.NewUserContorller(us)

	user.GET("info", middleware.AuthMiddleware(redis), uc.GetUserHeaderUC)
}
