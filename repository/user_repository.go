package repository

import (
	"quell-api/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user entity.User) error
	GetUserByEmail(email string) (entity.User, error)
	FindUserByEmail(email string) bool
	GetUserByID(id uint) (entity.User, error)
	UpdateUser(user entity.User) error
	DeleteUser(user entity.User) error
	CreatePremium(premium entity.UserPremium) error
	GetPremiumByID(id uint) (entity.UserPremium, error)
	FindPremiumByID(id uint) bool
	UpdatePremium(premium entity.UserPremium) error
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) CreateUser(user entity.User) error {
	result := r.DB.Create(&user).Error
	if result != nil {
		return result
	}
	return nil
}

func (r *userRepository) GetUserByEmail(email string) (entity.User, error) {
	var user entity.User
	result := r.DB.Where("email = ?", email).First(&user).Error
	if result != nil {
		return user, result
	}
	return user, nil
}

func (r *userRepository) FindUserByEmail(email string) bool {
	var user entity.User
	result := r.DB.Where("email = ?", email).First(&user).Error
	return result == nil
}

func (r *userRepository) GetUserByID(id uint) (entity.User, error) {
	var user entity.User
	result := r.DB.Where("id = ?", id).First(&user).Error
	if result != nil {
		return user, result
	}
	return user, nil
}

func (r *userRepository) UpdateUser(user entity.User) error {
	result := r.DB.Save(&user).Error
	if result != nil {
		return result
	}
	return nil
}

func (r *userRepository) DeleteUser(user entity.User) error {
	result := r.DB.Delete(&user).Error
	if result != nil {
		return result
	}
	return nil
}

func (r *userRepository) CreatePremium(premium entity.UserPremium) error {
	result := r.DB.Create(&premium).Error
	if result != nil {
		return result
	}
	return nil
}

func (r *userRepository) GetPremiumByID(id uint) (entity.UserPremium, error) {
	var premium entity.UserPremium
	result := r.DB.Where("user_id = ?", id).First(&premium).Error
	if result != nil {
		return premium, result
	}
	return premium, nil
}

func (r *userRepository) FindPremiumByID(id uint) bool {
	var premium entity.UserPremium
	result := r.DB.Where("user_id = ?", id).First(&premium).Error
	return result == nil
}

func (r *userRepository) UpdatePremium(premium entity.UserPremium) error {
	result := r.DB.Save(&premium).Error
	if result != nil {
		return result
	}
	return nil
}
