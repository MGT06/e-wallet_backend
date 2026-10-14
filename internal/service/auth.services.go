package service

import (
	"context"
	"fmt"
	"log"

	"github.com/MGT06/e-wallet_backend/internal/dto"
	msgerr "github.com/MGT06/e-wallet_backend/internal/message"
	"github.com/MGT06/e-wallet_backend/internal/repository"
	"github.com/MGT06/e-wallet_backend/pkg"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	conf  *pkg.ConfigHash
	redis *redis.Client
	ar    *repository.AuthRepo
}

func NewAuthService(AuthRepo *repository.AuthRepo, redis *redis.Client) *AuthService {
	return &AuthService{
		ar:    AuthRepo,
		conf:  pkg.NewConfigHash(),
		redis: redis,
	}
}

func (as *AuthService) RegisterAS(c context.Context, body dto.RegisterRequest) (err error) {

	if len(body.Email) == 0 || len(body.Name) == 0 || len(body.Password) <= 8 {
		return msgerr.ErrBodyEmpty
	}

	hashPassword, err := as.conf.HashPassword(body.Password)
	if err != nil {
		log.Printf("[AuthService.RegisterAS] HashPassword %s", err)
		return err
	}

	body.Password = hashPassword

	err = as.ar.RegisterAR(c, body)
	if err != nil {
		return err
	}
	return nil
}

func (as *AuthService) LoginAS(c context.Context, body dto.LoginRequest) (*dto.LoginResponse, error) {

	if len(body.Email) == 0 || len(body.Password) <= 8 {
		return nil, msgerr.ErrBodyEmpty
	}

	user, err := as.ar.GetUserByEmail(c, body)
	if err != nil {
		return nil, err
	}

	_, err = as.conf.Compare(body.Password, user.Password)
	if err != nil {
		log.Printf("[AuthService.LoginAS] Compare %s", err)
		return nil, err
	}

	jwt := pkg.NewJWTClaims(user.ID, 3)

	token, err := jwt.GeneretToken()
	if err != nil {
		log.Printf("[AuthService.LoginAS] GeneretToke %s", err)
		return nil, err
	}

	key := fmt.Sprintf("ewalet:token:%d", user.ID)
	as.redis.Set(c, key, key, 3)

	return &dto.LoginResponse{
		Token: token,
	}, nil
}
