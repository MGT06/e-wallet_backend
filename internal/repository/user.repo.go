package repository

import (
	"context"
	"errors"
	"log"

	"github.com/MGT06/e-wallet_backend/internal/dto"
	msgerr "github.com/MGT06/e-wallet_backend/internal/message"
	"github.com/MGT06/e-wallet_backend/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUSerRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (ur *UserRepo) GetUserHeaderUR(c context.Context, userID int64) (*models.Account, error) {
	query := `
        SELECT 
            fullname, 
            image_url
        FROM accounts
        WHERE id = $1
    `
	var user models.Account
	err := ur.db.QueryRow(c, query, userID).Scan(
		&user.Fullname,
		&user.ImageURL,
	)

	if err != nil {
		log.Printf("[UserRepo.GetUserHeaderUR] QueryRow : %s", err)
		return nil, msgerr.ErrUserNotExist
	}

	return &user, nil
}

func (ur *UserRepo) GetUserProfileUR(c context.Context, userID int64) (*models.Account, error) {
	query := `
        SELECT
            fullname, 
            image_url, 
            phone, 
            email
        FROM accounts
        WHERE id = $1
    `
	var user models.Account
	err := ur.db.QueryRow(c, query, userID).Scan(
		&user.Fullname,
		&user.ImageURL,
		&user.Phone,
		&user.Email,
	)

	if err != nil {
		log.Printf("[UserRepo.GetUserHeaderUR] QueryRow : %s", err)
		return nil, msgerr.ErrUserNotExist
	}

	return &user, nil
}

func (ur *UserRepo) ChangeProfileUR(c context.Context, userID int64, image_path string, body dto.UserChangeProfileRequest) (*models.Account, error) {
	query := `
        UPDATE accounts
        SET 
            fullname = COALESCE(NULLIF($1, ''), fullname),
            image_url = COALESCE(NULLIF($2, ''), image_url),
            phone = COALESCE(NULLIF($3, ''), phone),
            email = COALESCE(NULLIF($4, ''), email)
        WHERE id = $5
        RETURNING 
            fullname, 
            image_url, 
            phone, 
            email
    `
	args := []any{
		body.Fullname,
		image_path,
		body.Phone,
		body.Email,
		userID,
	}
	var profile models.Account
	err := ur.db.QueryRow(
		c,
		query,
		args...,
	).Scan(
		&profile.Fullname,
		&profile.ImageURL,
		&profile.Phone,
		&profile.Email,
	)

	if err != nil {
		log.Printf("[UserRepo.ChangeProfileUR] QueryRow : %s", err)
		return nil, err
	}

	return &profile, nil
}

func (ur *UserRepo) DeletePhotoProfileUR(c context.Context, userID int64) error {
	query := `
        UPDATE accounts
        SET 
            image_url = ''
        WHERE id = $1
    `
	tag, err := ur.db.Exec(c, query, userID)
	if err != nil {
		log.Printf("[UserRepo.DeletePhotoProfileUR] Exec : %s", err)
		var repoErr *pgconn.PgError
		if errors.As(err, &repoErr) {
			code, msg := repoErr.Code, repoErr.ConstraintName
			log.Printf("code : %s\nmessage : %s", code, msg)
			if code == "23505" && msg == "accounts_id_key" {
				log.Printf("code : %s\nmessage : %s cek", code, msg)
				return msgerr.ErrUserNotExist
			}
		}
		return err
	}

	if tag.RowsAffected() == 0 {
		log.Printf("[UserRepo.DeletePhotoProfileUR] RowsAffected : %s", err)
		return msgerr.ErrNoRowAffected
	}

	return nil
}

func (ur *UserRepo) ChangePasswordUR(c context.Context, userID int64, password string) error {
	query := `
        UPDATE accounts
        SET 
            password = $1
        WHERE id = $2
    `
	args := []any{password, userID}

	tag, err := ur.db.Exec(c, query, args...)
	if err != nil {
		log.Printf("[UserRepo.ChangePasswordUR] Exec : %s", err)
		var repoErr *pgconn.PgError
		if errors.As(err, &repoErr) {
			code, msg := repoErr.Code, repoErr.ConstraintName
			log.Printf("code : %s\nmessage : %s", code, msg)
			if code == "23505" && msg == "accounts_id_key" {
				log.Printf("code : %s\nmessage : %s cek", code, msg)
				return msgerr.ErrUserNotExist
			}
		}
		return err
	}

	if tag.RowsAffected() == 0 {
		log.Printf("[UserRepo.DeletePhotoProfileUR] RowsAffected : %s", err)
		return msgerr.ErrNoRowAffected
	}

	return nil

}

func (ur *UserRepo) GetUserUR(c context.Context, userID int64) (*models.Account, error) {
	query := `
        SELECT
            password,
			pin
        FROM accounts
        WHERE id = $1
    `
	var user models.Account
	err := ur.db.QueryRow(c, query, userID).Scan(
		&user.Password,
		&user.Pin,
	)

	if err != nil {
		log.Printf("[UserRepo.GetUserUR] QueryRow : %s", err)
		return nil, msgerr.ErrUserNotExist
	}

	return &user, nil
}

func (ur *UserRepo) SetPinUR(c context.Context, pin string, idUser int64) error {
	query := "UPDATE accounts SET pin = $1 WHERE id = $2"
	args := []any{pin, idUser}

	tag, err := ur.db.Exec(c, query, args...)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return msgerr.ErrNoRowAffected
	}

	return nil
}

func (ur *UserRepo) AddFavoriteUR(c context.Context, body dto.AddFavoriteRequest) error {
	query := `
        INSERT INTO favorites (account_id, favorite_account_id)
        VALUES ($1, $2)
    `
	args := []any{body.ID, body.FavoriteID}

	tag, err := ur.db.Exec(c, query, args...)
	if err != nil {
		log.Printf("[UserRepo.AddFavoriteUR] Exec : %s", err)

		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23503" {
				switch pgErr.ConstraintName {
				case "favorites_account_id_fkey":
					return msgerr.ErrUserNotExist
				case "favorites_favorite_account_id_fkey":
					return msgerr.ErrUSerFavNotFound
				default:
					return err
				}
			}
		}
		return err
	}

	if tag.RowsAffected() == 0 {
		return msgerr.ErrNoRowAffected
	}
	return nil
}

func (ur *UserRepo) ChangePinUR(c context.Context, userId int64, pin string) error {
	query := `
        UPDATE accounts 
		SET pin = $1 
		WHERE id = $2
    `
	args := []any{pin, userId}

	tag, err := ur.db.Exec(c, query, args...)
	if err != nil {
		log.Printf("[UserRepo.ChangePinUR] Exec : %s\n\n%s\n\n", err, pin)
		return err
	}

	if tag.RowsAffected() == 0 {
		log.Printf("[ChangePinUR.RowsAffected] : %s", err)
		return msgerr.ErrNoRowAffected
	}
	return nil
}

func (ur *UserRepo) HistoryUR(c context.Context, userId int64, query dto.HistoryTransactionRequest) (*[]models.HistoryMDL, error) {
	return nil, nil
}
