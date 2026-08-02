package repository

import (
	"context"

	"gorm.io/gorm"
)

// Reader defines read-only operations on a resource.
type Reader[T any] interface {
	FindByID(ctx context.Context, id uint) (*T, error)
}

// Writer defines write operations on a resource.
type Writer[T any] interface {
	Create(ctx context.Context, entity *T) error
	CreateBatch(ctx context.Context, entities []*T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error
	DeleteBatch(ctx context.Context, ids []uint) error
}

// Repository combines Reader and Writer (ISP: compose only what you need).
type Repository[T any] interface {
	Reader[T]
	Writer[T]
}

// GormRepository is a generic GORM implementation of Repository[T].
// Concrete repositories embed this and extend it with domain-specific queries.
type GormRepository[T any] struct {
	DB *gorm.DB
}

func NewGormRepository[T any](db *gorm.DB) *GormRepository[T] {
	return &GormRepository[T]{DB: db}
}

func (r *GormRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Create(entity).Error
}

func (r *GormRepository[T]) CreateBatch(ctx context.Context, entities []*T) error {
	if len(entities) == 0 {
		return nil
	}
	return r.DB.WithContext(ctx).Create(entities).Error
}

func (r *GormRepository[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	if err := r.DB.WithContext(ctx).First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *GormRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

func (r *GormRepository[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	return r.DB.WithContext(ctx).Delete(&entity, id).Error
}

func (r *GormRepository[T]) DeleteBatch(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	var entity T
	return r.DB.WithContext(ctx).Delete(&entity, ids).Error
}
