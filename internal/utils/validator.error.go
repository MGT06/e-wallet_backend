package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidationError(errs validator.ValidationErrors) string {
	for _, fe := range errs {
		field := fe.Field()
		switch fe.Tag() {
		case "required":
			return fmt.Sprintf("field %s wajib diisi.", field)
		case "email":
			return fmt.Sprintf("format %s tidak valid.", field)
		case "numeric":
			return fmt.Sprintf("field %s harus berupa angka.", field)
		case "min":
			return fmt.Sprintf("field %s minimal harus memiliki %s karakter.", field, fe.Param())
		case "max":
			return fmt.Sprintf("field %s maksimal memiliki %s karakter.", field, fe.Param())
		case "len":
			return fmt.Sprintf("field %s harus 6 angka.", field)
		default:
			return fmt.Sprintf("field %s tidak valid pada rule %s.", field, fe.Tag())
		}
	}
	return ""
}
