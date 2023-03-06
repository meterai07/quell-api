package repository

import (
	"quell-api/entity"
	"quell-api/models"

	"gorm.io/gorm"
)

type SavingCategoryRepository interface {
	FindAll() ([]models.SavingCategory, error)
	FindById(id uint) (models.SavingCategory, error)
	CreateSavingCategory(savingCategory entity.SavingCategory) error
	UpdateSavingCategory(savingCategory entity.SavingCategory, id uint) error
	DeleteSavingCategory(id uint) error
}

type savingCategoryRepository struct {
	db *gorm.DB
}

func NewSavingCategoryRepository(db *gorm.DB) SavingCategoryRepository {
	return &savingCategoryRepository{db}
}

func (r *savingCategoryRepository) FindAll() ([]models.SavingCategory, error) {
	var savingCategories []models.SavingCategory
	err := r.db.Find(&savingCategories).Error
	if err != nil {
		return nil, err
	}
	return savingCategories, nil
}

func (r *savingCategoryRepository) FindById(id uint) (models.SavingCategory, error) {
	var savingCategory models.SavingCategory
	err := r.db.Where("id = ?", id).First(&savingCategory).Error
	if err != nil {
		return savingCategory, err
	}
	return savingCategory, nil
}

func (r *savingCategoryRepository) CreateSavingCategory(savingCategory entity.SavingCategory) error {
	err := r.db.Create(&savingCategory).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *savingCategoryRepository) UpdateSavingCategory(savingCategory entity.SavingCategory, id uint) error {
	err := r.db.Where("id = ?", id).Updates(&savingCategory).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *savingCategoryRepository) DeleteSavingCategory(id uint) error {
	err := r.db.Where("id = ?", id).Delete(&models.SavingCategory{}).Error
	if err != nil {
		return err
	}
	return nil
}
