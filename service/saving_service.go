package service

import (
	"quell-api/entity"
	"quell-api/models"
	"quell-api/repository"
)

type SavingService interface {
	FindAll() ([]models.Saving, error)
	FindById(id uint) (models.Saving, error)
	CreateSaving(saving entity.Saving) error
	UpdateSaving(saving entity.Saving, id uint) error
	DeleteSaving(id uint) error
}

type savingService struct {
	repository repository.SavingRepository
}

func NewSavingService(repository repository.SavingRepository) SavingService {
	return &savingService{repository}
}

func (s *savingService) FindAll() ([]models.Saving, error) {
	return s.repository.FindAll()
}

func (s *savingService) FindById(id uint) (models.Saving, error) {
	return s.repository.FindById(id)
}

func (s *savingService) CreateSaving(saving entity.Saving) error {
	return s.repository.CreateSaving(saving)
}

func (s *savingService) UpdateSaving(saving entity.Saving, id uint) error {
	return s.repository.UpdateSaving(saving, id)
}

func (s *savingService) DeleteSaving(id uint) error {
	return s.repository.DeleteSaving(id)
}
