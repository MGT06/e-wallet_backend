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

type UserContorller struct {
	us *service.UserService
}

func NewUserContorller(us *service.UserService) *UserContorller {
	return &UserContorller{
		us: us,
	}
}

// User info
//
// @Summary			User info
// @Description		Get user info
// @Tags			User
// @Produce			json
// @Security		Bearer
// @Router			/user/info	[get]
// @Success			200		{object}	dto.ResponseSuccess
// @Failure			400		{object}	dto.ResponseError
// @Failure			500		{object}	dto.ResponseError
func (uc *UserContorller) GetUserHeaderUC(c *gin.Context) {
	res, exist := c.Get("ID")
	if !exist {
		log.Printf("[UserContorller.GetUserHeaderUC] Error get token from redis")
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	userID, ok := res.(int64)
	if !ok {
		log.Printf("[UserContorller.GetUserHeaderUC] Error change to string")
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	user, err := uc.us.GetUserHeaderUS(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, msgerr.ErrUserNotExist) {
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
		Message: "get user success",
		Data:    user,
	})
}

// User profile
//
// @Summary			User profile
// @Description		Get user profile
// @Tags			User
// @Produce			json
// @Security		Bearer
// @Router			/user/profile	[get]
// @Success			200		{object}	dto.ResponseSuccess
// @Failure			400		{object}	dto.ResponseError
// @Failure			500		{object}	dto.ResponseError
func (uc *UserContorller) GetUserProfileUC(c *gin.Context) {
	res, exist := c.Get("ID")
	if !exist {
		log.Printf("[UserContorller.GetUserProfileUC] Error get token from redis")
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	userID, ok := res.(int64)
	if !ok {
		log.Printf("[UserContorller.GetUserProfileUC] Error change to string")
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	user, err := uc.us.GetUserProfileUS(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, msgerr.ErrUserNotExist) {
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
		Message: "get user profile success",
		Data:    user,
	})
}

// Change profile user
//
// @Summary         Change user profile
// @Description     Update user profile details including name, phone, email, and profile picture
// @Tags            User
// @Accept          multipart/form-data
// @Produce         json
// @Security        Bearer
// @Param           name    formData    string  false   "Full Name" example("Alma'ruf Hidayat")
// @Param           phone   formData    string  false   "Phone Number" example("081973779380")
// @Param           email   formData    string  false   "Email Address" example("almarufhidayat99@gmail.com")
// @Param           image   formData    file    false   "Profile Picture (.jpg, .jpeg, .png, max 2MB)"
// @Router          /user/profile   [patch]
// @Success         200     {object}    dto.ResponseSuccess{data=dto.UserChangeProfileResponse}
// @Failure         400     {object}    dto.ResponseError
// @Failure         500     {object}    dto.ResponseError
func (uc *UserContorller) ChangeProfileUC(c *gin.Context) {
	res, exist := c.Get("ID")
	if !exist {
		log.Printf("[UserContorller.ChangeProfileUC] Error get token from redis")
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has occurred",
		})
		return
	}

	userIDInt64, ok := res.(int64)
	if !ok {
		log.Printf("[UserContorller.ChangeProfileUC] Error casting ID to int64")
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has occurred",
		})
		return
	}
	userID := int32(userIDInt64)

	var body dto.UserChangeProfileRequest
	if err := c.ShouldBindWith(&body, binding.FormMultipart); err != nil {
		log.Printf("[UserContorller.ChangeProfileUC] Bind error: %s", err.Error())
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Success: false,
			Message: "invalid request format",
		})
		return
	}

	profile, err := uc.us.ChangeProfileUS(c.Request.Context(), userID, body)
	if err != nil {
		if errors.Is(err, msgerr.ErrExt) || errors.Is(err, msgerr.ErrImageToLarge) || errors.Is(err, msgerr.ErrUserNotExist) {
			c.JSON(http.StatusBadRequest, dto.ResponseError{
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
		Message: "change user profile success",
		Data:    profile,
	})
}
