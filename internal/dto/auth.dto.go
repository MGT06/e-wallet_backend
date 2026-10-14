package dto

type RegisterRequest struct {
	Name     string `json:"name" example:"Alma'ruf Hidayat"`
	Email    string `json:"email" binding:"required,email" example:"almarufhidayat99@gmail.com"`
	Password string `json:"password" binding:"min=8" example:"Admin@1234"`
	Role     string `json:"-"`
}
