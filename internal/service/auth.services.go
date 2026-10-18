package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
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

func (as *AuthService) ForgotPasswordAS(c context.Context, email string) error {
	keyNX := fmt.Sprintf("ewallet:otp:cooldown:%s", email)
	ok, err := as.redis.SetNX(c, keyNX, 1, 60*time.Second).Result()
	if err != nil {
		return err
	}

	if !ok {
		return msgerr.ErrTooManyRequests
	}

	emailUser, err := as.ar.ForgotPasswordAR(c, email)
	if err != nil {
		if errors.Is(err, msgerr.ErrEmailNotExist) {
			return nil
		}
		return err
	}

	go as.proccessOTP(emailUser)

	return nil
}

func (as *AuthService) proccessOTP(email string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	otp, err := utils.GenerateOTP()
	if err != nil {
		log.Println(err)
		return
	}

	hash, er := as.conf.HashPassword(otp)
	if er != nil {
		log.Println(er)
		return
	}

	keyotp := fmt.Sprintf("ewallet:otp:%s", email)
	ttl := 3 * time.Minute

	if err := utils.SetToRedis(ctx, as.redis, keyotp, hash, ttl); err != nil {
		log.Println(err)
		return
	}

	keyattempt := fmt.Sprintf("ewallet:otp:attempts:%s", email)
	utils.DelFromRedis(ctx, as.redis, keyattempt)

	html, e := mailer.RenderTemplate(ForgotData{
		Email: email,
		OTP:   otp,
	})
	if e != nil {
		utils.DelFromRedis(ctx, as.redis, keyotp)
		log.Println(e)
		return
	}

	m, fail := mailer.NewMailer()
	if fail != nil {
		utils.DelFromRedis(ctx, as.redis, keyotp)
		log.Println(fail)
		return
	}

	if err := m.SendMail(email, html); err != nil {
		utils.DelFromRedis(ctx, as.redis, keyotp)
		log.Println(err)
		return
	}
}

func (as *AuthService) CheckOTP(c context.Context, otp, email string) (string, error) {

	keyattempt := fmt.Sprintf("ewallet:otp:attempts:%s", email)
	pipe := as.redis.TxPipeline()

	incr := pipe.Incr(c, keyattempt)
	pipe.ExpireNX(c, keyattempt, 3*time.Minute)

	if _, err := pipe.Exec(c); err != nil {
		return "", err
	}

	keyotp := fmt.Sprintf("ewallet:otp:%s", email)

	if incr.Val() > 3 {
		utils.DelFromRedis(c, as.redis, keyotp)
		return "", msgerr.ErrTooManyAttempts
	}

	hashotp, er := utils.GetFromRedis[string](c, as.redis, keyotp)
	if er != nil {
		if errors.Is(er, redis.Nil) {
			return "", msgerr.ErrInvalidOTP
		}
		return "", er
	}

	if _, err := as.conf.Compare(otp, hashotp); err != nil {
		return "", msgerr.ErrInvalidOTP
	}

	deleted, e := as.redis.Del(c, keyotp).Result()
	if e != nil {
		return "", e
	}
	if deleted == 0 {
		return "", msgerr.ErrInvalidOTP
	}

	b := make([]byte, 60)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	token := base64.RawURLEncoding.EncodeToString(b)
	keytoken := fmt.Sprintf("ewallet:reset:%s", token)
	if err := utils.SetToRedis(c, as.redis, keytoken, email, 10*time.Minute); err != nil {
		return "", err
	}

	return token, nil
}

func (as *AuthService) ResetPassword(c context.Context, body dto.ResetPassRequest) error {
	keytoken := fmt.Sprintf("ewallet:reset:%s", body.Token)
	email, err := utils.GetFromRedis[string](c, as.redis, keytoken)
	if err != nil {
		return err
	}

	deleted, er := as.redis.Del(c, keytoken).Result()
	if er != nil {
		return er
	}
	if deleted == 0 {
		return msgerr.ErrInvalidToken
	}

	newPassword, e := as.conf.HashPassword(body.NewPassword)
	if e != nil {
		return e
	}

	if err := as.ar.ResetPassword(c, newPassword, email); err != nil {
		return err
	}

	return nil
}
