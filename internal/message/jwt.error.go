package msgerr

import "errors"

var (
	ErrMissingKey     = errors.New("jwt key not found")
	ErrUnexpectedSign = errors.New("unexpected signing method")
)
