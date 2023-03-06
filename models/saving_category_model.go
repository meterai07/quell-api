package models

type SavingCategory struct {
	Name string `gorm:"not null" json:"name"`
}
