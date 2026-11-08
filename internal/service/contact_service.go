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
	GetByContactUserID(ctx context.Context, contactUserID int64) (*response.Contact, error)
	List(ctx context.Context) ([]response.Contact, error)
	Add(ctx context.Context, userID int64, req request.AddContact) error
}

type contactServiceImpl struct {
	mysql.ContactRepository
	mysql.UserRepository
	mysql.RoomRepository
	tx mysql.TransactorRepository
}

func (contactSvc *contactServiceImpl) Add(ctx context.Context, userID int64, req request.AddContact) error {
	if _, err := contactSvc.UserRepository.FindByID(ctx, req.ContactUserID); err != nil {
		return err
	}

	return contactSvc.tx.WithInTransaction(ctx, func(ctxTx context.Context) error {

		contact := model.ContactModel{
			UserID:        userID,
			ContactUserID: req.ContactUserID,
		}

		if err := contactSvc.ContactRepository.Add(ctxTx, &contact); err != nil {
			return err
		}

		members := []model.RoomMemberModel{
			model.RoomMemberModel{UserID: userID},
			model.RoomMemberModel{UserID: req.ContactUserID},
		}

		room := model.RoomModel{
			Type:    pkg.RoomTypePrivate,
			Members: members,
		}

		if err := contactSvc.RoomRepository.Create(ctxTx, &room); err != nil {
			return err
		}

		return nil
	})
}

func (contactSvc *contactServiceImpl) GetContact(ctx context.Context, userID, contactUserID int64) (*response.Contact, error) {
	contactResult, err := contactSvc.ContactRepository.FindContact(ctx, userID, contactUserID)
	if err != nil {
		return nil, err
	}

	return &response.Contact{
		ID:            contactResult.ID,
		ContactUserID: contactResult.ContactUserID,
		UserName:      contactResult.User.UserName,
		ImagePath:     pkg.ImagePath + contactResult.User.Image,
	}, nil
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
			ImagePath:     pkg.ImagePath + contact.User.Image,
		})
	}
	return contacts, nil
}
