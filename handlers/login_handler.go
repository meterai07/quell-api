package handlers

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"quell-api/entity"
	"quell-api/models"
	"quell-api/sdk/crypto"
	sdk_jwt "quell-api/sdk/jwt"
	"quell-api/sdk/response"
	"quell-api/service"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
)

type user_Handler struct {
	userService service.Service
	postService service.PostService
}

func NewUserHandler(userService service.Service, postService service.PostService) *user_Handler {
	return &user_Handler{
		userService: userService,
		postService: postService,
	}
}

func (h *user_Handler) LoginHandler(c *gin.Context) {
	var body models.User_Login

	if err := c.ShouldBindJSON(&body); err != nil {
		response.Response(c, http.StatusBadRequest, "Invalid Body Email", nil)
		return
	}

	if err := validator.New().Struct(&body); err != nil {
		validationError := err.(validator.ValidationErrors)
		response.Response(c, http.StatusBadRequest, validationError.Error(), nil)
		return
	}

	var user entity.User
	result, err := h.userService.GetUserByEmail(body.Email)
	if err != nil {
		response.Response(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	user = result

	if !user.IsActive {
		response.Response(c, http.StatusUnauthorized, "Email Not Validated", nil)
		return
	}

	if err := crypto.CompareHash(user.Password, body.Password); err != nil {
		response.Response(c, http.StatusUnauthorized, "Password Not Validated", nil)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "Error while signing the token", nil)
		return
	}

	if user.IsPremium {
		s := gocron.NewScheduler(time.UTC)

		s.Every(1).Minute().Do(func() {
			posts, err := h.postService.FindAllPostsByUserID(user.ID)
			if err != nil {
				response.Response(c, http.StatusInternalServerError, "failed when find all data", nil)
				return
			}

			for _, post := range posts {
				if post.Type == "jadwal" || post.Type == "tugas" {
					deadline := post.Date
					now := time.Now()

					if deadline.After(now) {
						duration := deadline.Sub(now)
						if time.Duration(duration.Hours()) < 7*time.Hour && deadline.Year() == now.Year() && deadline.Month() == now.Month() && deadline.Day() == now.Day() {
							fmt.Println("test")
							if err := SendRemainderEmail(user.Email, post.Title, post.Type, post.Date); err != nil {
								response.Response(c, http.StatusInternalServerError, "failed when send email", nil)
								return
							}
						}
					}
				}
			}
		})

		s.StartAsync()
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "/", "", false, true)

	response.Response(c, http.StatusOK, "User signed in successfully", nil)
}

func SendRemainderEmail(email string, title string, typeTask string, date time.Time) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", from, password, host)
	subject := fmt.Sprintf("%s Reminder: %s on %s", typeTask, title, date.Format("January 2, 2006"))
	body := fmt.Sprintf("Hello,\n\nThis is a friendly reminder that you have a %s scheduled for %s. Don't forget!\n\nBest regards,\nYour Reminder Service", typeTask, date.Format("Monday, January 2, 2006"))

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

func (h *user_Handler) RegisterHandler(c *gin.Context) {
	var body models.User_Register

	if err := c.ShouldBindJSON(&body); err != nil {
		response.Response(c, http.StatusBadRequest, "Invalid Body Email", nil)
		c.Abort()
		return
	}

	if err := validator.New().Struct(&body); err != nil {
		validationError := err.(validator.ValidationErrors)
		response.Response(c, http.StatusBadRequest, validationError.Error(), nil)
		return
	}

	if result := h.userService.FindUserByEmail(body.Email); result {
		user, err := h.userService.GetUserByEmail(body.Email)
		if err != nil {
			response.Response(c, http.StatusNotFound, "User Not Found", nil)
			c.Abort()
			return
		}

		if user.IsActive {
			response.Response(c, http.StatusConflict, "Email Already Registered", nil)
			c.Abort()
			return
		}

		signedToken, err := sdk_jwt.GenerateToken(user)
		if err != nil {
			response.Response(c, http.StatusInternalServerError, "Internal Server Error When Generating Token", nil)
			c.Abort()
			return
		}

		user.Token = signedToken

		if err := h.userService.UpdateUser(user); err != nil {
			response.Response(c, http.StatusInternalServerError, "Internal Server Error When Updating User", nil)
			c.Abort()
			return
		}

		if err := SendValidationEmail(body.Email, signedToken); err != nil {
			response.Response(c, http.StatusInternalServerError, "Internal Server Error When Sending Email", nil)
			c.Abort()
			return
		}
		response.Response(c, http.StatusOK, "Email Validation Sent", nil)
		c.Abort()
		return
	}
	hashedPassword, err := crypto.HashValue(body.Password)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "Internal Server Error When Hashing Password", nil)
		c.Abort()
		return
	}

	signedToken, err := sdk_jwt.GenerateToken(entity.User{Email: body.Email})
	if err != nil {
		response.Response(c, http.StatusInternalServerError, "Internal Server Error When Generating Token", nil)
		c.Abort()
		return
	}

	user := entity.User{
		Username: body.Username,
		Email:    body.Email,
		Password: string(hashedPassword),
		Phone:    body.Phone,
		IsActive: false,
		Token:    signedToken,
	}

	if err := h.userService.CreateUser(user); err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			response.Response(c, http.StatusConflict, "Email Already Registered", nil)
			return
		}
		response.Response(c, http.StatusInternalServerError, "Internal Server Error When Creating User", nil)
		return
	}

	if err := SendValidationEmail(body.Email, signedToken); err != nil {
		response.Response(c, http.StatusInternalServerError, "Internal Server Error When Sending Email First Create Email", nil)
		c.Abort()
		return
	}

	response.Response(c, http.StatusCreated, "User Created", nil)
	c.Abort()
}

func LogoutHandler(c *gin.Context) {
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("Authorization", "", -1, "/", "", false, true)

	response.Response(c, http.StatusOK, "User Logged Out", nil)
}
