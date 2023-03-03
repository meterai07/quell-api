package handlers

import (
	"path/filepath"
	"quell-api/entity"
	"quell-api/sdk/response"
	"quell-api/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type attachmentHandler struct {
	service service.AttachmentService
}

func NewAttachmentHandler(service service.AttachmentService) *attachmentHandler {
	return &attachmentHandler{service}
}

func (h *attachmentHandler) UploadFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Response(c, 400, "failed when parsing id", nil)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.Response(c, 400, "Invalid File", err.Error())
		return
	}

	result := c.PostForm("name")
	if result == "" {
		response.Response(c, 400, "Invalid File", err.Error())
		return
	}

	// newName := generateNewName(file.Filename)

	link, err := h.service.UploadFile(file)
	if err != nil {
		response.Response(c, 500, "Internal Server Error", err.Error())
		return
	}

	response.Response(c, 200, "Success", link)

	newBody := entity.Attachment{
		Name:   result,
		UserID: c.MustGet("user").(entity.User).ID,
		PostID: uint(id),
		Url:    link,
	}

	err = h.service.CreateAttachment(newBody)
	if err != nil {
		response.Response(c, 500, "Internal Server Error", err.Error())
		return
	}
}

func generateNewName(filename string) string {
	uuid := uuid.New()
	ext := filepath.Ext(filename)
	return uuid.String() + ext
}

func (h *attachmentHandler) DeleteFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Response(c, 400, "failed when parsing id", nil)
		return
	}

	attid, err := strconv.ParseUint(c.Param("attid"), 10, 64)
	if err != nil {
		response.Response(c, 400, "failed when parsing id", nil)
		return
	}

	result, err := h.service.FindById(uint(attid))
	if err != nil {
		response.Response(c, 500, "Internal Server Error", err.Error())
		return
	}

	if result.PostID != uint(id) {
		response.Response(c, 400, "Invalid File", err.Error())
		return
	}

	data, err := h.service.DeleteFile(result.Url)

	if err != nil {
		response.Response(c, 500, "Internal Server Error", err.Error())
		return
	}

	if err := h.service.DeleteAttachment(uint(attid)); err != nil {
		response.Response(c, 500, "Internal Server Error", err.Error())
		return
	}

	response.Response(c, 200, "Success", data)
}
