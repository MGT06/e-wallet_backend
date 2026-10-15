package repository

import (
	"context"
	"log"

	msgerr "github.com/MGT06/e-wallet_backend/internal/message"
	"github.com/MGT06/e-wallet_backend/internal/models"
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
