package repository

import (
	"quell-api/entity"
	"quell-api/models"

	"gorm.io/gorm"
)

type SavingRepository interface {
	FindAll() ([]models.Saving, error)
	FindById(id uint) (models.Saving, error)
	GetTotalAmount(id uint) (int, error)
	CreateSaving(saving entity.Saving) error
	UpdateSaving(saving entity.Saving, id uint) error
	DeleteSaving(id uint) error
}

type savingRepository struct {
	db *gorm.DB
}

func NewSavingRepository(db *gorm.DB) SavingRepository {
	return &savingRepository{db}
}

func (s *savingRepository) FindAll() ([]models.Saving, error) {
	var savings []models.Saving
	result := s.db.Find(&savings).Error
	if result != nil {
		return savings, result
	}
	return savings, nil
}

func (s *savingRepository) FindById(id uint) (models.Saving, error) {
	var saving models.Saving
	result := s.db.Where("id = ?", id).First(&saving).Error
	if result != nil {
		return saving, result
	}
	return saving, nil
}

func (s *savingRepository) GetTotalAmount(id uint) (int, error) {
	var totalAmount int
	result := s.db.Table("savings").Where("user_id = ?", id).Select("SUM(amount)").Scan(&totalAmount).Error
	if result != nil {
		return totalAmount, result
	}
	return totalAmount, nil
}

func (s *savingRepository) CreateSaving(saving entity.Saving) error {
	result := s.db.Create(&saving).Error
	if result != nil {
		return result
	}
	return nil
}

func (s *savingRepository) UpdateSaving(saving entity.Saving, id uint) error {
	var savingUpdate entity.Saving
	result := s.db.Where("id = ?", id).First(&savingUpdate).Error
	if result != nil {
		return result
	}

	if saving.Name != "" {
		savingUpdate.Name = saving.Name
	}

	if saving.Amount != 0 {
		savingUpdate.Amount = saving.Amount
	}

	if saving.Description != "" {
		savingUpdate.Description = saving.Description
	}

	if saving.SavingCategoryID != (nil) {
		savingUpdate.SavingCategoryID = saving.SavingCategoryID
	}

	result = s.db.Save(&savingUpdate).Error
	if result != nil {
		return result
	}
	return nil
}

func (s *savingRepository) DeleteSaving(id uint) error {
	var saving entity.Saving
	result := s.db.Where("id = ?", id).First(&saving).Error
	if result != nil {
		return result
	}
	result = s.db.Delete(&saving).Error
	if result != nil {
		return result
	}
	return nil
}
