package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MGT06/e-wallet_backend/internal/dto"
	"github.com/MGT06/e-wallet_backend/internal/mailer"
	msgerr "github.com/MGT06/e-wallet_backend/internal/message"
	"github.com/MGT06/e-wallet_backend/internal/repository"
	"github.com/MGT06/e-wallet_backend/internal/utils"
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

	if len(body.Email) == 0 || len(body.Name) == 0 || len(body.Password) < 8 {
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

	if len(body.Email) == 0 || len(body.Password) < 8 {
		return nil, msgerr.ErrEmailNotExist
	}

	user, err := as.ar.GetUserByEmail(c, body)
	if err != nil {
		return nil, err
	}

	_, err = as.conf.Compare(body.Password, user.Password)
	if err != nil {
		log.Printf("[AuthService.LoginAS] Compare %s", err)
		return nil, msgerr.ErrEmailNotExist
	}

	jwt := pkg.NewJWTClaims(user.ID, 60)

	token, err := jwt.GeneretToken()
	if err != nil {
		log.Printf("[AuthService.LoginAS] GeneretToke %s", err)
		return nil, err
	}

	key := fmt.Sprintf("ewallet:token:%d", user.ID)
	err = as.redis.Set(c, key, user.ID, 60*time.Minute).Err()
	if err != nil {
		log.Printf("[AuthService.LoginAS] Redis Set %s", err)
		return nil, err
	}

	return &dto.LoginResponse{
		Token:    token,
		IsHasPin: user.IsHasPin,
	}, nil
}

func (as *AuthService) LogoutAS(c context.Context, userID int64) (err error) {
	key := fmt.Sprintf("ewallet:token:%d", userID)

	err = as.redis.Del(c, key).Err()
	if err != nil {
		log.Printf("[AuthService.Logout] Gagal menghapus key Redis: %v", err)
		return err
	}

	return nil
}

type ForgotData struct {
	Email string
	OTP   string
}

func (as *AuthService) ForgotPasswordAS(c context.Context, email string) (string, error) {
	emailUser, err := as.ar.ForgotPasswordAR(c, email)
	if err != nil {
		return "", err
	}

	otp, er := utils.GenerateOTP()
	if er != nil {
		return "", err
	}

	m, fail := mailer.NewMailer()
	if fail != nil {
		return "", fail
	}

	html, errors := mailer.RenderTemplate("email.template.html", ForgotData{
		Email: emailUser,
		OTP:   otp,
	})

	if errors != nil {
		return "", errors
	}

	if err := m.SendMail(email, html); err != nil {
		return "", err
	}

	hash, e := as.conf.HashPassword(otp)
	if e != nil {
		return "", e
	}

	key := fmt.Sprintf("ewallet:otp:%s", email)
	ttl := 3 * time.Minute

	if err := utils.SetToRedis(c, as.redis, key, hash, ttl); err != nil {
		return "", err
	}

	return emailUser, nil
}

func (as *AuthService) CheckOTP(c context.Context, otp, email string) error {
	key := fmt.Sprintf("ewallet:otp:%s", email)

	realotp, err := utils.GetFromRedis[string](c, as.redis, key)
	if err != nil {
		return err
	}

	if _, err := as.conf.Compare(otp, realotp); err != nil {
		return err
	}

	return nil
}

func (as *AuthService) ResetPassword(c context.Context, body dto.ResetPassRequest) error {
	newPassword, err := as.conf.HashPassword(body.NewPassword)
	if err != nil {
		return err
	}

	if err := as.ar.ResetPassword(c, newPassword, body.Email); err != nil {
		return err
	}

	return nil
}
