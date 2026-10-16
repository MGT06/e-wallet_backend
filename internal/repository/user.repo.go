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
        SELECT ,
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

func (ur *UserRepo) ChangeProfileUR(ctx context.Context, userID int64, image_path string, body dto.UserChangeProfileRequest) (*models.Account, error) {
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
		ctx,
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

func (ur *UserRepo) DeletePhotoProfileUR(ctx context.Context, userID int64) error {
	query := `
        UPDATE accounts
        SET 
            image_url = ''
        WHERE id = $1
    `
	tag, err := ur.db.Exec(ctx, query, userID)
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
