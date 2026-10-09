package pkg

type Kind string
type ErrorCode string

type AppErr struct {
	Kind Kind
	// Code    ErrorCode
	Message string
	Details any
}

func NewError(kind Kind, msg string, details any) error {
	return &AppErr{
		Kind:    kind,
		Message: msg,
		// Code:    code,
		Details: details,
	}
}

func (appErr *AppErr) Error() string {
	return appErr.Message
}

const (
	// KindInternal     Kind = "internal server error"
	KindBadRequest Kind = "bad request"
	// KindValidation   Kind = "validation errors"
	KindUnauthorized Kind = "unauthorized"
	KindNotFound     Kind = "not found"
)

// const (
// 	CodeValidation      ErrorCode = "VALIDATION_ERROR"
// 	CodeUnauthenticated ErrorCode = "UNAUTHENTICATED"
// )
