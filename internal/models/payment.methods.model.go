package models

import "time"

type PaymentMethodMDL struct {
	ID         int32     `db:"id"`
	Name       string    `db:"name"`
	Type       string    `db:"type"`
	XenditCode string    `db:"xendit_code"`
	IconURL    *string   `db:"icon_url"`
	TaxPercent float64   `db:"tax_percent"`
	CreatedAt  time.Time `db:"created_at"`
}
