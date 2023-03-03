package entity

import "gorm.io/gorm"

type Saving struct {
	gorm.Model
	Name             string `gorm:"not null" validate:"required,min=5,max=30"`
	Description      string `gorm:"not null" validate:"required,min=5,max=1000"`
	Amount           int    `gorm:"not null" validate:"required"`
	SavingCategoryID *uint
	UserID           uint
	Type             string `gorm:"not null" validate:"required, oneof=income expense"`
}

type SavingCategory struct {
	gorm.Model
	Name    string   `gorm:"not null" validate:"required,min=5,max=30"`
	Savings []Saving `gorm:"foreignKey:SavingCategoryID"`
}
