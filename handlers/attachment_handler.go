package handlers

import (
	"net/http"
	"quell-api/entity"
	"quell-api/sdk/response"
	"quell-api/service"
	"strconv"

	"github.com/gin-gonic/gin"
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
		response.Response(c, http.StatusBadRequest, "failed when parsing id", nil)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.Response(c, http.StatusBadRequest, "Invalid File", err.Error())
		return
	}

	result := c.PostForm("name")
	if result == "" {
		response.Response(c, http.StatusBadRequest, "Invalid File", err.Error())
		return
	}

	// newName := generateNewName(file.Filename)

	link, err := h.service.UploadFile(file)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	response.Response(c, http.StatusOK, "Success", link)

	newBody := entity.Attachment{
		Name:   result,
		UserID: c.MustGet("user").(entity.User).ID,
		PostID: uint(id),
		Url:    link,
	}
	// test
	err = h.service.CreateAttachment(newBody)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}
}

func (h *attachmentHandler) DeleteFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Response(c, http.StatusBadRequest, "failed when parsing id", nil)
		return
	}

	attid, err := strconv.ParseUint(c.Param("attid"), 10, 64)
	if err != nil {
		response.Response(c, http.StatusBadRequest, "failed when parsing id", nil)
		return
	}

	result, err := h.service.FindById(uint(attid))
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	if result.PostID != uint(id) {
		response.Response(c, http.StatusBadRequest, "Invalid File", err.Error())
		return
	}

	data, err := h.service.DeleteFile(result.Url)

	if err != nil {
		response.Response(c, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	if err := h.service.DeleteAttachment(uint(attid)); err != nil {
		response.Response(c, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	response.Response(c, http.StatusOK, "Success", data)
}
