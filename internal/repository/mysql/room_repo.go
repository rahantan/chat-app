package mysql

import (
	"context"
	"realtime-chat/internal/model"
)

type RoomRepository interface {
	Create(ctx context.Context, room *model.RoomModel) error
	DeleteByID(ctx context.Context, roomID int64) error
	FindGroupByID(ctx context.Context, roomID int64) (*model.RoomModel, error)
	// List(ctx context.Context, userID int64) ([]model.RoomModel, error)
	// GetByID(ctx context.Context) (*model.RoomModel, error)
	IsExist(ctx context.Context, roomID int64) (bool, error)
}

type roomRepoImpl struct {
	baseRepo BaseRepsoitory
}

// func (repo *roomRepoImpl) List(ctx context.Context, userID int64) ([]model.RoomModel, error) {
// 	var rooms []model.RoomModel
// 	query := repo.baseRepo.Conn(ctx).Table("rooms r")
// 	querySelect := query.Select(``)
// 	join := querySelect.Joins("LEFT JOIN room_members rm ON rm.room_id=r.id AND ")
// }

// func (repo *roomRepoImpl) random() {
// 	err := db.Table("rooms r").
// 		Select(`r.id, r.type,
// 		COALESCE(r.title, u.name)            AS name,
// 		COALESCE(r.image_url, u.avatar_url)  AS image,
// 		m.content AS last_message, m.created_at AS last_message_at`).
// 		Joins("JOIN room_members me ON me.room_id = r.id AND me.user_id = ?", userID).
// 		Joins("LEFT JOIN room_members o ON o.room_id = r.id AND o.user_id <> ? AND r.type = 'private'", userID).
// 		Joins("LEFT JOIN users u ON u.id = o.user_id").
// 		Joins("LEFT JOIN messages m ON m.id = r.last_message_id").
// 		Order("m.created_at DESC").
// 		Scan(&items).Error
// }
