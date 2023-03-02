package entity

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	Gross_amount int
}

type Itemdetails struct {
	gorm.Model
	Price    int
	Quantity int
	Name     string
}

type Payment struct {
	Enable_callback bool
	Callback_url    string
}
