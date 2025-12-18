package repository

import (
	resp "app/internal/dto/response"
	"context"

	"gorm.io/gorm"
)

/* --------------------------------- Interface --------------------------------- */
type RProduct interface {
	FindByUUID(ctx context.Context, uuid string) (*resp.RespProduct, error)
}

type rProduct struct {
	db *gorm.DB
}

func NewRProduct(db *gorm.DB) RProduct {
	return &rProduct{db}
}

func (r *rProduct) FindByUUID(ctx context.Context, uuid string) (*resp.RespProduct, error) {
	var product resp.RespProduct
	if err := r.db.WithContext(ctx).Table("products p").
		Select("p.id, p.uuid, p.sku, p.name, p.category_id, c.uuid as category_uuid, c.name as category_name, p.price, p.stock, p.status, p.description, p.created_at, p.updated_at").
		Joins("LEFT JOIN categories c ON p.category_id = c.id").
		Where("p.uuid = ?", uuid).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}
