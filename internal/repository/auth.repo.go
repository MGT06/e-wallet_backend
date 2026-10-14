package repository

import (
	"context"
	"errors"
	"log"

	"github.com/MGT06/e-wallet_backend/internal/dto"
	msgerr "github.com/MGT06/e-wallet_backend/internal/message"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DbConn interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type AuthRepo struct{}

func NewAuthRepo() *AuthRepo {
	return &AuthRepo{}
}

func (ar *AuthRepo) RegisterAR(c context.Context, db DbConn, body dto.RegisterRequest) (err error) {
	query := `
		INSERT INTO accounts (fullname, email, password)
		VALUES ($1, $2, $3)
	`
	args := []any{body.Name, body.Email, body.Password}

	tag, err := db.Exec(c, query, args...)
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
