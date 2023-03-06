package service

import (
	"quell-api/entity"
	"quell-api/models"
	"quell-api/repository"
)

type SavingCategoryService interface {
	FindAll() ([]models.SavingCategory, error)
	FindById(id uint) (models.SavingCategory, error)
	CreateSavingCategory(savingCategory entity.SavingCategory) error
	UpdateSavingCategory(savingCategory entity.SavingCategory, id uint) error
	DeleteSavingCategory(id uint) error
}

type savingCategoryService struct {
	repository repository.SavingCategoryRepository
}

func NewSavingCategoryService(repository repository.SavingCategoryRepository) SavingCategoryService {
	return &savingCategoryService{repository}
}

func (s *savingCategoryService) FindAll() ([]models.SavingCategory, error) {
	return s.repository.FindAll()
}

func (s *savingCategoryService) FindById(id uint) (models.SavingCategory, error) {
	return s.repository.FindById(id)
}

func (s *savingCategoryService) CreateSavingCategory(savingCategory entity.SavingCategory) error {
	return s.repository.CreateSavingCategory(savingCategory)
}

func (s *savingCategoryService) UpdateSavingCategory(savingCategory entity.SavingCategory, id uint) error {
	return s.repository.UpdateSavingCategory(savingCategory, id)
}

func (s *savingCategoryService) DeleteSavingCategory(id uint) error {
	return s.repository.DeleteSavingCategory(id)
}
