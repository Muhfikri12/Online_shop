package model

type Category struct {
	ID   uint   `gorm:"primarykey"`
	UUID string `gorm:"type:uuid;varchar(255);uniqueIndex"`
	Name string `gorm:"type:varchar(255);not null"`
	TimeStamp
}
