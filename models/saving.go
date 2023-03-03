package models

type Saving struct {
	Name             string `json:"name" validate:"required,min=5,max=30"`
	Description      string `json:"description" validate:"required,min=5,max=1000"`
	Amount           int    `json:"amount" validate:"required"`
	Type             string `json:"type" validate:"required, oneof=income expense"`
	SavingCategoryID *uint  `json:"category_id" validate:"required"`
}
