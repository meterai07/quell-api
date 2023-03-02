package service

import "quell-api/repository"

type PaymentService interface {
	// CreatePayment(payment interface{}) error
}

type paymentService struct {
	repository repository.PaymentRepository
}

func NewPaymentService(repository repository.PaymentRepository) PaymentService {
	return &paymentService{repository}
}
