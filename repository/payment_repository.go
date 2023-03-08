package repository

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"quell-api/entity"
	"strings"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	FindAll() ([]entity.UserTransaction, error)
	FindById(id uint) (entity.UserTransaction, error)
	CreatePayment(payment entity.UserTransaction) error
	UpdatePayment(payment entity.UserTransaction, id uint) error
	DeletePayment(id uint) error
	CreatePaymentMidTrans(payment entity.Payload) (interface{}, error)
}

type paymentRepository struct {
	typePayment string
	db          *gorm.DB
}

func NewPaymentRepository(db *gorm.DB, typePayment string) *paymentRepository {
	return &paymentRepository{
		typePayment: typePayment,
		db:          db,
	}
}

func (r *paymentRepository) FindAll() ([]entity.UserTransaction, error) {
	var payments []entity.UserTransaction
	result := r.db.Find(&payments).Error
	if result != nil {
		return payments, result
	}
	return payments, nil
}

func (r *paymentRepository) FindById(id uint) (entity.UserTransaction, error) {
	var payment entity.UserTransaction
	result := r.db.Where("id = ?", id).First(&payment).Error
	if result != nil {
		return payment, result
	}
	return payment, nil
}

func (r *paymentRepository) CreatePayment(payment entity.UserTransaction) error {
	result := r.db.Create(&payment).Error
	if result != nil {
		return result
	}
	return nil
}

func (r *paymentRepository) UpdatePayment(payment entity.UserTransaction, id uint) error {
	var paymentUpdate entity.UserTransaction
	result := r.db.Where("id = ?", id).First(&paymentUpdate).Error
	if result != nil {
		return result
	}
	paymentUpdate.Status = payment.Status
	result = r.db.Save(&paymentUpdate).Error
	if result != nil {
		return result
	}
	return nil
}

func (r *paymentRepository) DeletePayment(id uint) error {
	var payment entity.UserTransaction
	result := r.db.Where("id = ?", id).First(&payment).Error
	if result != nil {
		return result
	}
	result = r.db.Delete(&payment).Error
	if result != nil {
		return result
	}
	return nil
}

func (r *paymentRepository) CreatePaymentMidTrans(payment entity.Payload) (interface{}, error) {
	data, err := json.Marshal(payment)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(data))

	payload := strings.NewReader(string(data))

	req, err := http.NewRequest("POST", "https://api.sandbox.midtrans.com/v2/charge", payload)
	if err != nil {
		return nil, err
	}

	fmt.Println(req)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(os.Getenv("SERVER_KEY")))))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	fmt.Println(res)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var responseBody any
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}
	return responseBody, nil
}
