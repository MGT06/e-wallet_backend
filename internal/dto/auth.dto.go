package dto

// Register
type RegisterRequest struct {
	Name     string `json:"name" example:"Alma'ruf Hidayat"`
	Email    string `json:"email" binding:"required,email" example:"almarufhidayat99@gmail.com"`
	Password string `json:"password" binding:"min=8" example:"Admin@1234"`
}

// Login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"almarufhidayat99@gmail.com"`
	Password string `json:"password" binding:"min=8" example:"Admin@1234"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Forgot
type ForgotRequest struct {
	Email string `json:"email" binding:"required,email"`
}

//Check OTP
type CheckOTPRequest struct {
	OTP string	`json:"otp" binding:"required,min=6,max=6"`
	Email string `json:"email" binding:"required,email"`
}