package mysql

import (
	"context"
	"errors"
	"realtime-chat/internal/model"
	"realtime-chat/internal/pkg"
	"time"

	"gorm.io/gorm"
)

type RoomRepository interface {
	IsExists(ctx context.Context, roomID int64) (bool, error)                            // done
	FindByUserIDWithDetail(ctx context.Context, userID int64) ([]model.RoomModel, error) //ad
	FindByIDWithMembers(ctx context.Context, roomID int64) (*model.RoomModel, error)
	FindByID(ctx context.Context, roomID int64) (*model.RoomModel, error)
	FindPrivateRoomByUsers(ctx context.Context, userID, targetUserID int64) (int64, error)

	Create(ctx context.Context, room *model.RoomModel) error // done
	DeleteByID(ctx context.Context, roomID int64) error      // done

	UpdateTime(ctx context.Context, roomID int64) error
}

type roomRepoImpl struct {
	baseRepo BaseRepsoitory
}

func NewRoomRepository(base BaseRepsoitory) RoomRepository {
	return &roomRepoImpl{
		baseRepo: base,
	}
}
func (rRepo *roomRepoImpl) FindPrivateRoomByUsers(ctx context.Context, userID, targetUserID int64) (int64, error) {
	var roomID int64

	err := rRepo.baseRepo.Conn(ctx).
		Table("rooms r").
		Select("r.id").
		Joins("JOIN room_members me ON me.room_id = r.id").
		Joins("JOIN room_members target ON target.room_id = r.id").
		Where("me.user_id = ?", userID).
		Where("target.user_id = ?", targetUserID).
		Where("r.type = ?", pkg.RoomTypePrivate).
		Take(&roomID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, model.ErrRoomNotFound
		}
		return 0, err
	}
	return roomID, nil
}

func (rRepo *roomRepoImpl) IsExists(ctx context.Context, roomID int64) (bool, error) {
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
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return model.ErrDuplicatePrivateRoom
		}
		return err
	}
	return nil
}

func (rRepo *roomRepoImpl) FindByUserIDWithDetail(ctx context.Context, userID int64) ([]model.RoomModel, error) {
	db := rRepo.baseRepo.Conn(ctx)

	roomIDs := db.Model(&model.RoomMemberModel{}).Select("room_id").Where("user_id = ?", userID)

	rooms := make([]model.RoomModel, 0)
	if err := db.Where("id IN (?)", roomIDs).
		Preload("Members", func(tx *gorm.DB) *gorm.DB {
			return tx.Joins("JOIN rooms ON rooms.id = room_members.room_id").
				Where("rooms.type = ?", pkg.RoomTypePrivate)
		}).
		Preload("Members.User").
		Find(&rooms).Error; err != nil {
		return nil, err
	}
	if len(rooms) == 0 {
		return rooms, nil
	}

	// Kumpulkan ID room
	ids := make([]int64, len(rooms))
	for i := range rooms {
		ids[i] = rooms[i].ID
	}

	// Ambil message terakhir tiap room (1 query)
	lastMsgs := make([]model.MessageModel, 0, len(ids))
	if err := rRepo.baseRepo.Conn(ctx).
		Where(`id IN (
		SELECT MAX(id) FROM messages
		WHERE room_id IN ?
		GROUP BY room_id
	)`, ids).
		Preload("Sender", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "username")
		}).
		Find(&lastMsgs).Error; err != nil {
		return nil, err
	}

	lastMsgMap := make(map[int64]*model.MessageModel, len(lastMsgs))
	for i := range lastMsgs {
		lastMsgMap[lastMsgs[i].RoomID] = &lastMsgs[i]
	}
	for i := range rooms {
		rooms[i].LastMessage = lastMsgMap[rooms[i].ID]
	}
	return rooms, nil
}

func (rRepo *roomRepoImpl) FindByIDWithMembers(ctx context.Context, roomID int64) (*model.RoomModel, error) {
	var room model.RoomModel
	if err := rRepo.baseRepo.Conn(ctx).Preload("Members.User").Take(&room, "id = ?", roomID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrRoomNotFound
		}
		return nil, err
	}
	return &room, nil
}

func (rRepo *roomRepoImpl) FindByID(ctx context.Context, roomID int64) (*model.RoomModel, error) {
	var room *model.RoomModel
	if err := rRepo.baseRepo.Conn(ctx).Where("id = ?", roomID).Take(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrRoomNotFound
		}
		return nil, err
	}

	return room, nil
}

func (rRepo *roomRepoImpl) UpdateTime(ctx context.Context, roomID int64) error {
	result := rRepo.baseRepo.Conn(ctx).Model(&model.RoomModel{}).Where("id = ?", roomID).UpdateColumn("updated_at", time.Now())
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return model.ErrRoomNotFound
	}
	return nil
}
