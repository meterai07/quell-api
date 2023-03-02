package service

import (
	"quell-api/entity"
	"quell-api/models"
	"quell-api/repository"
)

type PostService interface {
	FindAll() ([]entity.Post, error)
	FindById(id uint) (entity.Post, error)
	CreatePost(post entity.Post) error
	UpdatePost(post models.Post_Upload, id uint) error
	DeletePost(id uint) error
}

type postService struct {
	repository repository.PostRepository
}

func NewPostService(repository repository.PostRepository) PostService {
	return &postService{repository}
}

func (s *postService) FindAll() ([]entity.Post, error) {
	return s.repository.FindAll()
}

func (s *postService) FindById(id uint) (entity.Post, error) {
	return s.repository.FindById(id)
}

func (s *postService) CreatePost(post entity.Post) error {
	return s.repository.CreatePost(post)
}

func (s *postService) UpdatePost(post models.Post_Upload, id uint) error {
	return s.repository.UpdatePost(post, id)
}

func (s *postService) DeletePost(id uint) error {
	return s.repository.DeletePost(id)
}
