package dto

import "time"

type UpdateTopupStatusParams struct {
	ID        int64
	Status    string
	PaidAt    time.Time
	InvoiceID string
}

type XenditInvoiceRequest struct {
	ID             string   `json:"id"`
	ExternalID     string   `json:"external_id"`
	Status         string   `json:"status"`
	PaidAmount     *float64 `json:"paid_amount"`
	PaidAt         string   `json:"paid_at"`
	PaymentMethod  string   `json:"payment_method"`
	PaymentChannel string   `json:"payment_channel"`
	Currency       string   `json:"currency"`
}
