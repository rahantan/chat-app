package model

type RefreshToken struct {
	ID        int64  `gorm:"column:id"`
	UserID    int64  `gorm:"column:user_id"`
	TokenHash string `gorm:"column:token_hash"`
}

func (*RefreshToken) TableName() string {
	return "refresh_token"
}
