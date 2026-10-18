package router

import (
	"github.com/MGT06/e-wallet_backend/internal/controller"
	"github.com/MGT06/e-wallet_backend/internal/middleware"
	"github.com/MGT06/e-wallet_backend/internal/repository"
	"github.com/MGT06/e-wallet_backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func TransactionRouter(router *gin.Engine, pool *pgxpool.Pool, redis *redis.Client, resty *resty.Client) {
	tr := repository.NewTransactionRepo()
	ts := service.NewTransactionService(tr, resty, pool)
	tc := controller.NewTransactionContorller(ts)

	transaction := router.Group("transaction")

	transaction.GET("people", middleware.AuthMiddleware(redis), tc.FindPeopleUC)
	transaction.GET("topup-methode", middleware.AuthMiddleware(redis), tc.PaymentMethodeTC)
	transaction.POST("topup", middleware.AuthMiddleware(redis), tc.TopupsUC)
	transaction.POST("transfer", middleware.AuthMiddleware(redis), tc.TransferUC)
}
