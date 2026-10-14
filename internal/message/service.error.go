package msgerr

import "errors"

var ErrBodyEmpty = errors.New("payload body harus diisi")
