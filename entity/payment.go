package entity

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	Gross_amount int
}

type TransactionDetails struct {
	gorm.Model
	Price    int
	Quantity int
	Name     string
}
type Payment struct {
	Enable_callback bool
	Callback_url    string
}

type Customer struct {
	gorm.Model
	Email string
	Name  string
}

type Payload struct {
	gorm.Model
	Paymenttype        string
	Transaction        Transaction
	Transactiondetails TransactionDetails
	Customerdetails    Customer
}
