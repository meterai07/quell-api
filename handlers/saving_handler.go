package handlers

import (
	"net/http"
	"quell-api/entity"
	"quell-api/models"
	"quell-api/sdk/response"
	"quell-api/service"
	"strconv"

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

func (s *saving_handler) UpdateSavingHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Response(c, http.StatusBadRequest, "failed when parsing id", nil)
		return
	}

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

	err = s.savingService.UpdateSaving(newBody, uint(id))
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "failed when update saving", nil)
		return
	}
	response.Response(c, http.StatusOK, "success", nil)
}

func (s *saving_handler) GetTotalAmountHandler(c *gin.Context) {
	totalAmount, err := s.savingService.GetTotalAmount(c.MustGet("user").(entity.User).ID)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "failed when get total amount", nil)
		return
	}
	response.Response(c, http.StatusOK, "success", totalAmount)
}
