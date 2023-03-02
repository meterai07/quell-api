package repository

import (
	"mime/multipart"
	"quell-api/entity"
	"quell-api/initializers"

	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
)

type AttachmentRepository interface {
	FindAll() ([]entity.Attachment, error)
	FindById(id uint) (entity.Attachment, error)
	FindAttachmentById(id uint) (entity.Attachment, error)
	CreateAttachment(attachment entity.Attachment) error
	UpdateAttachment(attachment entity.Attachment, id uint) error
	DeleteAttachment(id uint) error

	UploadFile(file *multipart.FileHeader) (string, error)
	DeleteFile(linkFile string) (interface{}, error)
}

type attachmentRepository struct {
	supClient supabasestorageuploader.SupabaseClientService
}

func NewAttachmentRepository(supClient supabasestorageuploader.SupabaseClientService) AttachmentRepository {
	return &attachmentRepository{supClient}
}

func (r *attachmentRepository) FindAll() ([]entity.Attachment, error) {
	var attachments []entity.Attachment
	err := initializers.DB.Find(&attachments).Error
	if err != nil {
		return attachments, err
	}
	return attachments, nil
}

func (r *attachmentRepository) FindById(id uint) (entity.Attachment, error) {
	var attachment entity.Attachment
	err := initializers.DB.Where("id = ?", id).First(&attachment).Error
	if err != nil {
		return attachment, err
	}
	return attachment, nil
}

func (r *attachmentRepository) FindAttachmentById(id uint) (entity.Attachment, error) {
	var attachment entity.Attachment
	err := initializers.DB.Where("id = ?", id).First(&attachment).Error
	if err != nil {
		return attachment, err
	}
	return attachment, nil
}

func (r *attachmentRepository) CreateAttachment(attachment entity.Attachment) error {
	err := initializers.DB.Create(&attachment).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *attachmentRepository) UpdateAttachment(attachment entity.Attachment, id uint) error {
	var attachmentUpdate entity.Attachment
	err := initializers.DB.Where("id = ?", id).First(&attachmentUpdate).Error
	if err != nil {
		return err
	}
	attachmentUpdate.Name = attachment.Name
	attachmentUpdate.Url = attachment.Url
	err = initializers.DB.Save(&attachmentUpdate).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *attachmentRepository) DeleteAttachment(id uint) error {
	var attachment entity.Attachment
	err := initializers.DB.Where("id = ?", id).First(&attachment).Error
	if err != nil {
		return err
	}
	err = initializers.DB.Delete(&attachment).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *attachmentRepository) UploadFile(file *multipart.FileHeader) (string, error) {
	link, err := initializers.SupabaseClient.Upload(file)
	if err != nil {
		return "", err
	}
	return link, nil
}

func (r *attachmentRepository) DeleteFile(linkFile string) (interface{}, error) {
	data, err := initializers.SupabaseClient.DeleteFile(linkFile)
	if err != nil {
		return "", err
	}
	return data, nil
}
