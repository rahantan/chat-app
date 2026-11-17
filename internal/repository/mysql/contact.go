package mysql

import (
	"context"
	"errors"
	"realtime-chat/internal/model"

	"gorm.io/gorm"
)

type ContactRepository interface {
	// FindAll(ctx context.Context, userID int64) ([]model.ContactResponse, error)
	FindAll(ctx context.Context, userID int64) ([]model.ContactModel, error)
	FindContact(ctx context.Context, userID, contactUserID int64) (*model.ContactModel, error)
	Delete(ctx context.Context, userID, contactUserID int64) error
	IsExist(ctx context.Context, userID, contactUserID int64) (bool, error)
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

// func (contactRepo *contactRepositoryImpl) FindAll(ctx context.Context, userID int64) ([]model.ContactResponse, error) {
// 	var contacts []model.ContactResponse

// 	err := contactRepo.baseRepo.Conn(ctx).Table("contacts c").Select(`
// 		c.id,
// 		c.contact_user_id,
// 		u.username,
// 		u.avatar,
// 		rm.room_id
// 	`).
// 		Joins("JOIN users u ON u.id = c.contact_user_id").
// 		Joins("LEFT JOIN room_members rm ON rm.user_id= c.contact_user_id").
// 		Joins("LEFT JOIN room_members me ON me.room_id=rm.room_id AND me.user_id= c.user_id").
// 		Joins("LEFT JOIN rooms r ON r.id = rm.room_id AND r.type = ?", "private").
// 		Where("c.user_id = ?", userID).
// 		Where("(me.user_id IS NOT NULL AND r.id IS NOT NULL) OR rm.id IS NULL").
// 		Scan(&contacts).Error

// 	if err != nil {
// 		return nil, err
// 	}
// 	return contacts, nil
// }

// func (contactRepo *contactRepositoryImpl) FindAll(ctx context.Context, userID int64) ([]model.ContactModel, error) {
// 	var contacts []model.ContactModel
// 	if err := contactRepo.baseRepo.Conn(ctx).Joins("User").Find(&contacts, "user_id = ?", userID).Error; err != nil {
// 		return nil, err
// 	}
// 	return contacts, nil
// }

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
