package model

import "time"

type RoomModel struct {
	ID        int64             `gorm:"column:id"`
	Type      string            `gorm:"column:type"`
	Name      string            `gorm:"column:name"`
	Image     string            `gorm:"column:image"`
	CreatedAt time.Time         `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time         `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Members   []RoomMemberModel `gorm:"foreignKey:RoomID"`
}

type RoomMemberModel struct {
	ID     int64  `gorm:"column:id"`
	Role   string `gorm:"column:role"`
	RoomID int64  `gorm:"column:room_id"`
	UserID int64  `gorm:"column:user_id"`

	User UserModel `gorm:"foreignKey:UserID;references:ID"`
}
