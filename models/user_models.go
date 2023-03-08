package models

type User_Register struct {
	Username string `json:"username" validate:"required,min=5,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,min=10,max=13"`
	Password string `json:"password" validate:"required,min=8"`
}

type User_Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
