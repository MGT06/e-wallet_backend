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
			return fmt.Sprintf("Field %s wajib diisi.", field)
		case "email":
			return fmt.Sprintf("Format email pada field %s tidak valid.", field)
		case "min":
			return fmt.Sprintf("Field %s minimal harus memiliki %s karakter.", field, fe.Param())
		case "max":
			return fmt.Sprintf("Field %s maksimal memiliki %s karakter.", field, fe.Param())
		default:
			return fmt.Sprintf("Field %s tidak valid pada rule %s.", field, fe.Tag())
		}
	}
	return ""
}
