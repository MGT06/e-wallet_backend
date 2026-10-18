package msgerr

import "errors"

var ErrNoRowAffected = errors.New("no row affected")
var ErrEmailExist = errors.New("email alredy exist")
var ErrEmailNotExist = errors.New("email or password is wrong")
var ErrUserNotExist = errors.New("user not found")
var ErrUSerFavNotFound = errors.New("favorite_account_id not found")
var ErrDataNotFound = errors.New("data not found")
var	ErrInvoiceMismatch = errors.New("mismatch invoice id")
