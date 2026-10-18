package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/MGT06/e-wallet_backend/internal/dto"
	msgerr "github.com/MGT06/e-wallet_backend/internal/message"
	"github.com/MGT06/e-wallet_backend/internal/service"
	"github.com/MGT06/e-wallet_backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type WebhookController struct {
	ws *service.WebhookService
}

func NewWebhookController(ws *service.WebhookService) *WebhookController {
	return &WebhookController{
		ws: ws,
	}
}

func (wc *WebhookController) Xendit(ctx *gin.Context) {
	token := ctx.GetHeader("x-callback-token")

	var body dto.XenditInvoiceRequest
	if err := ctx.ShouldBindBodyWith(&body, binding.JSON); err != nil {
		log.Println(err)
		found, isErr := errors.AsType[validator.ValidationErrors](err)
		if isErr {
			err := utils.ValidationError(found)
			ctx.JSON(http.StatusBadRequest, dto.ResponseError{
				Success: false,
				Message: err,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	rawBody, exists := ctx.Get(gin.BodyBytesKey)
	if !exists {
		log.Println("kosong")

		ctx.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	if err := wc.ws.HandleXenditInvoice(ctx.Request.Context(), body, token, rawBody.([]byte)); err != nil {
		log.Println(err)
		if errors.Is(err, msgerr.ErrInvalidCallbackToken) {
			ctx.JSON(http.StatusUnauthorized, dto.ResponseError{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		if errors.Is(err, msgerr.ErrDuplicateWebhook) {
			ctx.JSON(http.StatusNoContent, dto.ResponseError{})
			return
		}
		if errors.Is(err, msgerr.ErrInvoiceMismatch) {
			ctx.JSON(http.StatusUnprocessableEntity, dto.ResponseError{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		if errors.Is(err, msgerr.ErrMismatchPaidAmount) {
			ctx.JSON(http.StatusUnprocessableEntity, dto.ResponseError{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	ctx.JSON(http.StatusNoContent, dto.ResponseSuccess{})
}
