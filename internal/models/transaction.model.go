package models

import (
	"time"
)

type Transaction struct {
	ID            int64      `db:"id"`
	AccountID     int64      `db:"account_id"`
	Type          string     `db:"type"`
	Amount        int64      `db:"amount"`
	BalanceBefore int64      `db:"balance_before"`
	BalanceAfter  int64      `db:"balance_after"`
	TopupID       *int64     `db:"topup_id"`
	TransferID    *int64     `db:"transfer_id"`
	DeleteAt      *time.Time `db:"delete_at"`
	CreatedAt     time.Time  `db:"created_at"`
}

type HistoryMDL struct {
	ID        int64   `db:"id"`
	Type      string  `db:"type"`
	Amount    int64   `db:"amount"`
	Fullname  string  `db:"fullname"`
	ImageURL  *string `db:"image_url"`
	Phone     *string `db:"phone"`
	TotalData float64 `db:"total"`
}
