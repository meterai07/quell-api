package service

import (
	"quell-api/entity"
	"quell-api/repository"
)

type PaymentService interface {
	FindAll(id uint) ([]entity.UserTransaction, error)
	FindById(id string) (entity.UserTransaction, error)
	CreatePayment(payment entity.UserTransaction) error
	UpdatePayment(payment entity.UserTransaction, id uint) error
	DeletePayment(id uint) error
	CreatePaymentMidTrans(payment entity.Payload) (interface{}, error)
}

type paymentService struct {
	repository repository.PaymentRepository
}

func NewPaymentService(repository repository.PaymentRepository) PaymentService {
	return &paymentService{repository}
}

func (s *paymentService) FindAll(id uint) ([]entity.UserTransaction, error) {
	return s.repository.FindAll(id)
}

func (s *paymentService) FindById(id string) (entity.UserTransaction, error) {
	return s.repository.FindById(id)
}

func (s *paymentService) CreatePayment(payment entity.UserTransaction) error {
	return s.repository.CreatePayment(payment)
}

func (s *paymentService) UpdatePayment(payment entity.UserTransaction, id uint) error {
	return s.repository.UpdatePayment(payment, id)
}

func (s *paymentService) DeletePayment(id uint) error {
	return s.repository.DeletePayment(id)
}

func (s *paymentService) CreatePaymentMidTrans(payment entity.Payload) (interface{}, error) {
	return s.repository.CreatePaymentMidTrans(payment)
}
