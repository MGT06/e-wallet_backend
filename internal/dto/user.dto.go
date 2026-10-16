package dto

import "mime/multipart"

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
	Email    *string               `form:"email" binding:"omitempty,email,max=20" example:"almarufhidayat99@gmail.com"`
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