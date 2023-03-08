package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" validate:"required,min=5,max=30"`
	Email    string `gorm:"unique;not null" validate:"required,email"`
	Password string `gorm:"not null" validate:"required,min=8"`
	IsActive bool   `gorm:"default:false"`
	Token    string `gorm:"unique;not null"`
	Posts    []Post `gorm:"foreignKey:UserID"`
}

// type UserTransaction struct {
// 	gorm.Model
// 	UserID uint
// 	Amount int
// }
