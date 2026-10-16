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

// func (tr *TransactionRepo) FindPeopleUR(c context.Context, userID int64, body dto.FindAllPeopleRequest) (*[]models.FindPeopleMDL, error) {
// 	query := `
// SELECT
//     a.id,
//     a.fullname,
//     a.phone,
//     a.image_url,
//     (f.favorite_account_id IS NOT NULL) AS isFavorit
// FROM accounts a
// LEFT JOIN favorites f
//     ON f.favorite_account_id = a.id
//     AND f.account_id = $1
// WHERE a.id != $1
//     AND (
//         LOWER(a.fullname) LIKE LOWER('%' || $2 || '%')
//         OR LOWER(a.phone) LIKE LOWER('%' || $3 || '%')
//     )
//     `

// 	args := []any{userID, body.Fullname, body.Phone}
// 	rows, err := tr.db.Query(c, query, args...)
// 	if err != nil {
// 		log.Printf("[UserRepo.FindPeopleAR] Query : %s", err.Error())
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	users := make([]models.FindPeopleMDL, 0)

// 	for rows.Next() {
// 		var user models.FindPeopleMDL
// 		err := rows.Scan(
// 			&user.ID,
// 			&user.Fullname,
// 			&user.Phone,
// 			&user.ImageURL,
// 			&user.IsFavorit,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		users = append(users, user)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return &users, nil
// }

func (tr *TransactionRepo) FindPeopleUR(c context.Context, userID int64, body dto.FindAllPeopleRequest) (*[]models.FindPeopleMDL, error) {
	var queryBuilder strings.Builder
	queryBuilder.WriteString(`
		SELECT 
			a.id,
			a.fullname,
			a.phone,
			a.image_url,
			(f.favorite_account_id IS NOT NULL) AS isFavorit
		FROM accounts a
		LEFT JOIN favorites f 
			ON f.favorite_account_id = a.id 
			AND f.account_id = $1
		WHERE a.id != $1
	`)

	args := []any{userID}
	paramIndex := 2

	hasFullname := body.Fullname != ""
	hasPhone := body.Phone != ""

	fmt.Printf("\n\nKondisi Fullname : %v\n", hasFullname)
	fmt.Printf("Kondisi Phone : %v", hasPhone)

	if hasFullname && hasPhone {
		queryBuilder.WriteString(fmt.Sprintf(" AND LOWER(a.fullname) LIKE LOWER($%d) AND LOWER(a.phone) LIKE LOWER($%d)", paramIndex, paramIndex+1))
		args = append(args, "%"+body.Fullname+"%", "%"+body.Phone+"%")
	} else if hasFullname {
		queryBuilder.WriteString(fmt.Sprintf(" AND LOWER(a.fullname) LIKE LOWER($%d)", paramIndex))
		args = append(args, "%"+body.Fullname+"%")
	} else if hasPhone {
		queryBuilder.WriteString(fmt.Sprintf(" AND LOWER(a.phone) LIKE LOWER($%d)", paramIndex))
		args = append(args, "%"+body.Phone+"%")
	}

	fmt.Printf("\n\nHasil Query : %v", &queryBuilder)

	rows, err := tr.db.Query(c, queryBuilder.String(), args...)
	if err != nil {
		log.Printf("[UserRepo.FindPeopleUR] Query : %s", err.Error())
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
