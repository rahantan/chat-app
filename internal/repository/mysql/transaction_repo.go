package mysql

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

type TransactorRepository interface {
	WithInTransaction(ctx context.Context, fn func(ctxTx context.Context) error) error
}

type gormConn struct {
	*gorm.DB
}

func NewTransactorRepository(db *gorm.DB) TransactorRepository {
	return &gormConn{
		DB: db,
	}
}

func (db *gormConn) WithInTransaction(ctx context.Context, fn func(c context.Context) error) error {
	if _, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return fn(ctx)
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, txKey{}, tx))
	})
}
