package mysql

import (
	"context"
	"errors"
	"realtime-chat/internal/model"

	"gorm.io/gorm"
)

type MessageRepository interface {
	FindByRoomIDWithDetail(ctx context.Context, roomID int64) ([]model.MessageModel, error)
	FindByID(ctx context.Context, msgID int64) (*model.MessageModel, error)
	Delete(ctx context.Context, userID, roomID, msgID int64) error
	Add(ctx context.Context, msg *model.MessageModel) error
}

type messageRepositoryImpl struct {
	baseRepo BaseRepsoitory
}

func NewMessageRepository(repo BaseRepsoitory) MessageRepository {
	return &messageRepositoryImpl{
		baseRepo: repo,
	}
}
func (msgRepo *messageRepositoryImpl) FindByRoomIDWithDetail(ctx context.Context, roomID int64) ([]model.MessageModel, error) {
	messages := make([]model.MessageModel, 0)

	if err := msgRepo.baseRepo.Conn(ctx).Joins("Sender").Preload("Attachments").Find(&messages, "room_id = ?", roomID).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

func (msgRepo *messageRepositoryImpl) FindByID(ctx context.Context, msgID int64) (*model.MessageModel, error) {
	var message model.MessageModel

	if err := msgRepo.baseRepo.Conn(ctx).Take(&message, "id = ?", msgID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrMessageNotFound
		}
		return nil, err
	}

	return &message, nil
}

func (msgRepo *messageRepositoryImpl) Delete(ctx context.Context, userID, roomID, msgID int64) error {
	result := msgRepo.baseRepo.Conn(ctx).Delete(&model.MessageModel{}, "user_id = ? AND room_id = ? AND id = ? ", userID, roomID, msgID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return model.ErrMessageNotFound
	}
	return nil
}

func (msgRepo *messageRepositoryImpl) Add(ctx context.Context, msg *model.MessageModel) error {
	if err := msgRepo.baseRepo.Conn(ctx).Create(&msg).Error; err != nil {
		return err
	}
	return nil
}
