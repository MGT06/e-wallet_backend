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
