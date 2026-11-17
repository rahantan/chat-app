package mysql

import (
	"context"
	"errors"
	"realtime-chat/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(ctx context.Context, userID int64) (*model.UserModel, error)
	Create(ctx context.Context, user *model.UserModel) error
	FindByEmail(ctx context.Context, email string) (*model.UserModel, error)
	FindByIDs(ctx context.Context, userIDs []int64) ([]model.UserModel, error)
}

type userRepositoryImpl struct {
	baseRepo BaseRepsoitory
}

func NewUserRepository(repo BaseRepsoitory) UserRepository {
	return &userRepositoryImpl{
		baseRepo: repo,
	}
}

func (userRepo *userRepositoryImpl) FindByIDs(ctx context.Context, userIDs []int64) ([]model.UserModel, error) {
	users := make([]model.UserModel, 0, len(userIDs))

	if err := userRepo.baseRepo.Conn(ctx).Find(&users, "id IN ?", userIDs).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (userRepo *userRepositoryImpl) FindByID(ctx context.Context, userID int64) (*model.UserModel, error) {
	var user model.UserModel
	if err := userRepo.baseRepo.Conn(ctx).Take(&user, "id=?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (userRepo *userRepositoryImpl) Create(ctx context.Context, user *model.UserModel) error {
	if err := userRepo.baseRepo.Conn(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (userRepo *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.UserModel, error) {
	var user model.UserModel
	if err := userRepo.baseRepo.Conn(ctx).Take(&user, "email=?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
