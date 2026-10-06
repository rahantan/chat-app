package service

import (
	"context"
	"simple-auth/internal/dto"
	"simple-auth/internal/model"
	"simple-auth/internal/pkg"
	mysqlRepo "simple-auth/internal/repository/mysql"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) error
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.UserResponse, error)
}

type authServiceImpl struct {
	mysqlRepo.UserRepository
}

func NewAuthService(userRepo mysqlRepo.UserRepository) AuthService {
	return &authServiceImpl{
		UserRepository: userRepo,
	}
}

func (authService *authServiceImpl) Register(ctx context.Context, req *dto.RegisterRequest) error {

	pass, err := bcrypt.GenerateFromPassword([]byte(req.Pass), bcrypt.DefaultCost)
	if err != nil {
		return pkg.NewError(pkg.KindInternal, err.Error(), nil)
	}

	err = authService.UserRepository.Create(ctx, &model.UserModel{
		UserName: req.UserName,
		Email:    req.Email,
		Password: string(pass),
	})

	return err
}

func (authService *authServiceImpl) Login(ctx context.Context, req *dto.LoginRequest) (*dto.UserResponse, error) {
	user, err := authService.UserRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, pkg.ErrLogin
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Pass)); err != nil {
		return nil, pkg.ErrLogin
	}

	return &dto.UserResponse{
		ID:       user.ID,
		UserName: user.UserName,
		Email:    user.Email,
	}, nil
}
