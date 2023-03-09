package handlers

import (
	"net/http"
	"os"
	"quell-api/entity"
	"quell-api/sdk/crypto"
	"quell-api/sdk/response"

	"github.com/gin-gonic/gin"
)

func ValidatePayment(c *gin.Context) {
	var validatePayment entity.ValidatePayment

	if err := c.ShouldBindJSON(&validatePayment); err != nil {
		response.Response(c, http.StatusBadRequest, "Failed to bind json", nil)
		return
	}

	makeSignatureKey := validatePayment.OrderID + validatePayment.StatusCode + validatePayment.GrossAmount + os.Getenv("SERVER_KEY")
	encodeSignatureKey, err := crypto.HashValueSha512(makeSignatureKey)
	if err != nil {
		response.Response(c, http.StatusBadRequest, "Failed to encode signature key", err.Error())
		return
	}

	if err := crypto.CompareHash(encodeSignatureKey, validatePayment.SignatureKey); err != nil {
		response.Response(c, http.StatusBadRequest, "Failed to validate signature key", err.Error())
		return
	}
}
