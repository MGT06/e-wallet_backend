package repository

import (
	"context"
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

func (ar *AuthRepo) RegisterAR(c context.Context, db DbConn, body dto.RegisterRequest) (err error) {
	query := `
		INSERT INTO accounts (fullname, email, password)
		VALUES ($1, $2, $3)
	`
	args := []any{body.Name, body.Email, body.Password}

	tag, err := db.Exec(c, query, args...)
	if err != nil {
		log.Printf("[AuthRepo.RegisterAR] Exec : %s", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return msgerr.ErrNoRowAffected
	}
	return nil
}
