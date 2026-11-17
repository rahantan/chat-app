package service

import (
	"context"
	"realtime-chat/internal/dto/request"
	"realtime-chat/internal/dto/response"
	"realtime-chat/internal/model"
	"realtime-chat/internal/pkg"
	"realtime-chat/internal/repository/mysql"
)

type ContactService interface {
	// GetContact(ctx context.Context, userID, contactUserID int64) (*response.Contact, error)
	List(ctx context.Context, userID int64) ([]response.Contact, error)
	Add(ctx context.Context, userID int64, req *request.AddContact) error
	Delete(ctx context.Context, userID, contactID int64) error
}

type contactServiceImpl struct {
	mysql.ContactRepository
	mysql.UserRepository
}

func NewContactService(contact mysql.ContactRepository, user mysql.UserRepository) ContactService {
	return &contactServiceImpl{
		ContactRepository: contact,
		UserRepository:    user,
	}
}
func (contactSvc *contactServiceImpl) Delete(ctx context.Context, userID, contactUserID int64) error {
	if err := contactSvc.ContactRepository.Delete(ctx, userID, contactUserID); err != nil {
		return err
	}
	return nil
}

func (contactSvc *contactServiceImpl) Add(ctx context.Context, userID int64, req *request.AddContact) error {
	if userID == req.ContactUserID {
		return pkg.NewErr(pkg.ErrCodeBadRequest, "cannot add yourself as a contact")
	}

	if _, err := contactSvc.UserRepository.FindByID(ctx, req.ContactUserID); err != nil {
		return err
	}

	exist, err := contactSvc.ContactRepository.IsExist(ctx, userID, req.ContactUserID)
	if err != nil {
		return err
	}

	if exist {
		return pkg.NewErr(pkg.ErrCodeConflict, "contact already exists")
	}

	contact := model.ContactModel{
		UserID:        userID,
		ContactUserID: req.ContactUserID,
	}

	if err := contactSvc.ContactRepository.Add(ctx, &contact); err != nil {
		return err
	}

	return nil
}

func (contactSvc *contactServiceImpl) List(ctx context.Context, userID int64) ([]response.Contact, error) {
	contactResults, err := contactSvc.ContactRepository.FindAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	contacts := make([]response.Contact, 0)
	for _, contact := range contactResults {
		contacts = append(contacts, response.Contact{
			ID:            contact.ID,
			ContactUserID: contact.ContactUserID,
			UserName:      contact.User.UserName,
			Avatar:        pkg.AvatarImageUrl + contact.User.Avatar,
		})
	}

	return contacts, nil
}
