package mysql

import (
	"context"
	"errors"
	"simple-auth/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RefreshTokenRepository interface {
	FindByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	DeleteByUserID(ctx context.Context, userID int64) error
	UpSertByUserID(ctx context.Context, token *model.RefreshToken) error
}

type refreshTokenImpl struct {
	*gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenImpl{
		DB: db,
	}
}

func (rToken *refreshTokenImpl) FindByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	var tokenResult model.RefreshToken
	if err := rToken.DB.WithContext(ctx).Take(&tokenResult, "token_hash = ?", tokenHash).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid token")
		}
		return nil, err
	}
	return &tokenResult, nil
}
func (rToken *refreshTokenImpl) DeleteByUserID(ctx context.Context, userID int64) error {
	result := rToken.DB.WithContext(ctx).Delete(&model.RefreshToken{}, "user_id = ?", userID)
	if result.RowsAffected < 0 {
		return nil
	}
	return result.Error
}

func (rToken *refreshTokenImpl) UpSertByUserID(ctx context.Context, token *model.RefreshToken) error {
	result := rToken.DB.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"token_hash"}),
		}).Create(token)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
