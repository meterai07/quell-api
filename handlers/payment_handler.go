package handlers

import (
	"net/http"
	"os"
	"quell-api/entity"
	"quell-api/sdk/crypto"
	"quell-api/sdk/response"
	"quell-api/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	paymentService service.PaymentService
	userService    service.Service
}

func NewPaymentHandler(paymentService service.PaymentService, userService service.Service) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		userService:    userService,
	}
}

func (p *PaymentHandler) PremiumPayment(c *gin.Context) {
	month, err := strconv.ParseInt(c.Query("month"), 10, 64)
	if err != nil {
		response.Response(c, http.StatusBadRequest, "Failed to parse month", err.Error())
		return
	}

	var listData []entity.ItemDetailsContent
	var ItemDetailsContent entity.ItemDetailsContent
	var transactionDetails entity.TransactionDetailsContent
	var customerDetails entity.CustomerDetails
	var gopay entity.Gopay
	var payload entity.Payload

	ItemDetailsContent = entity.ItemDetailsContent{
		ID:       uuid.New().String(),
		Name:     "Premium",
		Price:    15000,
		Quantity: int(month),
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
		Lifetime:    int(month),
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

	result, err := p.paymentService.FindById(validatePayment.OrderID)
	if err != nil {
		response.Response(c, http.StatusBadRequest, "Failed to find order id", err.Error())
		return
	}

	makeSignatureKey := validatePayment.OrderID + validatePayment.StatusCode + validatePayment.GrossAmount + os.Getenv("SERVER_KEY")
	encodeSignatureKey, err := crypto.HashValueSha512(makeSignatureKey)
	if err != nil {
		response.Response(c, http.StatusBadRequest, "Failed to encode signature key", err.Error())
		return
	}

	if err := crypto.CompareHashSha512(encodeSignatureKey, validatePayment.SignatureKey); err != nil {
		response.Response(c, http.StatusBadRequest, "Failed to validate signature key", err.Error())
		return
	}

	if validatePayment.TransactionStatus == "settlement" {
		result.Status = "SUCCESS"

		if err := p.paymentService.UpdatePayment(result, result.ID); err != nil {
			response.Response(c, http.StatusBadRequest, "Failed to update payment", err.Error())
			return
		}

		result, err := p.paymentService.FindById(validatePayment.OrderID)
		if err != nil {
			response.Response(c, http.StatusBadRequest, "Failed to find order id", err.Error())
			return
		}

		if result.Status == "SUCCESS" {
			user, err := p.userService.GetUserByID(result.UserID)
			if err != nil {
				response.Response(c, http.StatusBadRequest, "Failed to get user", err.Error())
				return
			}

			user.IsPremium = true

			if err := p.userService.UpdateUser(user); err != nil {
				response.Response(c, http.StatusBadRequest, "Failed to update user", err.Error())
				return
			}

			userPremium := entity.UserPremium{
				StartDate: time.Now(),
				EndDate:   time.Now().AddDate(0, result.Lifetime, 0),
				UserID:    result.UserID,
			}

			if err := p.userService.CreatePremium(userPremium); err != nil {
				response.Response(c, http.StatusBadRequest, "Failed to create premium", err.Error())
				return
			}
		}
	}
}

func (p *PaymentHandler) GetTransaction(c *gin.Context) {
	id := c.MustGet("user").(entity.User).ID

	result, err := p.paymentService.FindAll(id)
	if err != nil {
		response.Response(c, http.StatusBadRequest, "Failed to get transaction", err.Error())
		return
	}

	response.Response(c, http.StatusOK, "success", result)
}
