package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/MGT06/e-wallet_backend/internal/dto"
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
// @Param        	search  query     string  false  "search people"
// @Param        	limit   query     int  	  false  "limit get data"
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
// @Param        	body  body	dto.TopupsRequest  false  "amount and payment method"
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
