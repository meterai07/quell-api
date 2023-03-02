package service

import (
	"quell-api/entity"
	"quell-api/models"
	"quell-api/repository"
)

type CategoryService interface {
	FindAll() ([]models.Category, error)
	FindById(id uint) (models.Category, error)
	CreateCategory(category entity.Category) error
	UpdateCategory(category entity.Category, id uint) error
	DeleteCategory(id uint) error
}

type categoryService struct {
	repository repository.CategoryRepository
}

func NewCategoryService(repository repository.CategoryRepository) CategoryService {
	return &categoryService{repository}
}

func (s *categoryService) FindAll() ([]models.Category, error) {
	return s.repository.FindAll()
}

func (s *categoryService) FindById(id uint) (models.Category, error) {
	return s.repository.FindById(id)
}

func (s *categoryService) CreateCategory(category entity.Category) error {
	return s.repository.CreateCategory(category)
}

func (s *categoryService) UpdateCategory(category entity.Category, id uint) error {
	return s.repository.UpdateCategory(category, id)
}

func (s *categoryService) DeleteCategory(id uint) error {
	return s.repository.DeleteCategory(id)
}
