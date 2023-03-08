package handlers

import (
	"net/http"
	"quell-api/entity"
	"quell-api/sdk/response"
	"quell-api/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
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
		Order_ID:     uuid.New().String(),
		Gross_Amount: ItemDetailsContent.Price * ItemDetailsContent.Quantity,
	}

	listData = append(listData, ItemDetailsContent)

	customerDetails = entity.CustomerDetails{
		First_name: c.MustGet("user").(entity.User).Username,
		Last_name:  "user",
		Email:      c.MustGet("user").(entity.User).Email,
		Phone:      c.MustGet("user").(entity.User).Phone,
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

	result, err := p.paymentService.CreatePaymentMidTrans(payload)
	if err != nil {
		response.Response(c, http.StatusBadRequest, "Failed to create payment", err.Error())
		return
	}

	userTransaction := entity.UserTransaction{
		GrossAmount: transactionDetails.Gross_Amount,
		OrderID:     transactionDetails.Order_ID,
		Status:      "PENDING",
		UserID:      c.MustGet("user").(entity.User).ID,
	}

	err = p.paymentService.CreatePayment(userTransaction)
	if err != nil {
		response.Response(c, http.StatusBadRequest, "Failed to create user transaction", err.Error())
		return
	}

	response.Response(c, http.StatusOK, "success", result)
}

func (p *PaymentHandler) PremiumPaymentValidate(c *gin.Context) {
	// endpoint ketika transaksi telah dibayar
	var validatePayment entity.ValidatePayment

	if err := c.ShouldBindJSON(&validatePayment); err != nil {
		response.Response(c, http.StatusBadRequest, "Failed to bind json", nil)
		return
	}

	if err := validator.New().Struct(&validatePayment); err != nil {
		validationError := err.(validator.ValidationErrors)
		response.Response(c, http.StatusBadRequest, validationError.Error(), nil)
		return
	}
}
