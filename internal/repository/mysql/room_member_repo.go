package mysql

import (
	"context"
	"realtime-chat/internal/model"
)

type RoomMemberRepository interface {
	FindByUserIDAndRoomID(ctx context.Context, userID, roomID int64) (*model.RoomMemberModel, error)
	FindByRoomID(ctx context.Context) ([]model.RoomMemberModel, error)
	AddMembers(ctx context.Context, memberIDs []int64) error
}
