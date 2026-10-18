package models

import (
	"time"
)

type Topup struct {
	ID              int64      `db:"id"`
	ReferenceNo     string     `db:"reference_no"`
	AccountID       int64      `db:"account_id"`
	PaymentMethodID int        `db:"payment_method_id"`
	Amount          int64      `db:"amount"`
	Tax             int64      `db:"tax"`
	TotalAmount     int64      `db:"total_amount"`
	XenditInvoiceID string     `db:"xendit_invoice_id"`
	PaymentUrl      string     `db:"payment_url"`
	ExpiresAt       string     `db:"expires_at"`
	Status          string     `db:"status"`
	PaidAt          *time.Time `db:"paid_at"`
	CreatedAt       time.Time  `db:"created_at"`
}
