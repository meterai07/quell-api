package handlers

import (
	"net/http"
	"quell-api/entity"
	"quell-api/sdk/response"

	"github.com/gin-gonic/gin"
)

func ValidatePayment(c *gin.Context) {
	var validatePayment entity.ValidatePayment

	if err := c.ShouldBindJSON(&validatePayment); err != nil {
		response.Response(c, http.StatusBadRequest, "Failed to bind json", nil)
		return
	}
}
