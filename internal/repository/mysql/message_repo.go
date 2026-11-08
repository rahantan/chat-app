package mysql

import (
	"context"
	"realtime-chat/internal/model"
)

type MessageRepository interface {
	Add(ctx context.Context, msg *model.MessageModel) error
	FindAllByRoomID(ctx context.Context, roomID int64) ([]model.MessageModel, error)
	FindByID(ctx context.Context, roomID int64) (*model.MessageModel, error)
	Delete(ctx context.Context, userID, roomID, msgID int64) error
}
