package pkg

import (
	"net/http"
	"realtime-chat/internal/dto/response"
)

func GetStatusCode(err error) int {
	errCode, ok := err.(*response.Error)
	if ok {
		switch errCode.Code {
		case ErrCodeBadRequest:
			return http.StatusBadRequest
		case ErrCodeUnauthorized, ErrCodeInvalidCredentials, ErrCodeTokenExpired:
			return http.StatusUnauthorized
		case ErrCodeForbidden:
			return http.StatusForbidden
		case ErrCodeNotFound:
			return http.StatusNotFound
		case ErrCodeConflict:
			return http.StatusConflict
		case ErrCodeValidationError:
			return http.StatusUnprocessableEntity
		}
	}
	return http.StatusInternalServerError
}
