package resp

import "time"

type RespProduct struct {
	ID           uint      `json:"-"`
	UUID         string    `json:"uuid"`
	SKU          string    `json:"sku"`
	Name         string    `json:"name"`
	CategoryID   int       `json:"-"`
	CategoryUUID string    `json:"category_uuid"`
	CategoryName string    `json:"category_name"`
	Price        int       `json:"price"`
	Stock        int       `json:"stock"`
	Status       string    `json:"status"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
