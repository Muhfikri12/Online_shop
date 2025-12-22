package seeder

import (
	"app/internal/model"
	"app/pkg/toolkit"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB) error {
	if err := SeedUser(db); err != nil {
		return err
	}

	return nil
}

func SeedUser(db *gorm.DB) error {

	hashPassword, err := toolkit.HashPassword("admin")
	if err != nil {
		return err
	}

	user := model.User{
		ID:       1,
		UUID:     uuid.New().String(),
		Username: "admin",
		Password: hashPassword,
		Email:    "admin@admin.com",
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}
