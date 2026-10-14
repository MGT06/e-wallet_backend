package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/MGT06/e-wallet_backend/internal/dto"
	msgerr "github.com/MGT06/e-wallet_backend/internal/message"
	"github.com/MGT06/e-wallet_backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AuthContorller struct {
	as *service.AuthService
}

func NewAuthContorller(as *service.AuthService) *AuthContorller {
	return &AuthContorller{
		as: as,
	}
}

// Register
//
// @Summary			Create New Account
// @Description		Create new account for full access
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Router			/auth/register	[post]
// @Param			data	body 		dto.RegisterRequest true "Body data new user"
// @Success			200		{object}	dto.ResponseSuccess
// @Failure			400		{object}	dto.ResponseError
// @Failure			409		{object}	dto.ResponseError
// @Failure			500		{object}	dto.ResponseError
func (ac *AuthContorller) RegisterAH(c *gin.Context) {
	var body dto.RegisterRequest
	if err := c.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Printf("[AuthHandler.Register] ShouldBindWith: %v\n", err)
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Success: false,
			Message: "invalid request payload",
		})
		return
	}

	if err := ac.as.RegisterAS(c.Request.Context(), body); err != nil {
		if errors.Is(err, msgerr.ErrNoRowAffected) {
			c.JSON(http.StatusBadRequest, dto.ResponseError{
				Success: false,
				Message: "Invalid input data",
			})
			return
		}

		if errors.Is(err, msgerr.ErrEmailExist) {
			c.JSON(http.StatusConflict, dto.ResponseError{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.ResponseSuccess{
		Success: true,
		Message: "User registered successfully",
	})
}
