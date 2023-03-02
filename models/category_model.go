package models

type Category struct {
	Name string `json:"name" validate:"required,min=5,max=30"`
}
