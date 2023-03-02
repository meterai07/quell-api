package models

type Post_Upload struct {
	Title      string `json:"title" validate:"required,min=5,max=30"`
	Content    string `json:"content" validate:"required,min=5,max=1000"`
	Date       string `json:"date" validate:"required"`
	Type       string `json:"type" validate:"required, oneof=jadwal tugas"`
	CategoryID *uint  `json:"category_id" validate:"required"`
}
