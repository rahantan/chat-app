package pkg

import "time"

var (
	ErrUnauthorized         = NewError(KindUnauthorized, "authentication required", nil)
	ErrNotFound             = NewError(KindNotFound, "not found", nil)
	ErrCustomInvalidPayload = NewError(KindBadRequest, "invalid payload", nil)
	ErrLogin                = NewError(KindNotFound, "invalid email or password", nil)
)

var (
	AccessTokenExpireTime  = 15 * time.Minute
	RefreshTokenExpireTime = 7 * 24 * time.Hour
)
