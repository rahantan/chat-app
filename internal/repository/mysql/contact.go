package mysql

import (
	"context"
	"realtime-chat/internal/model"
)

type ContactRepository interface {
	// FindAll(ctx context.Context, userID int64) ([]model.ContactResponse, error)
	// FindContact(ctx context.Context, userID, contactUserID int64) (*model.ContactModel, error)
	IsExist(ctx context.Context, userID, contactUserID int64) (bool, error)
	FindAll(ctx context.Context, userID int64) ([]model.ContactModel, error)
	Delete(ctx context.Context, userID, contactUserID int64) error
	Add(ctx context.Context, contact *model.ContactModel) error
}

type contactRepositoryImpl struct {
	baseRepo BaseRepsoitory
}

func NewContactRepository(base BaseRepsoitory) ContactRepository {
	return &contactRepositoryImpl{
		baseRepo: base,
	}
}

func (contactRepo *contactRepositoryImpl) IsExist(ctx context.Context, userID, contactUserID int64) (bool, error) {
	var count int64
	if err := contactRepo.baseRepo.Conn(ctx).Model(&model.ContactModel{}).Where("user_id =? AND contact_user_id = ?", userID, contactUserID).Count(&count).Error; err != nil {
		return false, err
	}
	return count == 1, nil
}

func (contactRepo *contactRepositoryImpl) FindAll(ctx context.Context, userID int64) ([]model.ContactModel, error) {
	var contacts []model.ContactModel
	if err := contactRepo.baseRepo.Conn(ctx).Joins("User").Find(&contacts, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}

	return contacts, nil
}

func (contactRepo *contactRepositoryImpl) Add(ctx context.Context, contact *model.ContactModel) error {
	if err := contactRepo.baseRepo.Conn(ctx).Create(contact).Error; err != nil {
		return err
	}

	return nil
}

func (contactRepo *contactRepositoryImpl) Delete(ctx context.Context, userID, contactUserID int64) error {
	result := contactRepo.baseRepo.Conn(ctx).Delete(&model.ContactModel{}, "user_id = ? AND id = ?", userID, contactUserID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return model.ErrContactNotFound
	}
	return nil
}
