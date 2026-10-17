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

func TransactionRouter(router *gin.Engine, pool *pgxpool.Pool, redis *redis.Client) {
	tr := repository.NewTransactionRepo(pool)
	ts := service.NewTransactionService(tr)
	tc := controller.NewTransactionContorller(ts)

	transaction := router.Group("transaction")

	transaction.GET("people", middleware.AuthMiddleware(redis), tc.FindPeopleUC)
	transaction.GET("topup-methode", middleware.AuthMiddleware(redis), tc.PaymentMethodeTC)
}
