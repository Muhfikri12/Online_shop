package repository

import (
	"app/internal/model"
	"context"
	"time"

	"gorm.io/gorm"
)

// RAuth handles persistence of refresh tokens.
type RAuth interface {
	CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error
	FindByTokenID(ctx context.Context, tokenID string) (model.RefreshToken, error)
	RevokeToken(ctx context.Context, id int64, replacedByTokenID *int64) error
	RevokeAllByUserID(ctx context.Context, userID int) error
}

type rAuth struct {
	db *gorm.DB
}

func NewRAuth(db *gorm.DB) RAuth {
	return &rAuth{db: db}
}

func (r *rAuth) CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *rAuth) FindByTokenID(ctx context.Context, tokenID string) (model.RefreshToken, error) {
	var t model.RefreshToken
	err := r.db.WithContext(ctx).Where("token_id = ?", tokenID).First(&t).Error
	return t, err
}

func (r *rAuth) RevokeToken(ctx context.Context, id int64, replacedByTokenID *int64) error {
	now := time.Now()
	update := map[string]interface{}{
		"revoked":    true,
		"revoked_at": &now,
	}
	if replacedByTokenID != nil {
		update["replaced_by_token_id"] = replacedByTokenID
	}
	return r.db.WithContext(ctx).Model(&model.RefreshToken{}).Where("id = ?", id).Updates(update).Error
}

func (r *rAuth) RevokeAllByUserID(ctx context.Context, userID int) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("user_id = ? AND revoked = FALSE", userID).
		Updates(map[string]interface{}{
			"revoked":    true,
			"revoked_at": &now,
		}).Error
}

