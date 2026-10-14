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
func (ac *AuthContorller) RegisterAC(c *gin.Context) {
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

// Login
//
// @Summary			Log In
// @Description		Log in for get full access
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Router			/auth/login	[post]
// @Param			data	body dto.LoginRequest true "Body data to access EventHub"
// @Success			200		{object}	dto.ResponseSuccess
// @Failure			400		{object}	dto.ResponseError
func (ac *AuthContorller) LoginAC(ctx *gin.Context) {
	var body dto.LoginRequest
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Printf("[AuthContorller.LoginAC] Bind error: %v\n", err.Error())
		ctx.JSON(http.StatusBadRequest, dto.ResponseError{
			Success: false,
			Message: "invalid request payload",
		})
		return
	}

	res, err := ac.as.LoginAS(ctx.Request.Context(), body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ResponseError{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.ResponseSuccess{
		Success: true,
		Message: "User login successfully",
		Data:    res,
	})
}
