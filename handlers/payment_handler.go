package handlers

import (
	"net/http"
	"quell-api/entity"
	"quell-api/repository"
	"quell-api/sdk/response"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService repository.PaymentRepository
}

func NewPaymentHandler(paymentService repository.PaymentRepository) *PaymentHandler {
	return &PaymentHandler{paymentService}
}

func PremiumPayment(c *gin.Context) {
	var transactionDetails entity.TransactionDetails

	if err := c.ShouldBindJSON(&transactionDetails); err != nil {
		response.Response(c, http.StatusBadRequest, "Failed to bind json", nil)
	}

}

func (p *PaymentHandler) NewPaymentHandler(c *gin.Context) {
	// var payload entity.Payload

	// paymentType := c.Query("paymentType")
	// customer := entity.Customer{
	// 	Email: c.MustGet("user").(entity.User).Email,
	// 	Name:  c.MustGet("user").(entity.User).Username,
	// }

	// var transaction entity.Transaction

	// gopayContent := map[string]any{}
	// shopeepayContent := map[string]any{}
	// if paymentType == "gopay" {
	// 	gopayContent["enable_callback"] = true
	// 	gopayContent["callback_url"] = "https://midtrans.com"
	// } else if paymentType == "shopeepay" {
	// 	shopeepayContent["enable_callback"] = true
	// 	shopeepayContent["callback_url"] = "https://midtrans.com"
	// }

	// payload = entity.Payload{
	// 	PaymentType: paymentType,
	// 	TransactionDetails: ,
	// 	Itemdetails: ,
	// 	CustomerDetails: customer,

	// result, err := p.paymentService.CreatePayment(payload)
	// if err != nil {
	// 	response.Response(c, http.StatusBadRequest, "Failed to create payment", nil)
	// }

	// response.Response(c, http.StatusOK, "success", result)
}
