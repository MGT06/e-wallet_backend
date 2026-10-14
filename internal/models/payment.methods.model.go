package models

import (
	"time"
)

type PaymentMethod struct {
	ID         int32     `db:"id"`
	Name       string    `db:"name"`
	Type       string    `db:"type"`
	IconURL    *string   `db:"icon_url"`
	TaxPercent float64   `db:"tax_percent"`
	CreatedAt  time.Time `db:"created_at"`
}
