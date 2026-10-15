package service

import (
	"context"
	"fmt"
	"log"
	"time"

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

	if user.DeletedAt != nil {
		return nil, msgerr.ErrAccounDeleted
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

	key := fmt.Sprintf("ewallet:token:%d", user.ID)
	err = as.redis.Set(c, key, true, 3*time.Minute).Err()
	if err != nil {
		log.Printf("[AuthService.LoginAS] Redis Set %s", err)
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
	}, nil
}

func (as *AuthService) LogoutAS(c context.Context, userID int64) (err error) {

	key := fmt.Sprintf("ewallet:token:%d", userID)

	err = as.CheckSessionAS(c, key)
	if err != nil {
		return err
	}

	err = as.redis.Del(c, key).Err()
	if err != nil {
		log.Printf("[AuthService.Logout] Gagal menghapus key Redis: %v", err)
		return err
	}

	return nil
}

func (as *AuthService) CheckSessionAS(c context.Context, key string) (err error) {
	key, err = as.redis.Get(c, key).Result()
	if err != nil {
		log.Printf("[AuthService.CheckSession] Redis Get %s", err)
		return msgerr.ErrTokenExpired
	}

	if key == "0" {
		return msgerr.ErrSesionLogout
	}

	return nil
}
