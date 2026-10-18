package router

import (
	"github.com/MGT06/e-wallet_backend/internal/controller"
	"github.com/MGT06/e-wallet_backend/internal/repository"
	"github.com/MGT06/e-wallet_backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func WebhookRouter(router *gin.Engine, pool *pgxpool.Pool) {
	wr := repository.NewWebhookRepo()
	ws := service.NewWebhookService(wr, pool)
	wc := controller.NewWebhookController(ws)

	webhook := router.Group("webhook")

	webhook.POST("xendit", wc.Xendit)
}
