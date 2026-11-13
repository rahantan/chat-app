package mysql

import (
	"context"
	"errors"
	"realtime-chat/internal/model"
	"realtime-chat/internal/pkg"

	"gorm.io/gorm"
)

type RoomMemberRepository interface {
	FindByUserIDAndRoomID(ctx context.Context, userID, roomID int64) (*model.RoomMemberModel, error)
	AddMembers(ctx context.Context, roomID int64, userIDs []int64) error
}

type roomMemberRepositoryImpl struct {
	baseRepo BaseRepsoitory
}

func NewRoomMemberRepo(repo BaseRepsoitory) RoomMemberRepository {
	return &roomMemberRepositoryImpl{
		baseRepo: repo,
	}
}

func (rmRepo *roomMemberRepositoryImpl) FindByUserIDAndRoomID(ctx context.Context, userID, roomID int64) (*model.RoomMemberModel, error) {
	var member model.RoomMemberModel
	if err := rmRepo.baseRepo.Conn(ctx).Take(&member, "user_id= ? AND room_id= ? ", userID, roomID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}
	return &member, nil
}

func (rmRepo *roomMemberRepositoryImpl) AddMembers(ctx context.Context, roomID int64, userIDs []int64) error {
	members := make([]model.RoomMemberModel, 0, len(userIDs))
	for _, id := range userIDs {
		members = append(members, model.RoomMemberModel{
			Role:   pkg.RoleMember,
			RoomID: roomID,
			UserID: id,
		})
	}

	if err := rmRepo.baseRepo.Conn(ctx).CreateInBatches(&members, 100).Error; err != nil {
		return err
	}

	return nil
}
