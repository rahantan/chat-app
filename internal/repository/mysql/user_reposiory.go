package mysql

import (
	"context"
	"errors"
	"simple-auth/internal/model"
	"simple-auth/internal/pkg"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(ctx context.Context, id int64) (*model.UserModel, error)
	Create(ctx context.Context, user *model.UserModel) error
	FindByEmail(ctx context.Context, email string) (*model.UserModel, error)
}

type userRepositoryImpl struct {
	*gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{
		DB: db,
	}
}

func (userRepo *userRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.UserModel, error) {
	var user model.UserModel
	err := userRepo.DB.WithContext(ctx).Take(&userRepo, "id=?", id).Error
	return &user, err
}

func (userRepo *userRepositoryImpl) Create(ctx context.Context, user *model.UserModel) error {
	if err := userRepo.DB.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (userRepo *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.UserModel, error) {
	var user model.UserModel
	if err := userRepo.DB.WithContext(ctx).Take(&user, "email=?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.NewError(pkg.KindNotFound, "not found email", nil)
		}
		return nil, err
	}

	return &user, nil
}
