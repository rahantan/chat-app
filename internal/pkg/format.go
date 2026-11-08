package pkg

import "time"

const (
	DateTimeLayout = "2006-01-02 15:04:05"
)

const (
	RefreshTokenCookiePath = "/api/public/refresh-token"
	AccessTokenCookiePath  = "/api/protect"
	ImagePath              = "/api/public/products/image?src="
	MsgFilePath            = "/api/public/messages/file?src="
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

const (
	RoomTypePrivate = "private"
	RoomTypeGroup   = "group"

	RoleMember = "member"
	RoleAdmin  = "admin"
)

const (
	MsgTypeText  = "text"
	MsgTypeImage = "image"
	MsgTypeDoc   = "document"
)
