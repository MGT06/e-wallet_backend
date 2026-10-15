package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MGT06/e-wallet_backend/internal/dto"
	msgerr "github.com/MGT06/e-wallet_backend/internal/message"
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

func (us *UserService) GetUserProfileUS(c context.Context, userID int64) (*dto.UserProfileResponse, error) {
	user, err := us.ur.GetUserHeaderUR(c, userID)
	if err != nil {
		return nil, err
	}
	return &dto.UserProfileResponse{
		Fullname: user.Fullname,
		ImageURL: user.ImageURL,
		Phone:    user.Phone,
		Email:    user.Email,
	}, nil
}

func (us *UserService) ChangeProfileUS(ctx context.Context, userID int32, body dto.UserChangeProfileRequest) (*dto.UserChangeProfileResponse, error) {

	var filename string
	dirPath := filepath.Join("public", "photo_profiles")

	if body.ImageURL != nil {
		file := body.ImageURL

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
			return nil, msgerr.ErrExt
		}

		if file.Size > 2*1024*1024 {
			return nil, msgerr.ErrImageToLarge
		}

		oldUser, err := us.ur.GetUserHeaderUR(ctx, int64(userID))
		if err == nil && oldUser.ImageURL != nil && *oldUser.ImageURL != "" {
			oldFilePath := filepath.Join("public", "photo_profiles", filepath.Base(*oldUser.ImageURL))
			if err := os.Remove(oldFilePath); err != nil && !os.IsNotExist(err) {
				log.Printf("[ChangeProfileUS.RemoveOldFile] : %s", err.Error())
			}
		}

		timestamp := time.Now().Unix()
		filename = fmt.Sprintf("%d_%d%s", userID, timestamp, ext)

		if err := os.MkdirAll(dirPath, 0755); err != nil {
			log.Printf("[ChangeProfileUS.path.Join] : %s", err.Error())
			return nil, msgerr.ErrWriteFile
		}

		targetFilePath := filepath.Join(dirPath, filename)

		src, err := file.Open()
		if err != nil {
			return nil, msgerr.ErrOpenFile
		}
		defer src.Close()

		data, err := io.ReadAll(src)
		if err != nil {
			log.Printf("[ChangeProfileUS.ReadAll] : %s", err.Error())
			return nil, msgerr.ErrReadFile
		}

		if err := os.WriteFile(targetFilePath, data, 0644); err != nil {
			log.Printf("[ChangeProfileUS.WriteFile] : %s", err.Error())
			return nil, msgerr.ErrWriteFile
		}
	}

	profile, err := us.ur.ChangeProfileUR(ctx, userID, filename, body)
	if err != nil {
		if filename != "" {
			_ = os.Remove(filepath.Join(dirPath, filename))
		}
		return nil, err
	}

	return &dto.UserChangeProfileResponse{
		Fullname: profile.Fullname,
		ImageURL: profile.ImageURL,
		Phone:    profile.Phone,
		Email:    profile.Email,
	}, nil
}
