package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/MGT06/e-wallet_backend/internal/dto"
	msgerr "github.com/MGT06/e-wallet_backend/internal/message"
	"github.com/MGT06/e-wallet_backend/internal/repository"
	"github.com/MGT06/e-wallet_backend/pkg"
)

type UserService struct {
	ur   *repository.UserRepo
	conf *pkg.ConfigHash
}

func NewUserService(userRepo *repository.UserRepo) *UserService {
	return &UserService{
		ur:   userRepo,
		conf: pkg.NewConfigHash(),
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
	user, err := us.ur.GetUserProfileUR(c, userID)
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

func (us *UserService) ChangeProfileUS(c context.Context, userID int64, body dto.UserChangeProfileRequest) (*dto.UserChangeProfileResponse, error) {
	if body.Phone != nil {
		if strings.HasPrefix(*body.Phone, "62") {
			*body.Phone = "0" + (*body.Phone)[2:]
		}

		if !strings.HasPrefix(*body.Phone, "0") {
			return nil, msgerr.ErrPhoneFormat
		}
	}

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

		oldUser, err := us.ur.GetUserHeaderUR(c, userID)
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

	profile, err := us.ur.ChangeProfileUR(c, userID, filename, body)
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

func (us *UserService) DeletePhotoProfileUS(c context.Context, userID int64) (*dto.DeletePhotoProfileResponse, error) {
	oldUser, err := us.ur.GetUserHeaderUR(c, userID)

	if *oldUser.ImageURL == "" {
		return nil, msgerr.ErrImgNotFound
	}

	oldFilePath := filepath.Join("public", "photo_profiles", filepath.Base(*oldUser.ImageURL))
	if err := os.Remove(oldFilePath); err != nil && !os.IsNotExist(err) {
		log.Printf("[DeletePhotoProfileUS.RemoveOldFile] : %s", err.Error())
	}

	err = us.ur.DeletePhotoProfileUR(c, userID)
	if err != nil {
		return nil, err
	}
	return &dto.DeletePhotoProfileResponse{
		ImageURL: new(string),
	}, nil
}

func (us *UserService) ChangePasswordUS(c context.Context, userID int64, body dto.UserChangePasswordRequest) error {
	user, err := us.ur.GetUserUR(c, userID)

	_, err = us.conf.Compare(body.OldPassword, user.Password)
	if err != nil {
		log.Printf("[AuthService.LoginAS] Compare %s", err)
		return msgerr.ErrOldPassword
	}

	newPassword, err := us.conf.HashPassword(body.NewPassword)
	if err != nil {
		log.Printf("[AuthService.RegisterAS] HashPassword %s", err)
		return err
	}

	err = us.ur.ChangePasswordUR(c, userID, newPassword)
	if err != nil {
		return nil
	}
	return nil
}

func (us *UserService) SetPinUS(c context.Context, idUser int64, pin string) error {
	pinhash, err := us.conf.HashPassword(pin)
	if err != nil {
		return err
	}

	if err := us.ur.SetPinUR(c, pinhash, idUser); err != nil {
		return err
	}

	return nil
}

func (us *UserService) ToggleFavoriteUS(c context.Context, userID int64, body dto.AddFavoriteRequest) (string, error) {
	body.ID = userID
	res, err := us.ur.ToggleFavoriteUR(c, body)
	if err != nil {
		return "", err
	}

	return res, nil
}

func (us *UserService) ChangePinUS(c context.Context, userID int64, body dto.ChangePinRequest) error {

	user, err := us.ur.GetUserUR(c, userID)

	_, err = us.conf.Compare(body.OldPin, *user.Pin)
	if err != nil {
		log.Printf("[AddFavoriteUS.Compare] : %s", err)
		return msgerr.ErrOldPin
	}

	newPin, err := us.conf.HashPassword(body.NewPin)

	if err != nil {
		return err
	}

	err = us.ur.ChangePinUR(c, userID, newPin)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) CheckPinUS(c context.Context, userID int64, body dto.CheckPinRequest) error {
	user, err := us.ur.GetUserUR(c, userID)

	_, err = us.conf.Compare(body.OldPin, *user.Pin)
	if err != nil {
		log.Printf("[AddFavoriteUS.Compare] : %s", err)
		return msgerr.ErrOldPin
	}
	return nil
}

func (us *UserService) HistoryTS(c context.Context, userID int64, query dto.HistoryTransactionRequest) (*dto.HistoryTransactionResponse, error) {
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

	page := query.Offset

	offset := (query.Offset - 1) * query.Limit

	query.Offset = offset

	res, err := us.ur.HistoryUR(c, userID, query)
	if err != nil {
		return nil, err
	}

	historys := make([]dto.HistoryTransaction, 0, len(*res))
	var totalData float64
	for _, u := range *res {
		totalData = u.TotalData
		history := dto.HistoryTransaction{
			Fullname: u.Fullname,
			ImageURL: u.ImageURL,
			Phone:    u.Phone,
			Amount:   u.Amount,
			ID:       u.ID,
			Type:     u.Type,
		}
		historys = append(historys, history)
	}

	totalPage := math.Ceil(totalData / float64(query.Limit))

	pathUrl := "/user/history"

	prev := url.Values{}
	next := url.Values{}
	if isByName {
		prev.Add("search", query.Fullname)
		next.Add("search", query.Fullname)
	}

	if !isByName && query.Phone != "" {
		prev.Add("search", query.Phone)
		next.Add("search", query.Phone)
	}

	prev.Add("limit", strconv.Itoa(query.Limit))
	next.Add("limit", strconv.Itoa(query.Limit))
	prev.Add("offset", strconv.Itoa(page-1))
	next.Add("offset", strconv.Itoa(page+1))

	prevPage := fmt.Sprintf("%s?%s", pathUrl, prev.Encode())
	nextPage := fmt.Sprintf("%s?%s", pathUrl, next.Encode())

	if page <= 1 {
		prevPage = ""
	}

	if page >= int(totalPage) {
		nextPage = ""
	}

	return &dto.HistoryTransactionResponse{
		Data:      historys,
		TotalData: int64(totalData),
		TotalPage: int64(totalPage),
		PrevPage:  prevPage,
		NextPage:  nextPage,
	}, nil
}

func (us *UserService) EwalletUS(c context.Context, userID int64) (*dto.EwalletResponse, error) {
	res, err := us.ur.EwalletUR(c, userID)
	if err != nil {
		return nil, err
	}
	return &dto.EwalletResponse{
		Balance: res.Balance,
		Income:  res.Income,
		Expanse: res.Expanse,
	}, nil
}

func (us *UserService) DeleteHistoryUS(c context.Context, query dto.DeleteHistoryRequest) error {
	err := us.ur.DeleteHistoryUR(c, query)
	if err != nil {
		return err
	}
	return nil
}
