package repository

import (
	"quell-api/entity"
	"quell-api/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll() ([]models.Category, error)
	FindById(id uint) (models.Category, error)
	CreateCategory(category entity.Category) error
	UpdateCategory(category entity.Category, id uint) error
	DeleteCategory(id uint) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) FindAll() ([]models.Category, error) {
	var categories []models.Category
	result := r.db.Find(&categories).Error
	if result != nil {
		return categories, result
	}
	return categories, nil
}

func (r *categoryRepository) FindById(id uint) (models.Category, error) {
	var category models.Category
	result := r.db.Where("id = ?", id).First(&category).Error
	if result != nil {
		return category, result
	}
	return category, nil
}

func (r *categoryRepository) CreateCategory(category entity.Category) error {

	result := r.db.Create(&category).Error
	if result != nil {
		return result
	}
	return nil
}

func (r *categoryRepository) UpdateCategory(category entity.Category, id uint) error {
	var categoryUpdate entity.Category
	result := r.db.Where("id = ?", id).First(&categoryUpdate).Error
	if result != nil {
		return result
	}
	categoryUpdate.Name = category.Name
	result = r.db.Save(&categoryUpdate).Error
	if result != nil {
		return result
	}
	return nil
}

func (r *categoryRepository) DeleteCategory(id uint) error {
	var category models.Category
	result := r.db.Where("id = ?", id).Delete(&category).Error
	if result != nil {
		return result
	}
	return nil
}
