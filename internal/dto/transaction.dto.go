package dto

import "time"

type FindAllPeople struct {
	ID        int64   `json:"id"`
	Fullname  string  `json:"fullname"`
	ImageURL  *string `json:"image_url"`
	Phone     *string `json:"phone"`
	IsFavorit bool    `json:"is_favorit"`
}

type FindAllPeopleResponse struct {
	Data []FindAllPeople `json:"data"`
}

type FindAllPeopleRequest struct {
	Fullname  string `form:"-"`
	Phone     string `form:"-"`
	Search    string `form:"search"`
	Limit     int    `form:"limit"`
	IsFavorit bool   `form:"is_favorit"`
}

type HistoryTransactionRequest struct {
	Fullname string `form:"-"`
	Phone    string `form:"-"`
	Search   string `form:"search"`
	Limit    int    `form:"limit"`
	Offset   int    `form:"page"`
	Type     string `form:"page"`
}

type HistoryTransaction struct {
	ID       int64   `json:"id"`
	Type     string  `json:"type"`
	Amount   int64   `json:"amount"`
	Fullname string  `json:"fullname"`
	ImageURL *string `json:"image_url"`
	Phone    *string `json:"phone"`
}

type HistoryTransactionResponse struct {
	TotalData int64                `json:"total_data"`
	TotalPage int64                `json:"total_page"`
	PrevPage  string               `json:"prev_page"`
	NextPage  string               `json:"next_page"`
	Data      []HistoryTransaction `json:"data"`
}

type TopupsRequest struct {
	Amount        int64  `json:"amount"`
	TotalAmount   int64  `json:"total_amount"`
	Tax           int64  `json:"tax"`
	XenditCode    string `json:"xendit_code"`
	PaymentMethod int    `json:"payment_method_id"`
}

type TopupsResponse struct {
	PaymentURL *string    `json:"payment_url"`
	ExpiresAt  *time.Time `json:"expires_at"`
}

type Recipient struct {
	ID       int64   `json:"id"`
	Fullname string  `json:"fullname"`
	ImageURL *string `json:"image_url"`
	Phone    *string `json:"phone"`
}

type TransferRequest struct {
	ReceiverID int64  `json:"receiver_id" binding:"required"`
	Amount     int64  `json:"amount" binding:"required"`
	Notes      string `json:"notes" binding:"max=255"`
}

type TransferResponse struct {
	ID          int64     `json:"id"`
	ReferenceNo string    `json:"reference_no"`
	Amount      int64     `json:"amount"`
	Notes       *string   `json:"notes"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	Sender      Recipient `json:"sender"`
	Receiver    Recipient `json:"receiver"`
}
