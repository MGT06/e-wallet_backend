package service

import (
	"context"
	"fmt"
	"os"
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

func (ts *TransactionService) FindPeopleTS(ctx context.Context, userID int64, query dto.FindAllPeopleRequest) (*dto.FindAllPeopleResponse, error) {
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

	res, err := ts.tr.FindPeopleTR(ctx, userID, query)
	if err != nil {
		return nil, err
	}

	users := make([]dto.FindAllPeople, 0, len(*res))
	var data int64
	for _, u := range *res {
		data = u.SisaData
		var url string
		if u.ImageURL != nil {
			url = fmt.Sprintf("/%s/%v", os.Getenv("PROFILE_URL"), *u.ImageURL)
		} else {
			url = ""
		}
		users = append(users, dto.FindAllPeople{
			ID:       u.ID,
			Fullname: u.Fullname,
			ImageURL: &url,
			Phone:    u.Phone,
		})
	}

	if query.Limit == "" {
		data = 0
	}

	return &dto.FindAllPeopleResponse{
		TotalData: int64(len(users)),
		SisaData:  data,
		Data:      users,
	}, nil
}
