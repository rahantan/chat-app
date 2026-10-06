package pkg

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidationErr(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errDetail := make(map[string]string)
		for _, validate := range validationErrors {
			switch validate.Tag() {
			case "email":
				errDetail[validate.Field()] = strings.ToLower("invalid email")
			case "required":
				errDetail[validate.Field()] = strings.ToLower(fmt.Sprintf("%s is required", validate.Field()))
			case "eqfield":
				errDetail[validate.Field()] = strings.ToLower(fmt.Sprintf("%s doesn't match", validate.Field()))
			}
		}
		return NewError(KindBadRequest, "validation errors", errDetail)
	}
	return err
}
