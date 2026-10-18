package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/MGT06/e-wallet_backend/internal/dto"
	msgerr "github.com/MGT06/e-wallet_backend/internal/message"
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

	if body.Limit >= 0 {
		fmt.Fprintf(&query, " LIMIT $%d", paramIndex)
		args = append(args, body.Limit)
		paramIndex += 1
	}

	finalQuery := query.String()

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

func (tr *TransactionRepo) TopupsTR(c context.Context, IdUser int64, body models.Topup) error {
	query := "INSERT INTO topups (reference_no, account_id, payment_method_id, amount, tax, total_amount) VALUES ($1, $2, $3, $4, $5, $6)"
	args := []any{body.ReferenceNo, body.AccountID, body.PaymentMethodID, body.Amount, body.Tax, body.TotalAmount}

	tag, err := tr.db.Exec(c, query, args...)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return	msgerr.ErrNoRowAffected
	}

	return nil
}

func (tr *TransactionRepo) UpdateTopupsTR(c context.Context, body models.Topup) (dto.TopupsResponse, error) {
	query := "UPDATE topups SET xendit_invoice_id = $1, payment_url = $2, expires_at = $3 WHERE reference_no = $4 RETURNING payment_url, expires_at"
	args := []any{*body.XenditID, *body.PaymentURL, *body.ExpiresAt, body.ReferenceNo}

	res := tr.db.QueryRow(c, query, args...)
	var resp dto.TopupsResponse
	if err := res.Scan(&resp.PaymentURL, &resp.ExpiresAt); err != nil {
		return dto.TopupsResponse{}, err
	}

	return resp, nil
}