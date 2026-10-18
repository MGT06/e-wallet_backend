package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"unicode"

	"github.com/MGT06/e-wallet_backend/internal/dto"
	msgerr "github.com/MGT06/e-wallet_backend/internal/message"
	"github.com/MGT06/e-wallet_backend/internal/models"
	"github.com/MGT06/e-wallet_backend/internal/repository"
	"github.com/MGT06/e-wallet_backend/internal/utils"
	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionService struct {
	tr    *repository.TransactionRepo
	resty *resty.Client
	db    *pgxpool.Pool
}

func NewTransactionService(tr *repository.TransactionRepo, resty *resty.Client, db *pgxpool.Pool) *TransactionService {
	return &TransactionService{
		tr:    tr,
		resty: resty,
		db:    db,
	}
}

func (ts *TransactionService) FindPeopleTS(c context.Context, userID int64, query dto.FindAllPeopleRequest) (*dto.FindAllPeopleResponse, error) {
	if query.Limit <= 0 {
		query.Limit = 10
	}
	var isByName bool
	if query.Search == "" {
		isByName = false
	}

	for _, r := range query.Search {
		if unicode.IsNumber(r) {
			isByName = false
		}
		if unicode.IsLetter(r) {
			isByName = true
		} else {
			isByName = false
		}
	}

	if isByName {
		query.Fullname = query.Search
	} else {
		query.Phone = query.Search
	}

	res, err := ts.tr.FindPeopleTR(c, ts.db, userID, query)
	if err != nil {
		return nil, err
	}

	users := make([]dto.FindAllPeople, 0, len(*res))
	for _, u := range *res {
		users = append(users, dto.FindAllPeople{
			ID:        u.ID,
			Fullname:  u.Fullname,
			ImageURL:  u.ImageURL,
			Phone:     u.Phone,
			IsFavorit: u.IsFavorit,
		})
	}

	return &dto.FindAllPeopleResponse{
		Data: users,
	}, nil
}

func (ts *TransactionService) PaymentMethodeTS(c context.Context) (*[]dto.PaymentMethodResponse, error) {
	res, err := ts.tr.PaymentMethodeTR(c, ts.db)
	if err != nil {
		return nil, err
	}

	payments := make([]dto.PaymentMethodResponse, 0, len(*res))
	for _, p := range *res {
		payment := dto.PaymentMethodResponse{
			ID:         p.ID,
			Name:       p.Name,
			Type:       p.Type,
			XenditCode: p.XenditCode,
			IconURL:    p.IconURL,
			TaxPercent: p.TaxPercent,
		}
		payments = append(payments, payment)
	}
	return &payments, nil
}

func (ts *TransactionService) TopupsTS(ctx context.Context, IdUser int64, body dto.TopupsRequest) (dto.TopupsResponse, error) {
	timestamp := time.Now()
	randstring := utils.GenerateRandomString(6)
	refno := fmt.Sprintf("TUP-%d-%s", timestamp.Unix(), randstring)

	if err := ts.tr.TopupsTR(ctx, ts.db, IdUser, models.Topup{
		ReferenceNo:     refno,
		AccountID:       IdUser,
		PaymentMethodID: body.PaymentMethod,
		Amount:          body.Amount,
		Tax:             body.Tax,
		TotalAmount:     body.TotalAmount,
	}); err != nil {
		return dto.TopupsResponse{}, err
	}

	var result map[string]any
	var errResult map[string]any

	sucURL := os.Getenv("XENDIT_SUCCESS_REDIRECT_URL")
	faiURL := os.Getenv("XENDIT_FAILURE_REDIRECT_URL")

	duration, _ := strconv.Atoi(os.Getenv("XENDIT_INVOICE_DURATION"))
	resp, err := ts.resty.R().
		SetContext(ctx).
		SetBasicAuth(os.Getenv("XENDIT_SECRET_KEY"), "").
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"external_id":          refno,
			"amount":               body.TotalAmount,
			"payment_methods":      []string{body.XenditCode},
			"invoice_duration":     duration,
			"success_redirect_url": sucURL,
			"failure_redirect_url": faiURL,
		}).SetResult(&result).SetError(&errResult).Post(fmt.Sprintf("%s/v2/invoices", os.Getenv("XENDIT_API_URL")))
	if err != nil {
		return dto.TopupsResponse{}, err
	}

	if resp.IsError() {
		return dto.TopupsResponse{}, fmt.Errorf("xendit error (status %d): %v", resp.StatusCode(), errResult)
	}

	xenditID, _ := result["id"].(string)
	url, _ := result["invoice_url"].(string)
	expire, _ := result["expiry_date"].(time.Time)
	exid, _ := result["external_id"].(string)

	res, err := ts.tr.UpdateTopupsTR(ctx, ts.db, models.Topup{
		XenditID:    &xenditID,
		PaymentURL:  &url,
		ExpiresAt:   &expire,
		ReferenceNo: exid,
	})

	if err != nil {
		return dto.TopupsResponse{}, err
	}

	return res, nil
}

func (ts *TransactionService) Transfer(ctx context.Context, senderID int64, req dto.TransferRequest) error {
	if senderID == req.ReceiverID {
		return msgerr.ErrSelfTransfer
	}

	timestamp := time.Now()
	randstring := utils.GenerateRandomString(6)
	refno := fmt.Sprintf("TRF-%d-%s", timestamp.Unix(), randstring)

	tx, err := ts.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Println(err.Error())
		}
	}()

	balances, err := ts.tr.GetBalances(ctx, ts.db, senderID, req.ReceiverID)
	if err != nil {
		return err
	}

	senderBefore:= balances[senderID]

	receiverBefore, ok := balances[req.ReceiverID]
	if !ok {
		return msgerr.ErrReceiverNotFound
	}

	if senderBefore < req.Amount {
		return msgerr.ErrInsufficientBalance
	}

	transferID, err := ts.tr.InsertTransfer(ctx, ts.db, models.Transfer{
		ReferenceNo: refno,
		SenderID:    senderID,
		ReceiverID:  req.ReceiverID,
		Amount:      req.Amount,
		Notes:       &req.Notes,
		Status:      "success",
	})
	if err != nil {
		return err
	}

	if err := ts.tr.AddBalance(ctx, ts.db, senderID, -req.Amount); err != nil {
		return err
	}

	if err := ts.tr.AddBalance(ctx, ts.db, req.ReceiverID, req.Amount); err != nil {
		return err
	}

	if err := ts.tr.InsertTransaction(ctx, ts.db, models.Transaction{
		AccountID:     senderID,
		Type:          "transfer",
		Amount:        -req.Amount,
		BalanceBefore: senderBefore,
		BalanceAfter:  senderBefore - req.Amount,
		TransferID:    &transferID,
	}); err != nil {
		return err
	}

	if err := ts.tr.InsertTransaction(ctx, ts.db, models.Transaction{
		AccountID:     req.ReceiverID,
		Type:          "topup",
		Amount:        req.Amount,
		BalanceBefore: receiverBefore,
		BalanceAfter:  receiverBefore + req.Amount,
		TransferID:    &transferID,
	}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
