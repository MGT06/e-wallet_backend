package models

import (
	"time"
)

type Transaction struct {
	ID            int64      `db:"id" json:"id"`
	AccountID     int64      `db:"account_id" json:"account_id"`
	Type          string     `db:"type" json:"type"`
	Amount        int64      `db:"amount" json:"amount"`
	BalanceBefore int64      `db:"balance_before" json:"balance_before"`
	BalanceAfter  int64      `db:"balance_after" json:"balance_after"`
	TopupID       *int64     `db:"topup_id" json:"topup_id"`
	TransferID    *int64     `db:"transfer_id" json:"transfer_id"`
	DeleteAt      *time.Time `db:"delete_at" json:"delete_at"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
}
