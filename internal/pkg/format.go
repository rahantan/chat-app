package pkg

import "time"

const (
	DateTimeLayout = "2006-01-02 15:04:05"
)

// AUTH PATH
const (
	RefreshTokenCookiePath = "/api/public/refresh-token"
	AccessTokenCookiePath  = "/api/protect"
)

// FILE PATH
const (
	BaseImagePath  = "./assets/image"
	UserImagePath  = BaseImagePath + "/user/"
	GroupImagePath = BaseImagePath + "/group/"
)

// URL
const (
	// ImagePath   = "/api/public/products/image?src="
	// MsgFilePath = "/api/public/messages/file?src="

	AvatarImageUrl = "/api/public/user/avatar?src="
	GroupImageUrl  = "/api/public/group/image?src="
	MsgFileUrl     = "/api/protect/messages/file?src="
)

// ROOM
const (
	RoomTypePrivate = "private"
	RoomTypeGroup   = "group"

	GroupRoleMember = "member"
	GroupRoleAdmin  = "admin"
)

// MESSAGE
const (
	MsgTypeText  = "text"
	MsgTypeImage = "image"
	MsgTypeDoc   = "document"
)

var (
	DefaultAvatarImage = "default.png"
	DefaultGroupImage  = "default.png"
)

var (
	AccessTokenExpireTime  = 30 * time.Minute
	RefreshTokenExpireTime = 7 * 24 * time.Hour
)
