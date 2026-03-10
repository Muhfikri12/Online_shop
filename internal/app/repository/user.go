package repository

import (
	"app/internal/model"
	"context"

	"gorm.io/gorm"
)

type RUser interface {
	FindByUsernameOrEmail(ctx context.Context, username string, email string) (model.User, error)
	FindByID(ctx context.Context, id int) (model.User, error)
	Create(ctx context.Context, user model.User) error
}

type rUser struct {
	db *gorm.DB
}

func NewRUser(db *gorm.DB) RUser {
	return &rUser{db: db}
}

func (r *rUser) FindByUsernameOrEmail(ctx context.Context, username string, email string) (model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("username = ? OR email = ?", username, email).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (r *rUser) FindByID(ctx context.Context, id int) (model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (r *rUser) Create(ctx context.Context, user model.User) error {
	return r.db.WithContext(ctx).Create(&user).Error
}
