package resp

type SessionRecord struct {
	ID       int    `gorm:"primarykey"`
	UUID     string `gorm:"uniqueIndex;varchar(255);not null"`
	Username string `gorm:"uniqueIndex;varchar(255);not null"`
	Email    string `gorm:"uniqueIndex;varchar(255);not null"`
}
