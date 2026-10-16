package dto

type ResponseSuccess struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"success"`
	Data    any    `json:"data,omitempty"`
}

type ResponseError struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"failed"`
	Data    any    `json:"data,omitempty"`
}
