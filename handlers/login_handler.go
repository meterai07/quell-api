package handlers

import (
	"fmt"
	"net/http"
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

		s.Every(5).Second().Do(func() {
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
							fmt.Println("Reminder activated for ", post.ID)
						}
					}
				}
			}
		})

		// s.Every(5).Second().Do(func(p *post_Handler) {
		// 	fmt.Println("masuk")
		// 	p.ActivateReminder(c)
		// })

		s.StartAsync()
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "/", "", false, true)

	response.Response(c, http.StatusOK, "User signed in successfully", nil)
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
			response.Response(c, http.StatusInternalServerError, "Internal Server Error When Getting User", nil)
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
