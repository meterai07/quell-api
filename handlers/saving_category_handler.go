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

type SavingCategoryHandler struct {
	savingCategoryService service.SavingCategoryService
}

func NewSavingCategoryHandler(savingCategoryService service.SavingCategoryService) *SavingCategoryHandler {
	return &SavingCategoryHandler{savingCategoryService}
}

func (h *SavingCategoryHandler) GetSavingCategoryHandler(c *gin.Context) {
	result, err := h.savingCategoryService.FindAll()
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "failed when get all saving category", nil)
		return
	}
	response.Response(c, http.StatusOK, "success", result)
}

func (h *SavingCategoryHandler) GetSavingCategoryByIdHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Response(c, http.StatusBadRequest, "failed when parsing id", nil)
		return
	}

	result, err := h.savingCategoryService.FindById(uint(id))
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "failed when get saving category by id", nil)
		return
	}
	response.Response(c, http.StatusOK, "success", result)
}

func (h *SavingCategoryHandler) CreateSavingCategoryHandler(c *gin.Context) {
	var body models.SavingCategory
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Response(c, http.StatusBadRequest, "failed when binding", nil)
		return
	}

	newBody := entity.SavingCategory{
		Name: body.Name,
	}

	err := h.savingCategoryService.CreateSavingCategory(newBody)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "failed when create saving category", nil)
		return
	}
	response.Response(c, http.StatusOK, "success", nil)
}

func (h *SavingCategoryHandler) UpdateSavingCategoryHandler(c *gin.Context) {
	var body models.SavingCategory
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Response(c, http.StatusBadRequest, "failed when binding", nil)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Response(c, http.StatusBadRequest, "failed when parsing id", nil)
		return
	}

	newBody := entity.SavingCategory{
		Name: body.Name,
	}

	err = h.savingCategoryService.UpdateSavingCategory(newBody, uint(id))
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "failed when update saving category", nil)
		return
	}
	response.Response(c, http.StatusOK, "success", nil)
}

func (h *SavingCategoryHandler) DeleteSavingCategoryHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Response(c, http.StatusBadRequest, "failed when parsing id", nil)
		return
	}

	err = h.savingCategoryService.DeleteSavingCategory(uint(id))
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "failed when delete saving category", nil)
		return
	}
	response.Response(c, http.StatusOK, "success", nil)
}
