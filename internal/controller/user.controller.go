package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/MGT06/e-wallet_backend/internal/dto"
	msgerr "github.com/MGT06/e-wallet_backend/internal/message"
	"github.com/MGT06/e-wallet_backend/internal/service"
	"github.com/gin-gonic/gin"
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
	res, ok := c.Get("ID")
	if !ok {
		log.Printf("[UserContorller.GetUserHeaderUC] Error get token from redis")
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Success: false,
			Message: "please login first",
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
			Message: "please login first",
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
