package pkg

// AUTH ERROR
var (
	ErrUnauthenticated        = NewErr(ErrCodeUnauthorized, "user is not authenticated")
	ErrInvalidToken           = NewErr(ErrCodeUnauthorized, "invalid token")
	ErrTokenExpired           = NewErr(ErrCodeUnauthorized, "token expired")
	ErrTokenRequired          = NewErr(ErrCodeUnauthorized, "token is required")
	ErrInvalidEmailOrPassword = NewErr(ErrCodeInvalidCredentials, "invalid email or password")
	ErrTokenRevoked           = NewErr(ErrCodeUnauthorized, "token has been revoked")
	ErrForbidden              = NewErr(ErrCodeForbidden, "you don't have permission to perform this action")
	ErrAccessDenied           = NewErr(ErrCodeForbidden, "you don't have access to this resource")
)

// GENERAL
var (
	ErrInvalidParamID = NewErr(ErrCodeBadRequest, "invalid parameter id")
	ErrInvalidPayload = NewErr(ErrCodeBadRequest, "invalid payload")
	ErrNotFound       = NewErr(ErrCodeNotFound, "not found")
)

// CONTACT ERROR
// MESSAGE ERROR
// GROUP ERROR
