package repository

import (
	"context"
	"errors"

	userDomain "erp-api/internal/domain/user"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) userDomain.Repository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *userDomain.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return userDomain.ErrUserAlreadyExists
		}
		return result.Error
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*userDomain.User, error) {
	var user userDomain.User
	
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, userDomain.ErrUserNotFound
		}
		return nil, result.Error
	}
	
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *userDomain.User) error {
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return userDomain.ErrUserNotFound
	}
	
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&userDomain.User{})
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return userDomain.ErrUserNotFound
	}
	
	return nil
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*userDomain.User, error) {
	var users []*userDomain.User
	
	result := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return users, nil
}

func (r *UserRepository) Count(ctx context.Context) (int, error) {
	var count int64
	
	result := r.db.WithContext(ctx).Model(&userDomain.User{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	
	return int(count), nil
}