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

type TransactionContorller struct {
	ts *service.TransactionService
}

func NewTransactionContorller(ts *service.TransactionService) *TransactionContorller {
	return &TransactionContorller{
		ts: ts,
	}
}

// Find People
//
// @Summary      	Find people
// @Description  	Get all users active
// @Tags         	Transaction
// @Accept       	json
// @Produce      	json
// @Param        	search  	 query     string     false  "search people"
// @Param        	limit   	 query     int  	  false  "limit get data"
// @Param        	is_favorit   query     bool  	  false  "is_Favorite get data"
// @Security     	Bearer
// @Router			/transaction/people		[get]
// @Success			200		{object}	dto.ResponseSuccess
// @Failure			400		{object}	dto.ResponseError
// @Failure			401		{object}	dto.ResponseError
// @Failure			500		{object}	dto.ResponseError
func (tc *TransactionContorller) FindPeopleUC(c *gin.Context) {
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

	var query dto.FindAllPeopleRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Success: false,
			Message: "invalid query parameters",
		})
		return
	}

	users, err := tc.ts.FindPeopleTS(c.Request.Context(), userID, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		Success: true,
		Message: "success get people list",
		Data:    users,
	})
}

// Get payments methode
//
// @Summary			Get payments methode
// @Description		Get payments methode
// @Tags			Transaction
// @Produce			json
// @Security		Bearer
// @Router			/transaction/topup-methode	[get]
// @Success			200		{object}	dto.ResponseSuccess
// @Failure			400		{object}	dto.ResponseError
// @Failure			500		{object}	dto.ResponseError
func (tc *TransactionContorller) PaymentMethodeTC(c *gin.Context) {
	methode, err := tc.ts.PaymentMethodeTS(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		Success: true,
		Message: "get payment methode success",
		Data:    methode,
	})
}

// Topup
//
// @Summary      	user topup
// @Tags         	Transaction
// @Accept       	json
// @Produce      	json
// @Param        	body  body	dto.TopupsRequest  true  "amount and payment method"
// @Security     	Bearer
// @Router			/transaction/topup	[post]
// @Success			200		{object}	dto.ResponseSuccess
// @Failure			400		{object}	dto.ResponseError
// @Failure			401		{object}	dto.ResponseError
// @Failure			500		{object}	dto.ResponseError
func (tc *TransactionContorller) TopupsUC(c *gin.Context) {
	res, exist := c.Get("ID")
	if !exist {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	var body dto.TopupsRequest
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

	resp, err := tc.ts.TopupsTS(c.Request.Context(), res.(int64), body)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		Success: true,
		Message: "Topups Pending",
		Data:    resp,
	})
}

// Transfer
//
// @Summary      	user transfer
// @Tags         	Transaction
// @Accept       	json
// @Produce      	json
// @Param        	body  body	dto.TransferRequest  false  "amount and payment method"
// @Security     	Bearer
// @Router			/transaction/transfer	[post]
// @Success			201		{object}	dto.ResponseSuccess
// @Failure			400		{object}	dto.ResponseError
// @Failure			404		{object}	dto.ResponseError
// @Failure			422		{object}	dto.ResponseError
// @Failure			500		{object}	dto.ResponseError
func (tc *TransactionContorller) TransferUC(ctx *gin.Context) {
	res, exist := ctx.Get("ID")
	if !exist {
		ctx.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	var req dto.TransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		found, isErr := errors.AsType[validator.ValidationErrors](err)
		if isErr {
			err := utils.ValidationError(found)
			ctx.JSON(http.StatusBadRequest, dto.ResponseError{
				Success: false,
				Message: err,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.ResponseError{
			Success: false,
			Message: "a system error has accurred",
		})
		return
	}

	err := tc.ts.Transfer(ctx.Request.Context(), res.(int64), req)
	if err != nil {
		log.Println(err)
		switch {
		case errors.Is(err, msgerr.ErrSelfTransfer):
			ctx.JSON(http.StatusBadRequest, dto.ResponseError{
				Success: false,
				Message: "You cannot transfer to your own account",
			})
		case errors.Is(err, msgerr.ErrTooManyAttempts):
			ctx.JSON(http.StatusTooManyRequests, dto.ResponseSuccess{
				Success: true,
				Message: err.Error(),
			})
		case errors.Is(err, msgerr.ErrInvalidInput):
			ctx.JSON(http.StatusBadRequest, dto.ResponseError{
				Success: false,
				Message: "wrong pin",
			})
		case errors.Is(err, msgerr.ErrReceiverNotFound):
			ctx.JSON(http.StatusNotFound, dto.ResponseError{
				Success: false,
				Message: "Receiver not found",
			})
		case errors.Is(err, msgerr.ErrInsufficientBalance):
			ctx.JSON(http.StatusUnprocessableEntity, dto.ResponseError{
				Success: false,
				Message: "Insufficient balance",
			})
		default:
			ctx.JSON(http.StatusInternalServerError, dto.ResponseError{
				Success: false,
				Message: "a system error has accurred",
			})
		}
		return
	}

	ctx.JSON(http.StatusCreated, dto.ResponseSuccess{
		Success: true,
		Message: "Transfer Success",
	})
}

// Get Receiver
//
// @Summary      	Get  receiver
// @Description  	Get  receiver
// @Tags         	User
// @Produce			json
// @Security     	Bearer
// @Param        	user_id   query     int64  	  page  "Get  receiver"
// @Router			/transaction		[get]
// @Success			200		{object}	dto.ResponseSuccess
// @Failure			400		{object}	dto.ResponseError
// @Failure			401		{object}	dto.ResponseError
// @Failure			500		{object}	dto.ResponseError
func (tc *TransactionContorller) GetReceiverTC(c *gin.Context) {
	var query dto.ReceiverRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseError{
			Success: false,
			Message: "invalid query parameters",
		})
		return
	}

	log.Printf("GetReceiverTC.ID TARGET : %d", query.UserID)

	users, err := tc.ts.GetReceiverTS(c.Request.Context(), query.UserID)
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
		Message: "success get people list",
		Data:    users,
	})
}
