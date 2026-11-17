package model

import (
	"time"
)

type MessageAttachmentModel struct {
	ID       int64  `gorm:"column:id;primaryKey"`
	MsgID    int64  `gorm:"column:message_id"`
	FileName string `gorm:"column:file_url"`
	FileType string `gorm:"column:file_type"`
	FileSize string `gorm:"column:file_size"`
}

func (*MessageAttachmentModel) TableName() string {
	return "message_attachments"
}

type MessageModel struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	SenderID  int64     `gorm:"column:user_id"`
	RoomID    int64     `gorm:"column:room_id"`
	MsgType   string    `gorm:"column:message_type"`
	Content   string    `gorm:"column:content"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`

	Sender      UserModel                `gorm:"foreignKey:SenderID;references:ID"`
	Attachments []MessageAttachmentModel `gorm:"foreignKey:MsgID;references:ID"`
}

func (*MessageModel) TableName() string {
	return "messages"
}
