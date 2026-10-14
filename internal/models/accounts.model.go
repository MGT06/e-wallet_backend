package models

import "time"

type Account struct {
	ID        int64      `db:"id"`
	Fullname  string     `db:"fullname"`
	ImageURL  *string    `db:"image_url"`
	Phone     *string    `db:"phone"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	Pin       *string    `db:"pin"`
	Balance   int64      `db:"balance"`
	DeletedAt *time.Time `db:"deleted_at"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
}
