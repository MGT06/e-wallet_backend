package models

import (
	"encoding/json"
	"time"
)

type WebhookEvent struct {
	ID          int             `db:"id"`
	EventID     string          `db:"event_id"`
	EventStatus *string         `db:"event_status"`
	Payload     json.RawMessage `db:"payload"`
	CreatedAt   time.Time       `db:"created_at"`
}
