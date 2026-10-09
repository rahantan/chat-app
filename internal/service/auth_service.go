package service

import (
	"context"
	"errors"
	"simple-auth/config"
	"simple-auth/internal/dto"
	"simple-auth/internal/model"
	"simple-auth/internal/pkg"
	mysqlRepo "simple-auth/internal/repository/mysql"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) error
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
	LogOut(ctx context.Context, userID int64) error
	RefreshToken(ctx context.Context, jwtToken string) (*dto.AuthResponse, error)
	VerifyToken(secretKey string, tokenStr string) (*dto.UserResponse, error)
}

type authServiceImpl struct {
	mysqlRepo.UserRepository
	mysqlRepo.RefreshTokenRepository
	*config.Jwt
}

func NewAuthService(userRepo mysqlRepo.UserRepository, rToken mysqlRepo.RefreshTokenRepository, jwt *config.Jwt) AuthService {
	return &authServiceImpl{
		UserRepository:         userRepo,
		RefreshTokenRepository: rToken,
		Jwt:                    jwt,
	}
}

func (authService *authServiceImpl) Register(ctx context.Context, req *dto.RegisterRequest) error {

	pass, err := bcrypt.GenerateFromPassword([]byte(req.Pass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = authService.UserRepository.Create(ctx, &model.UserModel{
		UserName: req.UserName,
		Email:    req.Email,
		Password: string(pass),
	})

	return err
}

func (authService *authServiceImpl) LogOut(ctx context.Context, userID int64) error {
	return authService.RefreshTokenRepository.DeleteByUserID(ctx, userID)
}

func (authService *authServiceImpl) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {

	result, err := authService.UserRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pkg.ErrUserNotFound) {
			return nil, pkg.ErrLogin
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(req.Pass)); err != nil {
		return nil, pkg.ErrLogin
	}

	user := dto.UserResponse{
		ID:       result.ID,
		UserName: result.UserName,
		Email:    result.Email,
	}

	refreshToken, err := authService.createToken(authService.Jwt.REFRESH_TOKEN_SECRET_KEY, pkg.RefreshTokenExpireTime, &user)
	if err != nil {
		return nil, err
	}

	accessToken, err := authService.createToken(authService.Jwt.ACCESS_TOKEN_SECRET_KEY, pkg.AccessTokenExpireTime, &user)
	if err != nil {
		return nil, err
	}

	if err := authService.RefreshTokenRepository.UpSertByUserID(ctx, &model.RefreshToken{UserID: user.ID, TokenHash: pkg.HashToken(refreshToken)}); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		UserResponse: user,
	}, nil
}

func (authService *authServiceImpl) RefreshToken(ctx context.Context, token string) (*dto.AuthResponse, error) {
	userClaim, err := authService.VerifyToken(authService.Jwt.REFRESH_TOKEN_SECRET_KEY, token)
	if err != nil {
		return nil, err
	}

	exists, err := authService.RefreshTokenRepository.ExsistByTokenHash(ctx, pkg.HashToken(token))
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, pkg.NewError(pkg.KindUnauthorized, "token has been revoked", nil)
	}

	newRefreshToken, err := authService.createToken(authService.REFRESH_TOKEN_SECRET_KEY, pkg.RefreshTokenExpireTime, userClaim)
	if err != nil {
		return nil, err
	}

	newAccessToken, err := authService.createToken(authService.Jwt.ACCESS_TOKEN_SECRET_KEY, pkg.AccessTokenExpireTime, userClaim)
	if err != nil {
		return nil, err
	}

	if err := authService.RefreshTokenRepository.UpSertByUserID(ctx, &model.RefreshToken{UserID: userClaim.ID, TokenHash: pkg.HashToken(newRefreshToken)}); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		RefreshToken: newRefreshToken,
		AccessToken:  newAccessToken,
		UserResponse: *userClaim,
	}, nil
}

func (authService *authServiceImpl) VerifyToken(secretKey string, tokenStr string) (*dto.UserResponse, error) {
	claims := &dto.JwtPayLoad{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, pkg.NewError(pkg.KindUnauthorized, "invalid signature algorithm", nil)
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, pkg.NewError(pkg.KindUnauthorized, "token is expired", nil)
	}

	if !token.Valid {
		return nil, pkg.NewError(pkg.KindUnauthorized, "invalid token", nil)
	}

	return &claims.UserResponse, nil
}

func (authService *authServiceImpl) createToken(secretKey string, expireTime time.Duration, user *dto.UserResponse) (string, error) {
	start := time.Now()
	payload := dto.JwtPayLoad{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(start.Add(expireTime)),
		},
		UserResponse: *user,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte(secretKey))
	return token, err
}
