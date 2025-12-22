package model

type User struct {
	ID       int    `gorm:"primarykey"`
	UUID     string `gorm:"uniqueIndex;varchar(255);not null"`
	Username string `gorm:"uniqueIndex;varchar(255);not null"`
	Password string `gorm:"varchar(255);not null"`
	Email    string `gorm:"uniqueIndex;varchar(255);not null"`
	TimeStamp
}
