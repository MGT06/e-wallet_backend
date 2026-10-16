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

func (ts *TransactionService) FindPeopleTS(ctx context.Context, tserID int64, query dto.FindAllPeopleRequest) (*[]dto.FindAllPeopleResponse, error) {
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

	res, err := ts.tr.FindPeopleUR(ctx, tserID, query)
	if err != nil {
		return nil, err
	}

	tsers := make([]dto.FindAllPeopleResponse, 0, len(*res))
	for _, u := range *res {
		tsers = append(tsers, dto.FindAllPeopleResponse{
			ID:       u.ID,
			Fullname: u.Fullname,
			ImageURL: u.ImageURL,
			Phone:    u.Phone,
		})
	}

	return &tsers, nil
}
