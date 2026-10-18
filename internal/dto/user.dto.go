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
	Pin string `json:"pin" binding:"required,min=6,max=6"`
}

type PaymentMethodResponse struct {
	ID         int32   `db:"id"`
	Name       string  `db:"name"`
	Type       string  `db:"type"`
	XenditCode string  `db:"xendit_code"`
	IconURL    *string `db:"icon_url"`
	TaxPercent float64 `db:"tax_percent"`
}

type DeletePhotoProfileResponse struct {
	ImageURL *string `json:"image_url"`
}

type AddFavoriteRequest struct {
	ID         int64 `json:"-"`
	FavoriteID int64 `json:"favorite_account_id"`
}
