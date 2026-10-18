package service

import (
	"context"
	"unicode"

	"github.com/MGT06/e-wallet_backend/internal/dto"
	"github.com/MGT06/e-wallet_backend/internal/repository"
)

type TransactionService struct {
	tr *repository.TransactionRepo
}

func NewTransactionService(tr *repository.TransactionRepo) *TransactionService {
	return &TransactionService{
		tr: tr,
	}
}

func (ts *TransactionService) FindPeopleTS(c context.Context, userID int64, query dto.FindAllPeopleRequest) (*dto.FindAllPeopleResponse, error) {
	if query.Limit <= 0 {
		query.Limit = 10
	}
	var isByName bool
	if query.Search == "" {
		isByName = false
	}

	for _, r := range query.Search {
		if unicode.IsNumber(r) {
			isByName = false
		}
		if unicode.IsLetter(r) {
			isByName = true
		} else {
			isByName = false
		}
	}

	if isByName {
		query.Fullname = query.Search
	} else {
		query.Phone = query.Search
	}

	res, err := ts.tr.FindPeopleTR(c, userID, query)
	if err != nil {
		return nil, err
	}

	users := make([]dto.FindAllPeople, 0, len(*res))
	for _, u := range *res {
		users = append(users, dto.FindAllPeople{
			ID:        u.ID,
			Fullname:  u.Fullname,
			ImageURL:  u.ImageURL,
			Phone:     u.Phone,
			IsFavorit: u.IsFavorit,
		})
	}

	return &dto.FindAllPeopleResponse{
		Data: users,
	}, nil
}

func (ts *TransactionService) PaymentMethodeTS(c context.Context) (*[]dto.PaymentMethodResponse, error) {
	res, err := ts.tr.PaymentMethodeTR(c)
	if err != nil {
		return nil, err
	}

	payments := make([]dto.PaymentMethodResponse, 0, len(*res))
	for _, p := range *res {
		payment := dto.PaymentMethodResponse{
			ID:         p.ID,
			Name:       p.Name,
			Type:       p.Type,
			XenditCode: p.XenditCode,
			IconURL:    p.IconURL,
			TaxPercent: p.TaxPercent,
		}
		payments = append(payments, payment)
	}
	return &payments, nil
}
