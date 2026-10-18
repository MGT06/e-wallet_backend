package service

import (
	"context"
	"crypto/subtle"
	"errors"
	"log"
	"os"
	"time"

	"github.com/MGT06/e-wallet_backend/internal/dto"
	msgerr "github.com/MGT06/e-wallet_backend/internal/message"
	"github.com/MGT06/e-wallet_backend/internal/models"
	"github.com/MGT06/e-wallet_backend/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WebhookService struct {
	wr *repository.WebhookRepo
	db *pgxpool.Pool
}

func NewWebhookService(wr *repository.WebhookRepo, db *pgxpool.Pool) *WebhookService {
	return &WebhookService{
		wr: wr,
		db: db,
	}
}

func (ws *WebhookService) HandleXenditInvoice(ctx context.Context, body dto.XenditInvoiceRequest, token string, raw []byte) error {
	if subtle.ConstantTimeCompare([]byte(token), []byte(os.Getenv("XENDIT_WEBHOOK_TOKEN"))) != 1 {
		return msgerr.ErrInvalidCallbackToken
	}

	tx, err := ws.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Println(err.Error())
		}
	}()

	var status string
	switch body.Status {
	case "PAID":
		status = "success"
	case "EXPIRED":
		status = "expired"
	default:
		status = "failed"
	}

	isNew, er := ws.wr.InsertWebhook(ctx, tx, body.ID, status, raw)
	if er != nil {
		return er
	}

	if !isNew {
		return msgerr.ErrDuplicateWebhook
	}

	res, e := ws.wr.GetTopup(ctx, tx, body.ExternalID)
	if e != nil {
		if errors.Is(e, msgerr.ErrDataNotFound) {
			log.Println(e)
			if err := tx.Commit(ctx); err != nil {
				return err
			}
			return nil
		}
		return e
	}

	if *res.XenditID != body.ID {
		return msgerr.ErrInvoiceMismatch
	}

	switch body.Status {
	case "PAID":
		if res.Status == "pending" {
			if res.TotalAmount != int64(*body.PaidAmount) {
				return msgerr.ErrMismatchPaidAmount
			}
			before, err := ws.wr.GetBalance(ctx, tx, res.AccountID)
			if err != nil {
				return err
			}

			paidAt, er := time.Parse(time.RFC3339, body.PaidAt)
			if er != nil {
				paidAt = time.Now()
			}

			if err := ws.wr.UpdateStatus(ctx, tx, models.Topup{
				ID: res.ID,
				Status: status,
				PaidAt: &paidAt,
			}); err != nil {
				return err
			}

			if err := ws.wr.AddBalance(ctx, tx, res.AccountID, res.Amount); err != nil {
				return err
			}

			if err := ws.wr.InsertTransaction(ctx, tx, models.Transaction{
				AccountID:     res.AccountID,
				Amount:        res.Amount,
				BalanceBefore: before,
				BalanceAfter:  before + res.Amount,
				TopupID:       &res.ID,
			}); err != nil {
				return err
			}
		}
	case "EXPIRED":
		if res.Status == "pending" {
			if err := ws.wr.UpdateStatus(ctx, tx, models.Topup{
				ID:     res.ID,
				Status: status,
			}); err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}
