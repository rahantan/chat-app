package mysql

import (
	"context"
	"simple-auth/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RefreshTokenRepository interface {
	ExsistByTokenHash(ctx context.Context, tokenHash string) (bool, error)
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

func (rToken *refreshTokenImpl) ExsistByTokenHash(ctx context.Context, tokenHash string) (bool, error) {
	var count int64
	err := rToken.DB.Model(&model.RefreshToken{}).WithContext(ctx).Where("token_hash = ?", tokenHash).Count(&count).Error
	return count > 0, err
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
