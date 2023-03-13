package handlers

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"quell-api/entity"
	"quell-api/models"
	"quell-api/sdk/response"
	"quell-api/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

type post_Handler struct {
	post_Service service.PostService
}

func NewPostHandler(post_Service service.PostService) *post_Handler {
	return &post_Handler{
		post_Service,
	}
}

func (h *post_Handler) GetPostHandler(c *gin.Context) {
	posts, err := h.post_Service.FindAll()
	if err != nil {
		response.Response(c, http.StatusNotFound, "failed when find all data", nil)
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

var (
	s *gocron.Scheduler
)

func (h *user_Handler) Reminder(c *gin.Context) {
	typeReminder := c.Query("type")

	user := c.MustGet("user").(entity.User)

	if !user.IsPremium {
		response.Response(c, http.StatusUnauthorized, "You must Premium to Access this content", nil)
		return
	}

	if typeReminder == "start" {
		s = gocron.NewScheduler(time.UTC)

		posts, err := h.postService.FindAllPostsByUserID(user.ID)
		if err != nil {
			response.Response(c, http.StatusInternalServerError, "failed when find all data", nil)
			return
		}

		s.Every(2).Hour().Do(CheckPost, c, posts, user)

		s.StartAsync()
		response.Response(c, http.StatusOK, "success", nil)
		return
	}

	if typeReminder == "stop" {
		if s != nil {
			s.Remove(CheckPost)
			response.Response(c, http.StatusOK, "the scheduler stopped", nil)
			return
		}
	}

	response.Response(c, http.StatusBadRequest, "failed when stopping because the feature not enabled", nil)
}

func CheckPost(c *gin.Context, posts []entity.Post, user entity.User) {
	for _, post := range posts {
		if post.Type == "jadwal" || post.Type == "tugas" {
			deadline := post.Date
			now := time.Now()

			if deadline.After(now) {
				duration := deadline.Sub(now)
				if time.Duration(duration.Hours()) < 7*time.Hour && deadline.Year() == now.Year() && deadline.Month() == now.Month() && deadline.Day() == now.Day() {
					fmt.Println("email sent")
					if err := SendRemainderEmail(user.Email, post.Title, post.Type, post.Date); err != nil {
						response.Response(c, http.StatusInternalServerError, "failed when sending email", nil)
					}
				}
			}
		}
	}
}

func SendRemainderEmail(email string, title string, typeTask string, date time.Time) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", from, password, host)
	subject := fmt.Sprintf("%s Reminder: %s on %s", typeTask, title, date.Format("January 2, 2006"))
	body := fmt.Sprintf("Hello,\n\nThis is a friendly reminder that you have a %s scheduled for %s. Don't forget!\n\nBest regards,\nYour Reminder Service, Quell", typeTask, date.Format("Monday, January 2, 2006"))

	msg := []byte("To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(host+":"+port, auth, from, []string{email}, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}
