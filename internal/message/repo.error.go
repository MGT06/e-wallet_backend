package msgerr

import "errors"

var ErrNoRowAffected = errors.New("no row affected")
var ErrEmailExist = errors.New("email alredy exist")
