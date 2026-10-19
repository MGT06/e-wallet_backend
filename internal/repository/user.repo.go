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
			id,
            fullname, 
            image_url, 
            phone, 
            email
        FROM accounts
        WHERE id = $1
    `
	var user models.Account
	err := ur.db.QueryRow(c, query, userID).Scan(
		&user.ID,
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

func (ur *UserRepo) ToggleFavoriteUR(c context.Context, body dto.AddFavoriteRequest) (string, error) {
	query := `
        WITH target_check AS (
            SELECT account_id FROM favorites 
            WHERE account_id = $1 AND favorite_account_id = $2
        ),
        deleted AS (
            DELETE FROM favorites 
            WHERE account_id = $1 AND favorite_account_id = $2
            RETURNING 'deleted' AS action
        ),
        inserted AS (
            INSERT INTO favorites (account_id, favorite_account_id)
            SELECT $1, $2
            WHERE NOT EXISTS (SELECT 1 FROM target_check)
            RETURNING 'saved' AS action
        )
        SELECT action FROM inserted
        UNION ALL
        SELECT action FROM deleted
    `
	args := []any{body.ID, body.FavoriteID}

	var action string
	err := ur.db.QueryRow(c, query, args...).Scan(&action)
	if err != nil {
		log.Printf("[UserRepo.ToggleFavoriteUR] QueryRow : %s", err)

		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23503" {
				switch pgErr.ConstraintName {
				case "favorites_account_id_fkey":
					return "", msgerr.ErrUserNotExist
				case "favorites_favorite_account_id_fkey":
					return "", msgerr.ErrUSerFavNotFound
				default:
					return "", err
				}
			}
		}
		return "", err
	}

	return action, nil
}

func (ur *UserRepo) DelFavoriteUR(c context.Context, body dto.AddFavoriteRequest) error {
	query := `
        DELETE FROM favorites
        WHERE account_id = $1 AND favorite_account_id = $2
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

func (ur *UserRepo) HistoryUR(c context.Context, userId int64, body dto.HistoryTransactionRequest) (*[]models.HistoryMDL, error) {
	var query strings.Builder
	fmt.Fprintf(&query, `
		SELECT 
			t.id AS transaction_id,
			p.fullname AS fullname,
			p.phone AS phone,
			p.image_url AS image_url,
			t.amount,
			COUNT(tr.sender_id) OVER() AS total,
			CASE 
				WHEN tr.sender_id = $1 THEN 'sender'
				WHEN tr.receiver_id = $1 THEN 'receiver'
			END AS type
		FROM transaction t
		JOIN transfers tr ON t.transfer_id = tr.id
		JOIN accounts p ON p.id = CASE 
			WHEN tr.sender_id = $%d 
			THEN tr.receiver_id 
			ELSE tr.sender_id 
		END
		WHERE t.account_id = $%d
		  AND t.type = 'transfer'
		  AND t.delete_at IS NULL
	`, 1, 1)
	args := []any{userId}
	paramIndex := 2

	hasFullname := body.Fullname != ""
	hasPhone := body.Phone != ""

	if hasFullname {
		fmt.Fprintf(&query, " AND LOWER(p.fullname) LIKE LOWER($%d)", paramIndex)
		args = append(args, "%"+body.Fullname+"%")
		paramIndex += 1
	} else if hasPhone {
		fmt.Fprintf(&query, " AND p.phone LIKE ($%d)", paramIndex)
		args = append(args, "%"+body.Phone+"%")
		paramIndex += 1
	}
	if body.Limit >= 0 {
		fmt.Fprintf(&query, " LIMIT $%d", paramIndex)
		args = append(args, body.Limit)
		paramIndex += 1
	}

	if body.Offset >= 0 {
		fmt.Fprintf(&query, " OFFSET $%d", paramIndex)
		args = append(args, body.Offset)
	}

	finalQuery := query.String()

	rows, err := ur.db.Query(c, finalQuery, args...)
	if err != nil {
		log.Printf("[UserRepo.HistoryUR] Query : %s", err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	historys := make([]models.HistoryMDL, 0)

	for rows.Next() {
		var history models.HistoryMDL
		err := rows.Scan(
			&history.ID,
			&history.Fullname,
			&history.Phone,
			&history.ImageURL,
			&history.Amount,
			&history.TotalData,
			&history.Type,
		)
		if err != nil {
			return nil, err
		}
		historys = append(historys, history)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &historys, nil
}

func (ur *UserRepo) EwalletUR(c context.Context, userID int64) (*models.EwalletMDL, error) {
	query := `
		SELECT 
		    a.balance AS balance,
		    COALESCE(SUM(CASE WHEN t.type = 'topup' THEN t.amount ELSE 0 END), 0) AS income,
		    COALESCE(SUM(CASE WHEN t.type = 'transfer' THEN ABS(t.amount) ELSE 0 END), 0) AS expense
		FROM accounts a
		LEFT JOIN transaction t ON a.id = t.account_id AND t.delete_at IS NULL
		WHERE a.id = $1
		GROUP BY a.id, a.fullname, a.balance
    `
	var ewallet models.EwalletMDL

	err := ur.db.QueryRow(c, query, userID).Scan(
		&ewallet.Balance,
		&ewallet.Income,
		&ewallet.Expanse,
	)

	if err != nil {
		log.Printf("[UserRepo.GetUserHeaderUR] QueryRow : %s", err)
		return nil, msgerr.ErrUserNotExist
	}

	return &ewallet, nil
}
