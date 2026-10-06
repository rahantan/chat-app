package pkg

type Kind string

type AppErr struct {
	Kind    string
	Message string
	Details any
}

func NewError(kind Kind, msg string, details any) *AppErr {
	return &AppErr{
		Kind:    string(kind),
		Message: msg,
		Details: details,
	}
}

func (appErr *AppErr) Error() string {
	return appErr.Message
}

const (
	KindInternal     Kind = "internal server error"
	KindBadRequest   Kind = "bad request"
	KindValidation   Kind = "validation error"
	KindUnauthorized Kind = "unauthorized"
	KindNotFound     Kind = "not found"
)
