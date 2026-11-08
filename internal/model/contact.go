package model

import "time"

type ContactModel struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID        int64     `gorm:"column:user_id"`
	ContactUserID int64     `gorm:"column:contact_user_id"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`

	User UserModel `gorm:"foreignKey:ContactUserID;references:ID"`
}
