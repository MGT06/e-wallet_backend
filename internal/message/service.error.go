package msgerr

import "errors"

var ErrBodyEmpty = errors.New("payload body harus diisi")
var ErrTokenExpired = errors.New("session expired, login again")
var ErrSesionLogout = errors.New("session out, login again")
var ErrAccounDeleted = errors.New("account deleted")
var ErrEmailNoExists = errors.New("email is not exists")
var ErrExt = errors.New("ext must jpg/png/jpeg")
