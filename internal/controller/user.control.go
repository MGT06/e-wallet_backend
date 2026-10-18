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

	userID, ok := res.(int64)
	if !ok {
		log.Printf("[UserContorller.ChangeProfileUC] Error casting ID to int64")
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has occurred",
		})
		return
	}

	var body dto.UserChangeProfileRequest
	if err := c.ShouldBindWith(&body, binding.FormMultipart); err != nil {
		log.Printf("[UserContorller.ChangeProfileUC] Bind error: %s", err.Error())
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
			Message: "invalid request format",
		})
		return
	}

	profile, err := uc.us.ChangeProfileUS(c.Request.Context(), userID, body)
	if err != nil {
		if errors.Is(err, msgerr.ErrExt) || errors.Is(err, msgerr.ErrImageToLarge) || errors.Is(err, msgerr.ErrUserNotExist) || errors.Is(err, msgerr.ErrPhoneFormat) {
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

// Delete photo profile
//
// @Summary			Delete photo profile
// @Description		Delete foto profile
// @Tags			User
// @Produce			json
// @Security		Bearer
// @Router			/user/image	[delete]
// @Success			200		{object}	dto.ResponseSuccess
// @Failure			400		{object}	dto.ResponseError
// @Failure			500		{object}	dto.ResponseError
func (uc *UserContorller) DeletProfileUC(c *gin.Context) {
	res, exist := c.Get("ID")
	if !exist {
		log.Printf("[UserContorller.DeletProfileUC] Error get token from redis")
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	userID, ok := res.(int64)
	if !ok {
		log.Printf("[UserContorller.DeletProfileUC] Error change to string")
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	res, err := uc.us.DeletePhotoProfileUS(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, msgerr.ErrNoRowAffected) || errors.Is(err, msgerr.ErrImgNotFound) {
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
		Message: "deleted photo profile succcess",
		Data:    res,
	})
}

// Change password
//
// @Summary			Change password
// @Description		Change password user
// @Tags			User
// @Accept			json
// @Produce			json
// @Security		Bearer
// @Router			/user/password	[patch]
// @Param			body	body 		dto.UserChangePasswordRequest true "Body data change password"
// @Success			200		{object}	dto.ResponseSuccess
// @Failure			400		{object}	dto.ResponseError
// @Failure			500		{object}	dto.ResponseError
func (uc *UserContorller) ChangePasswordUC(c *gin.Context) {
	res, exist := c.Get("ID")
	if !exist {
		log.Printf("[UserContorller.DeletProfileUC] Error get token from redis")
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	userID, ok := res.(int64)
	if !ok {
		log.Printf("[UserContorller.DeletProfileUC] Error change to string")
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	var body dto.UserChangePasswordRequest
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

	err := uc.us.ChangePasswordUS(c.Request.Context(), userID, body)
	if err != nil {
		if errors.Is(err, msgerr.ErrNoRowAffected) || errors.Is(err, msgerr.ErrUserNotExist) || errors.Is(err, msgerr.ErrOldPassword) {
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
		Message: "change password succesfully",
	})
}

// Set Pin
//
// @Summary			Set Pin User
// @Tags			User
// @Accept			json
// @Produce			json
// @Security		Bearer
// @Router			/user/set-pin		[patch]
// @Param			body	body 		dto.SetPinRequest true	"pin user"
// @Success			204		{object}	dto.ResponseSuccess
// @Failure			400		{object}	dto.ResponseError
// @Failure			500		{object}	dto.ResponseError
func (uc *UserContorller) SetPinUC(c *gin.Context) {
	res, exist := c.Get("ID")
	if !exist {
		log.Printf("[UserContorller.DeletProfileUC] Error get token from redis")
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	userID, ok := res.(int64)
	if !ok {
		log.Printf("[UserContorller.DeletProfileUC] Error change to string")
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	var body dto.SetPinRequest
	if err := c.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Printf("[SetPinUC.ShouldBindWith] : %v", err)
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

	if err := uc.us.SetPinUS(c.Request.Context(), userID, body.Pin); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	c.JSON(http.StatusNoContent, dto.ResponseSuccess{
		Success: true,
		Message: "Set Pin Succesfully",
	})
}

// Add Favorite
//
// @Summary			Add Favorite
// @Tags			User
// @Accept			json
// @Produce			json
// @Security		Bearer
// @Router			/user/favorite		[post]
// @Param			body	body 		dto.AddFavoriteRequest true	"ID Account"
// @Success			204		{object}	dto.ResponseSuccess
// @Failure			400		{object}	dto.ResponseError
// @Failure			500		{object}	dto.ResponseError
func (uc *UserContorller) AddFavoriteUC(c *gin.Context) {
	res, exist := c.Get("ID")
	if !exist {
		log.Printf("[UserContorller.DeletProfileUC] Error get token from redis")
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	userID, ok := res.(int64)
	if !ok {
		log.Printf("[UserContorller.DeletProfileUC] Error change to string")
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	var body dto.AddFavoriteRequest
	if err := c.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Printf("[AddFavoriteUC.ShouldBindWith] : %v", err)
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

	if err := uc.us.AddFavoriteUS(c.Request.Context(), userID, body); err != nil {
		if errors.Is(err, msgerr.ErrUSerFavNotFound) || errors.Is(err, msgerr.ErrUserNotExist) {
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

	c.JSON(http.StatusNoContent, dto.ResponseSuccess{
		Success: true,
		Message: "add to favorite succesfully",
	})
}

// Change Pin
//
// @Summary			Change Pin
// @Description		Change pin user
// @Tags			User
// @Accept			json
// @Produce			json
// @Security		Bearer
// @Router			/user/change-pin	[patch]
// @Param			body	body 		dto.ChangePinRequest true "Body data change password"
// @Success			200		{object}	dto.ResponseSuccess
// @Failure			400		{object}	dto.ResponseError
// @Failure			500		{object}	dto.ResponseError
func (uc *UserContorller) ChangePinUC(c *gin.Context) {
	res, exist := c.Get("ID")
	if !exist {
		log.Printf("[UserContorller.ChangePinUC] Error get token from redis")
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	userID, ok := res.(int64)
	if !ok {
		log.Printf("[ChangePinUC] Error change to string")
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	var body dto.ChangePinRequest
	if err := c.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Printf("[ChangePinUC] ShouldBindWith: %v\n", err)
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

	err := uc.us.ChangePinUS(c.Request.Context(), userID, body)
	if err != nil {
		if errors.Is(err, msgerr.ErrOldPin) {
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
		Message: "change pin succesfully",
	})
}

// History
//
// @Summary      	History
// @Description  	Get all history transaction
// @Tags         	Transaction
// @Accept       	json
// @Produce      	json
// @Param        	search  query     string  false  "search people"
// @Param        	limit   query     int  	  false  "limit get data"
// @Param        	page   query     int  	  page  "limit get data"
// @Security     	Bearer
// @Router			/user/history		[get]
// @Success			200		{object}	dto.ResponseSuccess
// @Failure			400		{object}	dto.ResponseError
// @Failure			401		{object}	dto.ResponseError
// @Failure			500		{object}	dto.ResponseError
func (uc *UserContorller) HistoryUC(c *gin.Context) {
	res, exist := c.Get("ID")
	if !exist {
		log.Printf("[UserContorller.FindPeopleUC] Error get token from redis")
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	userID, ok := res.(int64)
	if !ok {
		log.Printf("[UserContorller.FindPeopleUC] Error change to string")
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	var query dto.HistoryTransactionRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	users, err := uc.us.HistoryTS(c.Request.Context(), userID, query)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		Success: true,
		Message: "success get hitory",
		Data:    users,
	})
}
