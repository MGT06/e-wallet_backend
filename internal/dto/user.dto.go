package dto

import (
	"mime/multipart"
)

type UserHeaderResponse struct {
	Fullname string  `json:"name"`
	ImageURL *string `json:"image_url"`
}

type UserProfileResponse struct {
	Fullname string  `json:"name"`
	ImageURL *string `json:"image_url"`
	Phone    *string `json:"phone"`
	Email    string  `json:"email"`
}

type UserChangeProfileRequest struct {
	Fullname *string               `form:"name" example:"Alma'ruf Hidayat" binding:"omitempty,min=3,max=100"`
	ImageURL *multipart.FileHeader `form:"image" swaggerignore:"true" binding:"omitempty"`
	Phone    *string               `form:"phone" example:"081973779380" binding:"omitempty,numeric,min=12,max=14"`
	Email    *string               `form:"email" binding:"omitempty,email,max=40" example:"almarufhidayat99@gmail.com"`
}

type UserChangeProfileResponse struct {
	Fullname string  `json:"name"`
	ImageURL *string `json:"image_url"`
	Phone    *string `json:"phone"`
	Email    string  `json:"email"`
}

type UserChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"min=8" example:"Admin@1234"`
	NewPassword string `json:"new_password" binding:"min=8" example:"Admin@1234"`
}

type SetPinRequest struct {
	Pin string `json:"pin" binding:"required,len=6,numeric" example:"123456"`
}

type PaymentMethodResponse struct {
	ID         int32   `json:"id"`
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	XenditCode string  `json:"xendit_code"`
	IconURL    *string `json:"icon_url"`
	TaxPercent float64 `json:"tax_percent"`
}

type DeletePhotoProfileResponse struct {
	ImageURL *string `json:"image_url"`
}

type AddFavoriteRequest struct {
	ID         int64 `json:"-"`
	FavoriteID int64 `json:"favorite_account_id"`
}

type ChangePinRequest struct {
	OldPin string `json:"old_pin" binding:"required,len=6,numeric" example:"123456"`
	NewPin string `json:"new_pin" binding:"required,len=6,numeric" example:"654321"`
}

type CheckPinRequest struct {
	OldPin string `json:"pin" binding:"required,len=6,numeric" example:"123456"`
}

type EwalletResponse struct {
	Balance int64 `json:"balance"`
	Income  int64 `json:"income"`
	Expanse int64 `json:"expense"`
}

type DeleteHistoryRequest struct {
	UserID    int64 `json:"-"`
	HistoryID int64 `json:"history_id"`
}

type TransactionSummaryRequest struct {
	Type string `json:"type"`
	Day  int64  `json:"filter"`
}

type TransactionSummaryResponse struct {
	Type   string `json:"type"`
	Amount int64  `json:"amount"`
	Date   string `json:"date"`
	Day    string `json:"day"`
	Month  string `json:"month"`
	Week   int    `json:"week"`
}
