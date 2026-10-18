package repository

import (
	"context"
	"errors"

	msgerr "github.com/MGT06/e-wallet_backend/internal/message"
	"github.com/MGT06/e-wallet_backend/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBTX interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type WebhookRepo struct{}

func NewWebhookRepo() *WebhookRepo {
	return &WebhookRepo{}
}

func (wp *WebhookRepo) GetBalance(ctx context.Context, db DBTX, userId int64) (int64, error) {
	query := "SELECT balance FROM accounts WHERE id = $1 FOR UPDATE"
	args := []any{userId}

	var user models.Account
	err := db.QueryRow(ctx, query, args...).Scan(&user.Balance)

	if err != nil {
		return 0, err
	}

	return user.Balance, err
}

func (wp *WebhookRepo) AddBalance(ctx context.Context, db DBTX, userId, amountPayment int64) error {
	query := "UPDATE accounts SET balance = balance + $2, updated_at = now() WHERE id = $1"
	args := []any{userId, amountPayment}

	tag, err := db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return msgerr.ErrNoRowAffected
	}

	return nil
}

func (wp *WebhookRepo) UpdateStatus(ctx context.Context, db DBTX, arg models.Topup) error {
	builder := squirrel.Update("topups").
		Set("status", arg.Status).
		Set("updated_at", squirrel.Expr("now()")).
		Where(squirrel.Eq{"id": arg.ID, "status": "pending"})

	if arg.Status == "success" {
		builder = builder.
			Set("paid_at", arg.PaidAt)
	}

	query, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	tag, er := db.Exec(ctx, query, args...)
	if er != nil {
		return er
	}

	if tag.RowsAffected() == 0 {
		return msgerr.ErrNoRowAffected
	}

	return nil
}

func (wp *WebhookRepo) GetTopup(ctx context.Context, db DBTX, referenceNo string) (models.Topup, error) {
	query := "SELECT id, reference_no, account_id, amount, total_amount, status::text, xendit_invoice_id FROM topups WHERE reference_no = $1 FOR UPDATE"
	args := []any{referenceNo}

	var tp models.Topup
	err := db.QueryRow(ctx, query, args...).Scan(&tp.ID, &tp.ReferenceNo, &tp.AccountID, &tp.Amount, &tp.TotalAmount, &tp.Status, &tp.XenditID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Topup{}, msgerr.ErrDataNotFound
		}
		return models.Topup{}, err
	}

	return tp, nil
}

func (wp *WebhookRepo) InsertTransaction(ctx context.Context, db DBTX, t models.Transaction) error {
	query := "INSERT INTO transaction (account_id, type, amount, balance_before, balance_after, topup_id) VALUES ($1, 'topup', $2, $3, $4, $5)"
	args := []any{t.AccountID, t.Amount, t.BalanceBefore, t.BalanceAfter, t.TopupID}

	tag, err := db.Exec(ctx, query, args...)

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return msgerr.ErrNoRowAffected
	}

	return nil
}

func (wp *WebhookRepo) InsertWebhook(ctx context.Context, db DBTX, eventId, eventStatus string, payload []byte) (isNew bool, err error) {
	query := "INSERT INTO webhook_events (event_id, event_status, payload) VALUES ($1, $2, $3::jsonb) ON CONFLICT (event_id, event_status) DO NOTHING"
	args := []any{eventId, eventStatus, string(payload)}

	tag, err := db.Exec(ctx, query, args...)

	if err != nil {
		return false, err
	}

	if tag.RowsAffected() == 0 {
		return false, nil
	}
	return true, nil
}
