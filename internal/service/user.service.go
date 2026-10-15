package service

import (
	"context"

	"github.com/MGT06/e-wallet_backend/internal/dto"
	"github.com/MGT06/e-wallet_backend/internal/repository"
)

type UserService struct {
	ur *repository.UserRepo
}

func NewUserService(userRepo *repository.UserRepo) *UserService {
	return &UserService{
		ur: userRepo,
	}
}

func (us *UserService) GetUserHeaderUS(c context.Context, userID int64) (*dto.UserHeaderResponse, error) {
	user, err := us.ur.GetUserHeaderUR(c, userID)
	if err != nil {
		return nil, err
	}
	return &dto.UserHeaderResponse{
		Fullname: user.Fullname,
		ImageURL: user.ImageURL,
	}, nil
}
