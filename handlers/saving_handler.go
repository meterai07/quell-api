package handlers

import (
	"net/http"
	"quell-api/entity"
	"quell-api/models"
	"quell-api/sdk/response"
	"quell-api/service"

	"github.com/gin-gonic/gin"
)

type saving_handler struct {
	savingService service.SavingService
}

func NewSavingHandler(savingService service.SavingService) *saving_handler {
	return &saving_handler{savingService}
}

func (s *saving_handler) CreateSavingHandler(c *gin.Context) {
	var body models.Saving
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Response(c, http.StatusBadRequest, "failed when binding", nil)
		return
	}

	checkType := body.Type
	if checkType != "income" && checkType != "expense" {
		response.Response(c, http.StatusBadRequest, "failed when checking type", nil)
		return
	}

	newBody := entity.Saving{
		Name:             body.Name,
		Description:      body.Description,
		Amount:           body.Amount,
		SavingCategoryID: body.SavingCategoryID,
		UserID:           c.MustGet("user").(entity.User).ID,
		Type:             body.Type,
	}
	err := s.savingService.CreateSaving(newBody)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "failed when create saving", nil)
		return
	}
	response.Response(c, http.StatusOK, "success", nil)
}
