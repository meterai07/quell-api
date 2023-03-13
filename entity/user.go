package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username         string            `gorm:"unique;not null" validate:"required,min=5,max=30"`
	Email            string            `gorm:"unique;not null" validate:"required,email"`
	Password         string            `gorm:"not null" validate:"required,min=8"`
	Phone            string            `gorm:"unique;not null" validate:"required,min=10,max=13"`
	IsActive         bool              `gorm:"default:false"`
	IsPremium        bool              `gorm:"default:false"`
	Token            string            `gorm:"unique;not null"`
	Posts            []Post            `gorm:"foreignKey:UserID"`
	UserTransactions []UserTransaction `gorm:"foreignKey:UserID"`
}

type UserTransaction struct {
	gorm.Model
	GrossAmount int    `gorm:"not null" validate:"required"`
	OrderID     string `gorm:"unique;not null" validate:"required"`
	Status      string `gorm:"not null" validate:"required,oneof=SUCCESS PENDING"`
	Lifetime    int    `gorm:"not null" validate:"required"`
	UserID      uint
}

type UserPremium struct {
	gorm.Model
	StartDate time.Time `gorm:"not null" validate:"required"`
	EndDate   time.Time `gorm:"not null" validate:"required"`
	UserID    uint
}
