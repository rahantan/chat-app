package model

// type UserContactModel struct {
// 	ID     int64  `gorm:"column:id;primaryKey;autoIncrement"`
// 	RoomID int64  `gorm:"column:room_id"`
// 	Status string `gorm:"column:status"`
// }

// type ContactModel struct {
// 	ID            int64     `gorm:"column:id;primaryKey;autoIncrement"`
// 	UserID        int64     `gorm:"column:user_id"`
// 	UserContactID int64     `gorm:"column:user_contact_id"`
// 	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
// 	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`

// 	UserContact UserContactModel `gorm:"foreignKey:UserContactID:references:ID"`
// 	User        UserModel        `gorm:"foreignKey:UserID;references:ID"`
// }

type ContactModel struct {
	ID            int64 `gorm:"column:id;primaryKey;autoIncrement"`
	UserID        int64 `gorm:"column:user_id"`
	ContactUserID int64 `gorm:"column:contact_user_id"`
	// CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	// UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`

	User UserModel `gorm:"foreignKey:ContactUserID;references:ID"`
}

func (*ContactModel) TableName() string {
	return "contacts"
}

type ContactResponse struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement"`
	ContactUserID int64     `gorm:"column:contact_user_id"`
	UserName      string    `gorm:"column:username"`
	Avatar        string    `gorm:"column:avatar"`
	RoomID        int64     `gorm:"column:room_id"`
	User          UserModel `gorm:"foreignKey:ContactUserID;references:ID"`
}
