package models

import (
	"time"
)

type Favorite struct {
	AccountID         int64     `db:"account_id"`
	FavoriteAccountID int64     `db:"favorite_account_id"`
	CreatedAt         time.Time `db:"created_at"`
}
