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
)

type PaymentRepository interface {
	// FindAll() ([]entity.PaymentHistory, error)
	// FindById(id uint) (entity.PaymentHistory, error)
	// CreatePayment(payment entity.PaymentHistory) error
	// UpdatePayment(payment entity.PaymentHistory, id uint) error
	// DeletePayment(id uint) error
	CreatePayment(payment entity.Payload) (interface{}, error)
}

type paymentRepository struct {
	typePayment string
}

func NewPaymentRepository(typePayment string) *paymentRepository {
	return &paymentRepository{
		typePayment: typePayment,
	}
}

func (r *paymentRepository) CreatePayment(payment entity.Payload) (interface{}, error) {
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
