package user

import "gorm.io/gorm"

type UserEntity struct {
	gorm.Model
	Name    string `gorm:"type:varchar(100);not null"`
	Email   string `gorm:"uniqueIndex;not null"`
	VelogId string `gorm:"type:varchar(100);not null"`
	Token   string `gorm:"type:varchar(512);not null;uniqueIndex"`
}
