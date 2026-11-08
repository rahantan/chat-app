package model

import "time"

type MessageAttachmentModel struct {
	ID       int64  `gorm:"column:id"`
	MsgID    int64  `gorm:"column:message_id"`
	FileName string `gorm:"column:file_url"`
	FileType string `gorm:"column:file_type"`
	FileSize string `gorm:"column:file_size"`
}

type MessageModel struct {
	ID        int64     `gorm:"column:id"`
	SenderID  int64     `gorm:"column:sender_id"`
	RoomID    int64     `gorm:"column:room_id"`
	MsgType   string    `gorm:"column:message_type"`
	Content   string    `gorm:"column:content"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`

	Sender      UserModel `gorm:"foreignKey:UserID;references:ID"`
	Attachments []MessageAttachmentModel
}
