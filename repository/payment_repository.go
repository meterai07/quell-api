package repository

type PaymentRepository interface {
	// CreatePayment(payment interface{}) error
}

type paymentRepository struct {
	typePayment string
}

func NewPaymentRepository() *paymentRepository {
	return &paymentRepository{}
}
