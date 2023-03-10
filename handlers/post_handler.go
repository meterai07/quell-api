package handlers

import (
	"fmt"
	"net/http"
	"quell-api/entity"
	"quell-api/models"
	"quell-api/sdk/response"
	"quell-api/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type post_Handler struct {
	post_Service service.PostService
	user_Service service.Service
}

func NewPostHandler(post_Service service.PostService, user_Service service.Service) *post_Handler {
	return &post_Handler{
		post_Service,
		user_Service,
	}
}

func (h *post_Handler) GetPostHandler(c *gin.Context) {
	posts, err := h.post_Service.FindAll()
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "failed when find all data", nil)
		return
	}
	response.Response(c, http.StatusOK, "success", posts)
}

func (h *post_Handler) GetPostByIdHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Response(c, http.StatusBadRequest, "failed when parsing id", nil)
		return
	}

	post, err := h.post_Service.FindById(uint(id))
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "failed when find by id", nil)
		return
	}
	response.Response(c, http.StatusOK, "success", post)
}

func (h *post_Handler) CreatePostHandler(c *gin.Context) {
	var body models.Post_Upload

	if err := c.ShouldBindJSON(&body); err != nil {
		response.Response(c, http.StatusBadRequest, "failed when binding", nil)
		return
	}

	// parse string to time
	parseDate, err := time.Parse("2006-01-02 15:04:05", body.Date)
	if err != nil {
		response.Response(c, http.StatusBadRequest, "failed when parsing date", nil)
		return
	}

	deadline := parseDate.Add(-7 * time.Hour)

	checkType := body.Type
	if checkType != "jadwal" && checkType != "tugas" {
		response.Response(c, http.StatusBadRequest, "failed when checking type", nil)
		return
	}

	newBody := entity.Post{
		Title:      body.Title,
		Content:    body.Content,
		Date:       deadline,
		Type:       body.Type,
		UserID:     c.MustGet("user").(entity.User).ID,
		CategoryID: body.CategoryID,
	}

	if err := h.post_Service.CreatePost(newBody); err != nil {
		response.Response(c, http.StatusInternalServerError, "failed when creating post", nil)
		return
	}

	response.Response(c, 200, "success", nil)
}

func (h *post_Handler) UpdatePostHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Response(c, http.StatusBadRequest, "failed when parsing id", nil)
		return
	}

	var body models.Post_Upload
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Response(c, http.StatusBadRequest, "failed when binding", nil)
		return
	}

	if err := h.post_Service.UpdatePost(body, uint(id)); err != nil {
		response.Response(c, http.StatusInternalServerError, "failed", nil)
		return
	}
	response.Response(c, http.StatusOK, "success", nil)
}

func (h *post_Handler) DeletePostHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Response(c, http.StatusBadRequest, "failed when parsing id", nil)
		return
	}

	err = h.post_Service.DeletePost(uint(id))
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "failed", nil)
		return
	}
	response.Response(c, http.StatusOK, "success", nil)
}

func (h *post_Handler) ActivateReminder(c *gin.Context) {
	user, err := h.user_Service.GetUserByID(c.MustGet("user").(entity.User).ID)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "failed when find user", nil)
		return
	}

	if !user.IsPremium {
		response.Response(c, http.StatusBadRequest, "user is not premium to activate this feature", nil)
		return
	}

	posts, err := h.post_Service.FindAllPostsByUserID(c.MustGet("user").(entity.User).ID)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "failed when find all data", nil)
		return
	}

	for _, post := range posts {
		if post.Type == "jadwal" || post.Type == "tugas" {
			deadline := post.Date
			now := time.Now()

			fmt.Println("Deadline:", deadline)
			fmt.Println("Now:", now)

			if deadline.After(now) {
				fmt.Println("Deadline is later than now")
			} else {
				fmt.Println("Deadline has passed")
			}
		}
	}

}
