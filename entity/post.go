package entity

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name  string `gorm:"not null" validate:"required,min=5,max=30"`
	Posts []Post `gorm:"foreignKey:CategoryID"`
}

type Post struct {
	gorm.Model
	Title       string    `gorm:"not null" validate:"required,min=5,max=30"`
	Content     string    `gorm:"not null" validate:"required,min=5,max=1000"`
	Date        time.Time `gorm:"not null" validate:"required" json:"date"`
	Type        string    `gorm:"not null" validate:"required, oneof=jadwal tugas"`
	UserID      uint
	CategoryID  *uint
	Attachments []Attachment `gorm:"foreignKey:PostID"`
}

type Attachment struct {
	gorm.Model
	Name   string `gorm:"not null" validate:"required"`
	UserID uint
	PostID uint
	Url    string `gorm:"not null" validate:"required"`
}
