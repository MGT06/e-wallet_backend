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

type AuthRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepo(db *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (ar *AuthRepo) RegisterAR(c context.Context, body dto.RegisterRequest) (err error) {
	query := `
		INSERT INTO accounts (fullname, email, password)
		VALUES ($1, $2, $3)
	`
	args := []any{body.Name, body.Email, body.Password}

	tag, err := ar.db.Exec(c, query, args...)
	if err != nil {
		log.Printf("[AuthRepo.RegisterAR] Exec : %s", err)
		var repoErr *pgconn.PgError
		if errors.As(err, &repoErr) {
			code, msg := repoErr.Code, repoErr.ConstraintName
			log.Printf("code : %s\nmessage : %s", code, msg)
			if code == "23505" && msg == "accounts_email_key" {
				log.Printf("code : %s\nmessage : %s cek", code, msg)
				return msgerr.ErrEmailExist
			}
		}
		return err
	}

	if tag.RowsAffected() == 0 {
		return msgerr.ErrNoRowAffected
	}
	return nil
}

func (ar *AuthRepo) GetUserByEmail(c context.Context, body dto.LoginRequest) (*models.Account, error) {
	query := `
        SELECT
            password,
            deleted_at
        FROM accounts
        WHERE email = $1
    `
	var user models.Account
	err := ar.db.QueryRow(c, query, body.Email).Scan(
		&user.Password,
		&user.DeletedAt,
	)

	if err != nil {
		log.Printf("[AuthRepo.GetUserByEmail] QueryRow : %s", err)
		return nil, err
	}

	return &user, nil
}

func (ar *AuthRepo) GetUserById(c context.Context, body dto.LoginRequest) (*models.Account, error) {
	query := `
        SELECT 
            id,
            fullname, 
            image_url, 
            phone, 
            email, 
            password, 
            pin, 
            balance, 
            deleted_at,
            created_at, 
            updated_at
        FROM accounts
        WHERE id = $1
    `
	var user models.Account
	err := ar.db.QueryRow(c, query, body.Email).Scan(
		&user.Fullname,
		&user.ImageURL,
		&user.Phone,
		&user.Email,
		&user.Pin,
		&user.Balance,
	)

	if err != nil {
		log.Printf("[AuthRepo.GetUserById] QueryRow : %s", err)
		return nil, err
	}

	return &user, nil
}
