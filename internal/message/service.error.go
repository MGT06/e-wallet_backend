package msgerr

import "errors"

var ErrBodyEmpty = errors.New("payload body harus diisi")
var ErrTokenExpired = errors.New("session expired, login again")
var ErrSesionLogout = errors.New("session out, login again")
var ErrAccounDeleted = errors.New("account deleted")
var ErrExt = errors.New("ext accept only jpg/png/jpeg")
var ErrImageToLarge = errors.New("image to large max 2 mb")
var ErrWriteFile = errors.New("cannot write file")
var ErrOpenFile = errors.New("can not open file image")
var ErrReadFile = errors.New("can not read file image")
var ErrPhoneFormat = errors.New("format phone harus berawalan 0/62")
var ErrImgNotFound = errors.New("photo profile belum ada")
var ErrOldPin = errors.New("invalid pin")
var ErrOldPassword = errors.New("invalid password")
var ErrDuplicateWebhook = errors.New("duplicate webhook")
var ErrMismatchPaidAmount = errors.New("transaction failed: amount does not match")
var ErrInvalidCallbackToken = errors.New("invalid callback token")
var ErrTooManyRequests = errors.New("too many requests")
var ErrInvalidOTP = errors.New("invalid or expired otp")
var ErrInvalidToken = errors.New("invalid or expired token")
var ErrTooManyAttempts = errors.New("too many attempts")
