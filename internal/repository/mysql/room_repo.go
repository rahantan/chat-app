package mysql

import (
	"context"
	"errors"
	"realtime-chat/internal/model"
	"realtime-chat/internal/pkg"

	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(ctx context.Context, room *model.RoomModel) error                   // done
	DeleteByID(ctx context.Context, roomID int64) error                        // done
	FindGroupByID(ctx context.Context, roomID int64) (*model.RoomModel, error) // done
	IsExist(ctx context.Context, roomID int64) (bool, error)                   // done

	FindAll(ctx context.Context, userID int64) ([]model.RoomModel, error)
	// GetByID(ctx context.Context) (*model.RoomModel, error)

	FindByID(ctx context.Context, roomID int64) (*model.RoomModel, error)
}

type roomRepoImpl struct {
	baseRepo BaseRepsoitory
}

// members
// members.User

func (rRepo *roomRepoImpl) FindGroupByID(ctx context.Context, roomID int64) (*model.RoomModel, error) {
	var room *model.RoomModel
	if err := rRepo.baseRepo.Conn(ctx).Take(&room, "id = ? AND type = ? ", roomID, pkg.RoomTypeGroup).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrRoomNotFound
		}
		return nil, err
	}
	return room, nil
}

func (rRepo *roomRepoImpl) IsExist(ctx context.Context, roomID int64) (bool, error) {
	var count int64
	if err := rRepo.baseRepo.Conn(ctx).Model(&model.RoomModel{}).Where("id = ?", roomID).Count(&count).Error; err != nil {
		return false, err
	}
	return count == 1, nil
}

func (rRepo *roomRepoImpl) DeleteByID(ctx context.Context, roomID int64) error {
	result := rRepo.baseRepo.Conn(ctx).Delete(&model.RoomModel{}, "id = ?", roomID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return model.ErrRoomNotFound
	}
	return nil
}

func (rRepo *roomRepoImpl) Create(ctx context.Context, room *model.RoomModel) error {
	if err := rRepo.baseRepo.Conn(ctx).Create(room).Error; err != nil {
		return err
	}
	return nil
}

func (rRepo *roomRepoImpl) FindAll(ctx context.Context, userID int64) ([]model.RoomModel, error) {
	db := rRepo.baseRepo.Conn(ctx)

	roomIDs := db.Model(&model.RoomMemberModel{}).Select("room_id").Where("user_id = ?", userID)

	privateRoomIDs := db.Model(&model.RoomModel{}).
		Select("rooms.id").
		Joins("JOIN room_members ON room_members.room_id = rooms.id").
		Where("rooms.type = ?", pkg.RoomTypePrivate).
		Where("room_members.user_id = ?", userID)

	rooms := make([]model.RoomModel, 0)
	if err := db.Model(&model.RoomModel{}).
		Where("id IN (?)", roomIDs).
		Preload("Members", "user_id <> ? AND room_id IN (?)", userID, privateRoomIDs).
		Preload("Members.User").
		Joins(`
			LEFT JOIN messages AS last_messages
			ON last_messages.id = (
				SELECT m.id
				FROM messages AS m
				WHERE m.room_id = rooms.id
				ORDER BY m.created_at DESC, m.id DESC
				LIMIT 1
			)
		`).
		Find(&rooms).Error; err != nil {
		return nil, err
	}

	return rooms, nil
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
