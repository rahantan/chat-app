package mysql

import (
	"context"
	"errors"
	"realtime-chat/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RefreshTokenRepository interface {
	DeleteByUserID(ctx context.Context, userID int64) error
	UpSertByUserID(ctx context.Context, token *model.RefreshToken) error
	ExsistByTokenHash(ctx context.Context, tokenHash string) (bool, error)
}

type refreshTokenRepositoryImpl struct {
	baseRepo BaseRepsoitory
}

func NewRefreshTokenRepository(baseRepo BaseRepsoitory) RefreshTokenRepository {
	return &refreshTokenRepositoryImpl{
		baseRepo: baseRepo,
	}
}

func (tokenRepo *refreshTokenRepositoryImpl) ExsistByTokenHash(ctx context.Context, tokenHash string) (bool, error) {
	if err := tokenRepo.baseRepo.Conn(ctx).Take(&model.RefreshToken{}, "token_hash = ?", tokenHash).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
func (tokenRepo *refreshTokenRepositoryImpl) DeleteByUserID(ctx context.Context, userID int64) error {
	result := tokenRepo.baseRepo.Conn(ctx).Delete(&model.RefreshToken{}, "user_id = ?", userID)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
func (tokenRepo *refreshTokenRepositoryImpl) UpSertByUserID(ctx context.Context, token *model.RefreshToken) error {
	onConflict := clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"token_hash"}),
	}

	if err := tokenRepo.baseRepo.Conn(ctx).Clauses(onConflict).Create(token).Error; err != nil {
		return err
	}
	return nil
}
