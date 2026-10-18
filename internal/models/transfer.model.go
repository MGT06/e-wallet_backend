package models

import (
	"time"
)

type Transfer struct {
	ID          int64     `db:"id" json:"id"`
	ReferenceNo string    `db:"reference_no" json:"reference_no"`
	SenderID    int64     `db:"sender_id" json:"sender_id"`
	ReceiverID  int64     `db:"receiver_id" json:"receiver_id"`
	Amount      int64     `db:"amount" json:"amount"`
	Notes       *string   `db:"notes" json:"notes"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type Recipient struct {
	ID       int64
	Fullname string
	ImageURL *string
	Phone    *string
}

type TransferDetail struct {
	Transfer
	Sender   Recipient
	Receiver Recipient
}