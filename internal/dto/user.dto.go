package dto

import "mime/multipart"

type UserHeaderResponse struct {
	Fullname string  `json:"name"`
	ImageURL *string `json:"image_url"`
}

type UserProfileResponse struct {
	Fullname string  `json:"name"`
	ImageURL *string `json:"image_url"`
	Phone    *string `db:"phone"`
	Email    string  `db:"email"`
}

type UserChangeProfileRequest struct {
	Fullname *string               `form:"name" example:"Alma'ruf Hidayat" binding:"omitempty,min=3,max=100"`
	ImageURL *multipart.FileHeader `form:"image" swaggerignore:"true" binding:"omitempty"`
	Phone    *string               `form:"phone" example:"081973779380" binding:"omitempty,min=12,max=14"`
	Email    *string               `form:"email" binding:"omitempty,email" example:"almarufhidayat99@gmail.com"`
}

type UserChangeProfileResponse struct {
	Fullname string  `json:"name"`
	ImageURL *string `json:"image_url"`
	Phone    *string `db:"phone"`
	Email    string  `db:"email"`
}
