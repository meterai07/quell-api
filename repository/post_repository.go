package repository

import (
	"quell-api/entity"
	"quell-api/models"

	"gorm.io/gorm"
)

type PostRepository interface {
	FindAll() ([]entity.Post, error)
	FindAllPostsByUserID(id uint) ([]entity.Post, error)
	FindById(id uint) (entity.Post, error)
	CreatePost(post entity.Post) error
	UpdatePost(post models.Post_Upload, id uint) error
	DeletePost(id uint) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db}
}

func (r *postRepository) FindAll() ([]entity.Post, error) {
	var posts []entity.Post
	result := r.db.Find(&posts).Error
	if result != nil {
		return posts, result
	}
	return posts, nil
}

func (r *postRepository) FindAllPostsByUserID(id uint) ([]entity.Post, error) {
	var posts []entity.Post
	result := r.db.Where("user_id = ?", id).Find(&posts).Error
	if result != nil {
		return posts, result
	}
	return posts, nil
}

func (r *postRepository) FindById(id uint) (entity.Post, error) {
	var post entity.Post
	result := r.db.Where("id = ?", id).First(&post).Error
	if result != nil {
		return post, result
	}
	return post, nil
}

func (r *postRepository) CreatePost(post entity.Post) error {
	result := r.db.Create(&post).Error
	if result != nil {
		return result
	}
	return nil
}

func (r *postRepository) UpdatePost(post models.Post_Upload, id uint) error {
	var postUpdate entity.Post
	result := r.db.Where("id = ?", id).First(&postUpdate).Error
	if result != nil {
		return result
	}

	if post.Title != "" {
		postUpdate.Title = post.Title
	}

	if post.Content != "" {
		postUpdate.Content = post.Content
	}

	if post.CategoryID != (nil) {
		postUpdate.CategoryID = post.CategoryID
	}

	result = r.db.Save(&postUpdate).Error
	if result != nil {
		return result
	}
	return nil
}

func (r *postRepository) DeletePost(id uint) error {
	var post entity.Post
	result := r.db.Where("id = ?", id).Delete(&post).Error
	if result != nil {
		return result
	}
	return nil
}
