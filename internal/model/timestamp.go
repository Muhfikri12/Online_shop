package model

import (
	"time"

	"gorm.io/gorm"
)

type TimeStamp struct {
	CreatedAt time.Time      `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamp"`
}
