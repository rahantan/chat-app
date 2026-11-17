package pkg

import (
	"fmt"
	"realtime-chat/internal/dto/response"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidationError(err error) error {
	if validationErr, ok := err.(validator.ValidationErrors); ok {
		errDetail := make([]response.FieldErr, 0)
		for _, validate := range validationErr {

			field := validate.Field()

			var msg string

			switch validate.Tag() {
			case "email":
				msg = "invalid email"

			case "required":
				msg = fmt.Sprintf("%s is required", validate.Field())

			case "min":
				msg = fmt.Sprintf("%s must be at least %s characters", validate.Field(), validate.Param())

			case "max":
				msg = fmt.Sprintf("%s must be at most %s characters", validate.Field(), validate.Param())

			case "eqfield":
				msg = fmt.Sprintf("%s doesn't match", validate.Field())

			case "numeric":
				msg = fmt.Sprintf("%s must be number", validate.Field())

			case "gt":
				msg = fmt.Sprintf("%s must be at least 1", validate.Field())

			}
			errDetail = append(errDetail, response.FieldErr{
				Field: field,
				Msg:   strings.ToLower(msg),
			})
		}
		return NewErr(ErrCodeValidationError, "validation errors", errDetail...)
	}
	return nil
}
