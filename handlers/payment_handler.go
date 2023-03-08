package handlers

import (
	"fmt"
	"net/http"
	"quell-api/entity"
	"quell-api/repository"
	"quell-api/sdk/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	paymentService repository.PaymentRepository
}

func NewPaymentHandler(paymentService repository.PaymentRepository) *PaymentHandler {
	return &PaymentHandler{paymentService}
}

func (p *PaymentHandler) PremiumPayment(c *gin.Context) {
	var listData []entity.ItemDetailsContent
	var ItemDetailsContent entity.ItemDetailsContent
	var transactionDetails entity.TransactionDetailsContent
	var customerDetails entity.CustomerDetails
	var gopay entity.Gopay
	var payload entity.Payload

	if err := c.ShouldBindJSON(&ItemDetailsContent); err != nil {
		response.Response(c, http.StatusBadRequest, "Failed to bind json", nil)
		return
	}

	if err := validator.New().Struct(&ItemDetailsContent); err != nil {
		validationError := err.(validator.ValidationErrors)
		response.Response(c, http.StatusBadRequest, validationError.Error(), nil)
		return
	}

	transactionDetails = entity.TransactionDetailsContent{
		Order_ID:     fmt.Sprintf("order %v", uuid.New()),
		Gross_Amount: ItemDetailsContent.Price * ItemDetailsContent.Quantity,
	}

	listData = append(listData, ItemDetailsContent)

	customerDetails = entity.CustomerDetails{
		First_name: c.MustGet("user").(entity.User).Username,
		Last_name:  "user",
		Email:      c.MustGet("user").(entity.User).Email,
		Phone:      "08123456789",
	}

	gopay = entity.Gopay{
		Enable_callback: true,
		Callback_url:    "https://midtrans.com",
	}

	payload = entity.Payload{
		Customer_details:    customerDetails,
		Gopay:               gopay,
		Item_details:        listData,
		Payment_type:        "gopay",
		Transaction_details: transactionDetails,
	}

	result, err := p.paymentService.CreatePayment(payload)
	if err != nil {
		response.Response(c, http.StatusBadRequest, "Failed to create payment", err.Error())
		return
	}

	response.Response(c, http.StatusOK, "success", result)

	// result disimpan didalam database
}
