package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/MGT06/e-wallet_backend/internal/dto"
	"github.com/MGT06/e-wallet_backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepo struct {
	db *pgxpool.Pool
}

func NewTransactionRepo(db *pgxpool.Pool) *TransactionRepo {
	return &TransactionRepo{
		db: db,
	}
}

func (tr *TransactionRepo) FindPeopleTR(c context.Context, userID int64, body dto.FindAllPeopleRequest) (*[]models.FindPeopleMDL, error) {
	var query strings.Builder
	fmt.Fprintf(&query, `
    SELECT 
        a.id,
        a.fullname,
        a.phone,
        a.image_url,
        (f.favorite_account_id IS NOT NULL) AS isFavorit
    FROM accounts a
    LEFT JOIN favorites f 
        ON f.favorite_account_id = a.id 
        AND f.account_id = $%d
    WHERE a.id != $%d`, 1, 1)

	args := []any{userID}
	paramIndex := 2

	hasFullname := body.Fullname != ""
	hasPhone := body.Phone != ""

	if hasFullname && hasPhone {
		fmt.Fprintf(&query, " AND LOWER(a.fullname) LIKE LOWER($%d) AND LOWER(a.phone) LIKE LOWER($%d)", paramIndex, paramIndex+1)
		args = append(args, "%"+body.Fullname+"%", "%"+body.Phone+"%")
	} else if hasFullname {
		fmt.Fprintf(&query, " AND LOWER(a.fullname) LIKE LOWER($%d)", paramIndex)
		args = append(args, "%"+body.Fullname+"%")
		paramIndex += 1
	} else if hasPhone {
		fmt.Fprintf(&query, " AND a.phone LIKE ($%d)", paramIndex)
		args = append(args, "%"+body.Phone+"%")
		paramIndex += 1
	}

	if body.IsFavorit {
		fmt.Fprintf(&query, " AND f.favorite_account_id IS NOT NULL ")
	}

	if body.Limit >= 0 {
		fmt.Fprintf(&query, " LIMIT $%d", paramIndex)
		args = append(args, body.Limit)
		paramIndex += 1
	}

	finalQuery := query.String()

	fmt.Println(finalQuery)

	rows, err := tr.db.Query(c, finalQuery, args...)
	if err != nil {
		log.Printf("[TransactionRepo.FindPeopleUR] Query : %s", err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	users := make([]models.FindPeopleMDL, 0)

	for rows.Next() {
		var user models.FindPeopleMDL
		err := rows.Scan(
			&user.ID,
			&user.Fullname,
			&user.Phone,
			&user.ImageURL,
			&user.IsFavorit,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &users, nil
}

func (tr *TransactionRepo) PaymentMethodeTR(c context.Context) (*[]models.PaymentMethodMDL, error) {
	query := `
		SELECT id, name, type, xendit_code, icon_url, tax_percent
		FROM payment_methods
	`
	rows, err := tr.db.Query(c, query)
	if err != nil {
		log.Printf("[TransactionRepo.PaymentMethodeTR] Query : %s", err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	payments := make([]models.PaymentMethodMDL, 0)

	for rows.Next() {
		var payment models.PaymentMethodMDL
		err := rows.Scan(
			&payment.ID,
			&payment.Name,
			&payment.Type,
			&payment.XenditCode,
			&payment.IconURL,
			&payment.TaxPercent,
		)
		if err != nil {
			log.Printf("[TransactionRepo.PaymentMethodeTR] Scan : %s", err.Error())
			return nil, err
		}
		payments = append(payments, payment)
	}
	return &payments, nil
}
