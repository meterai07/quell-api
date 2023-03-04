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

type categoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(service service.CategoryService) *categoryHandler {
	return &categoryHandler{service}
}

func (h *categoryHandler) GetCategoryHandler(c *gin.Context) {
	posts, err := h.service.FindAll()
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "failed when find all data", nil)
		return
	}
	response.Response(c, http.StatusOK, "success", posts)
}

func (h *categoryHandler) GetCategoryByIdHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Response(c, http.StatusBadRequest, "failed when parsing id", nil)
		return
	}

	post, err := h.service.FindById(uint(id))
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "failed when get posts by id", nil)
		return
	}
	response.Response(c, http.StatusOK, "success", post)
}

func (h *categoryHandler) CreateCategoryHandler(c *gin.Context) {
	var body models.Category
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Response(c, http.StatusBadRequest, "failed when binding", nil)
		return
	}

	newBody := entity.Category{
		Name: body.Name,
	}

	err := h.service.CreateCategory(newBody)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "failed when create category", nil)
		return
	}
	response.Response(c, 201, "success", nil)
}

func (h *categoryHandler) UpdateCategoryHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Response(c, http.StatusBadRequest, "failed when parsing id", nil)
		return
	}

	var body models.Category
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Response(c, http.StatusBadRequest, "failed when binding", nil)
		return
	}

	newBody := entity.Category{
		Name: body.Name,
	}

	err = h.service.UpdateCategory(newBody, uint(id))
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "failed when update category", nil)
		return
	}
	response.Response(c, http.StatusOK, "success", nil)
}

func (h *categoryHandler) DeleteCategoryHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Response(c, http.StatusBadRequest, "failed when parsing id", nil)
		return
	}

	err = h.service.DeleteCategory(uint(id))
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "failed when delete category", nil)
		return
	}
	response.Response(c, http.StatusOK, "success", nil)
}
