package mysql

import (
	"context"
	"realtime-chat/internal/model"
)

type ContactRepository interface {
	FindAll(ctx context.Context, userID int64) ([]model.ContactModel, error)
	FindContact(ctx context.Context, userID, contactUserID int64) (*model.ContactModel, error)
	Add(ctx context.Context, contact *model.ContactModel) error
}

type contactRepositoryImpl struct {
}
