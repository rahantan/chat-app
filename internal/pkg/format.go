package pkg

import "time"

const (
	RefreshTokenCookiePath = "/api/public/refresh-token"
	AccessTokenCookiePath  = "/api/protect"
)

var (
	AccessTokenExpireTime  = 1 * time.Minute
	RefreshTokenExpireTime = 7 * 24 * time.Hour
)

var (
	ErrUnauthorized         = NewError(KindUnauthorized, "authentication required", nil)
	ErrTokenRequired        = NewError(KindUnauthorized, "token is required", nil)
	ErrNotFound             = NewError(KindNotFound, "not found", nil)
	ErrCustomInvalidPayload = NewError(KindBadRequest, "invalid payload", nil)
	ErrLogin                = NewError(KindNotFound, "invalid email or password", nil)
	ErrUserNotFound         = NewError(KindNotFound, "user not found", nil)
)
