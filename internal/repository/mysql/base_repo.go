package mysql

import (
	"context"

	"gorm.io/gorm"
)

type BaseRepsoitory interface {
	Conn(ctx context.Context) *gorm.DB
}
type baseRepoImpl struct {
	*gorm.DB
}

func NewBaseRepository(db *gorm.DB) BaseRepsoitory {
	return &baseRepoImpl{
		DB: db,
	}
}

func (b *baseRepoImpl) Conn(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}

	return b.DB.WithContext(ctx)
}
