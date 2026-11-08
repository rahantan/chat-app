package model

import (
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement"`
	UserName  string         `gorm:"column:username;"`
	Email     string         `gorm:"column:email;unique"`
	Image     string         `gorm:"column:image"`
	Password  string         `gorm:"column:password;"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (*UserModel) TableName() string {
	return "users"
}
