package pkg

import "realtime-chat/internal/dto/response"

func NewErr(code response.Code, msg string, details ...response.FieldErr) error {
	return &response.Error{
		Msg:     msg,
		Code:    code,
		Details: details,
	}
}

// AUTH
const (
	ErrCodeUnauthorized       response.Code = "UNAUTHORIZED"
	ErrCodeForbidden          response.Code = "FORBIDDEN"
	ErrCodeTokenExpired       response.Code = "TOKEN_EXPIRED"
	ErrCodeInvalidCredentials response.Code = "INVALID_CREDENTIALS"
)

const (
	ErrCodeTimeOut         response.Code = "TIME_OUT"
	ErrCodeBadRequest      response.Code = "BAD_REQUEST"
	ErrCodeNotFound        response.Code = "NOT_FOUND"
	ErrCodeConflict        response.Code = "CONFLICT"
	ErrCodeValidationError response.Code = "VALIDATION_ERROR"
	ErrCodeInternalError   response.Code = "INTERNAL_ERROR"
)

// CLIENT
const ()

// ErrInvalidExtFile response.Code = "INVALID_EXT_FILE"
