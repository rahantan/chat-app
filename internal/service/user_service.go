package service

import (
	"context"
	"realtime-chat/internal/dto/response"
	"realtime-chat/internal/pkg"
	mysqlRepo "realtime-chat/internal/repository/mysql"
)

type UserService interface {
	GetById(ctx context.Context, userId int64) (*response.User, error)
}

type userServiceImpl struct {
	mysqlRepo.UserRepository
}

func NewUserService(user mysqlRepo.UserRepository) UserService {
	return &userServiceImpl{
		UserRepository: user,
	}
}

func (userService *userServiceImpl) GetById(ctx context.Context, userId int64) (*response.User, error) {
	user, err := userService.UserRepository.FindByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &response.User{
		ID:        userId,
		Email:     user.Email,
		UserName:  user.UserName,
		AvatarUrl: pkg.AvatarImageUrl + user.Avatar,
	}, nil
}
