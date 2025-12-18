package model

type Product struct {
	ID          uint   `gorm:"primarykey"`
	UUID        string `gorm:"type:uuid;varchar(255);uniqueIndex"`
	SKU         string `gorm:"type:varchar(255);not null;uniqueIndex"`
	Name        string `gorm:"type:varchar(255);not null"`
	CategoryID  int    `gorm:"type:int"`
	Price       int    `gorm:"type:int;not null"`
	Stock       int    `gorm:"type:int;not null"`
	Status      string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text;not null"`
	TimeStamp
}
