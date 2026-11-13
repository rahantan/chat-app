package mysql

import (
	"context"
	"errors"
	"realtime-chat/internal/model"

	"gorm.io/gorm"
)

type ContactRepository interface {
	FindAll(ctx context.Context, userID int64) ([]model.ContactModel, error)
	FindContact(ctx context.Context, userID, contactUserID int64) (*model.ContactModel, error)
	Add(ctx context.Context, contact *model.ContactModel) error
}

type contactRepositoryImpl struct {
	baseRepo BaseRepsoitory
}

func NewContact(base BaseRepsoitory) ContactRepository {
	return &contactRepositoryImpl{
		baseRepo: base,
	}
}

func (contactRepo *contactRepositoryImpl) FindAll(ctx context.Context, userID int64) ([]model.ContactModel, error) {
	var contacts []model.ContactModel
	if err := contactRepo.baseRepo.Conn(ctx).Joins("User").Find(&contacts, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return contacts, nil
}

func (contactRepo *contactRepositoryImpl) FindContact(ctx context.Context, userID, contactUserID int64) (*model.ContactModel, error) {
	var contact model.ContactModel

	if err := contactRepo.baseRepo.Conn(ctx).Joins("User").Find(&contact, "user_id = ? AND contact_user_id = ?", userID, contactUserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}

	return &contact, nil
}

func (contactRepo *contactRepositoryImpl) Add(ctx context.Context, contact *model.ContactModel) error {
	if err := contactRepo.baseRepo.Conn(ctx).Create(contact).Error; err != nil {
		return err
	}

	return nil
}
