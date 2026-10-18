package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/MGT06/e-wallet_backend/internal/dto"
	msgerr "github.com/MGT06/e-wallet_backend/internal/message"
	"github.com/MGT06/e-wallet_backend/internal/service"
	"github.com/MGT06/e-wallet_backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
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
// @Param			body	body 		dto.RegisterRequest true "Body data new user"
// @Success			200		{object}	dto.ResponseSuccess
// @Failure			400		{object}	dto.ResponseError
// @Failure			409		{object}	dto.ResponseError
// @Failure			500		{object}	dto.ResponseError
func (ac *AuthContorller) RegisterAC(c *gin.Context) {
	var body dto.RegisterRequest
	if err := c.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Printf("[AuthHandler.Register] ShouldBindWith: %v\n", err)
		found, isErr := errors.AsType[validator.ValidationErrors](err)
		if isErr {
			err := utils.ValidationError(found)
			c.JSON(http.StatusBadRequest, dto.ResponseError{
				Success: false,
				Message: err,
			})
			return
		}

		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Success: false,
			Message: "invalid request payload",
		})
		return
	}

	if err := ac.as.RegisterAS(c.Request.Context(), body); err != nil {
		log.Println(err)
		if errors.Is(err, msgerr.ErrNoRowAffected) {
			c.JSON(http.StatusBadRequest, dto.ResponseError{
				Success: false,
				Message: "Invalid input data",
			})
			return
		}

		if errors.Is(err, msgerr.ErrBodyEmpty) {
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

		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
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
// @Param			body	body 		dto.LoginRequest true "Body data to access ewalet"
// @Success			200		{object}	dto.ResponseSuccess
// @Failure			400		{object}	dto.ResponseError
func (ac *AuthContorller) LoginAC(c *gin.Context) {
	var body dto.LoginRequest
	if err := c.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Printf("[AuthContorller.LoginAC] Bind error: %v\n", err.Error())
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Success: false,
			Message: msgerr.ErrEmailNotExist.Error(),
		})
		return
	}

	res, err := ac.as.LoginAS(c.Request.Context(), body)
	if err != nil {
		log.Printf("[AuthContorller.LoginAS] Bind error: %v\n", err.Error())
		if errors.Is(err, msgerr.ErrEmailNotExist) {
			c.JSON(http.StatusBadRequest, dto.ResponseError{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		Success: true,
		Message: "User login successfully",
		Data:    res,
	})
}

// Forgot Password
//
// @Summary			Validation email for reset password
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Router			/auth/forgot	[post]
// @Param			body	body 		dto.ForgotRequest 	true 	"email user"
// @Success			200		{object}	dto.ResponseSuccess
// @Failure			400		{object}	dto.ResponseError
// @Failure			401		{object}	dto.ResponseError
// @Failure			500		{object}	dto.ResponseError
func (ac *AuthContorller) ForgotPasswordAC(c *gin.Context) {
	var body dto.ForgotRequest
	if err := c.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Success: false,
			Message: "invalid request payload",
		})
		return
	}

	emailUser, err := ac.as.ForgotPasswordAS(c.Request.Context(), body.Email)
	if err != nil {
		log.Println(err)
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, dto.ResponseError{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has occurred",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		Success: true,
		Data:    emailUser,
		Message: "email exists, please reset your password",
	})
}

// Logout
//
// @Summary			Log out
// @Description		Log out for delete access user
// @Tags			Auth
// @Produce			json
// @Security		Bearer
// @Router			/auth/logout	[delete]
// @Success			200		{object}	dto.ResponseSuccess
// @Failure			400		{object}	dto.ResponseError
func (ac *AuthContorller) LogoutAC(c *gin.Context) {
	res, exist := c.Get("ID")
	if !exist {
		log.Printf("[LogoutAC.Get ID] Error get token from redis")
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	userID, ok := res.(int64)
	if !ok {
		log.Printf("[LogoutAC.res] Error change to string")
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	err := ac.as.LogoutAS(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		Success: true,
		Message: "successfully logout ",
	})
}

// Check OTP
//
// @Summary			Check OTP user
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Router			/auth/check-otp		[post]
// @Param			body	body 		dto.CheckOTPRequest 	true 	"user"
// @Success			200		{object}	dto.ResponseSuccess
// @Success			400		{object}	dto.ResponseSuccess
// @Success			401		{object}	dto.ResponseSuccess
// @Failure			500		{object}	dto.ResponseError
func (ac *AuthContorller) CheckOTP(c *gin.Context) {
	var body dto.CheckOTPRequest
	if err := c.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Success: false,
			Message: "invalid request payload",
		})
		return
	}

	if err := ac.as.CheckOTP(c.Request.Context(), body.OTP, body.Email); err != nil {
		if errors.Is(err, msgerr.ErrInvalidInput) {
			c.JSON(http.StatusUnauthorized, dto.ResponseError{
				Success: false,
				Message: "OTP does not match",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		Success: true,
		Message: "OTP is match, please reset your password",
	})
}

// Reset Password
//
// @Summary			Reset Password User
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Router			/auth/reset-password	[post]
// @Param			body	body 		dto.ResetPassRequest 	true 	"email and new password user"
// @Success			204		{object}	dto.ResponseSuccess
// @Success			400		{object}	dto.ResponseSuccess
// @Failure			500		{object}	dto.ResponseError
func (ac *AuthContorller) ResetPasswordAC(c *gin.Context) {
	var body dto.ResetPassRequest
	if err := c.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err)
		found, isErr := errors.AsType[validator.ValidationErrors](err)
		if isErr {
			err := utils.ValidationError(found)
			c.JSON(http.StatusBadRequest, dto.ResponseError{
				Success: false,
				Message: err,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	if err := ac.as.ResetPassword(c.Request.Context(), body); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	c.JSON(http.StatusNoContent, dto.ResponseSuccess{
		Success: true,
		Message: "Succes Reset Password",
	})
}
