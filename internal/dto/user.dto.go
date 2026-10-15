package dto

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
