package dto

type UserHeaderResponse struct {
	Fullname string  `json:"name"`
	ImageURL *string `json:"image_url"`
}
