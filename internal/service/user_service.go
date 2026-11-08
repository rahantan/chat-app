package service

import (
	"context"
	"realtime-chat/internal/dto"
	mysqlRepo "realtime-chat/internal/repository/mysql"
)

type UserService interface {
	GetById(ctx context.Context, userId int64) (*dto.UserResponse, error)
}

type userServiceImpl struct {
	mysqlRepo.UserRepository
}

func (userService *userServiceImpl) GetById(ctx context.Context, userId int64) (*dto.UserResponse, error) {
	user, err := userService.UserRepository.FindByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &dto.UserResponse{
		ID:       userId,
		Email:    user.Email,
		UserName: user.UserName,
	}, nil
}
