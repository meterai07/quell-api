package service

import (
	"mime/multipart"
	"quell-api/entity"
	"quell-api/repository"
)

type AttachmentService interface {
	FindAll() ([]entity.Attachment, error)
	FindById(id uint) (entity.Attachment, error)
	CreateAttachment(attachment entity.Attachment) error
	UpdateAttachment(attachment entity.Attachment, id uint) error
	DeleteAttachment(id uint) error

	UploadFile(file *multipart.FileHeader) (string, error)
	DeleteFile(linkFile string) (interface{}, error)
}

type attachmentService struct {
	attachmentRepository repository.AttachmentRepository
}

func NewAttachmentService(attachmentRepository repository.AttachmentRepository) AttachmentService {
	return &attachmentService{attachmentRepository}
}

func (s *attachmentService) FindAll() ([]entity.Attachment, error) {
	return s.attachmentRepository.FindAll()
}

func (s *attachmentService) FindById(id uint) (entity.Attachment, error) {
	return s.attachmentRepository.FindById(id)
}

func (s *attachmentService) CreateAttachment(attachment entity.Attachment) error {
	return s.attachmentRepository.CreateAttachment(attachment)
}

func (s *attachmentService) UpdateAttachment(attachment entity.Attachment, id uint) error {
	return s.attachmentRepository.UpdateAttachment(attachment, id)
}

func (s *attachmentService) DeleteAttachment(id uint) error {
	return s.attachmentRepository.DeleteAttachment(id)
}

func (s *attachmentService) UploadFile(file *multipart.FileHeader) (string, error) {
	return s.attachmentRepository.UploadFile(file)
}

func (s *attachmentService) DeleteFile(linkFile string) (interface{}, error) {
	return s.attachmentRepository.DeleteFile(linkFile)
}
