package handlers

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"quell-api/entity"
	"quell-api/sdk/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func (h *user_Handler) ValidateHandler(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		response.Response(c, http.StatusBadRequest, "Invalid Email", nil)
		c.Abort()
		return
	}
	userToken := c.Query("token")
	if userToken == "" {
		response.Response(c, http.StatusBadRequest, "Invalid Token", nil)
		c.Abort()
		return
	}

	token, err := jwt.ParseWithClaims(userToken, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		response.Response(c, http.StatusBadRequest, err.Error(), nil)
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		response.Response(c, http.StatusBadRequest, "Invalid Token", nil)
		c.Abort()
		return
	}
	expTime := time.Unix(int64(claims["exp"].(float64)), 0)
	if time.Until(expTime) < 0 {
		var deleteUser entity.User
		result, err := h.userService.GetUserByEmail(email)
		if err != nil {
			response.Response(c, http.StatusBadRequest, "Email not registered", nil)
			c.Abort()
			return
		}
		deleteUser = result

		if err := h.userService.DeleteUser(deleteUser); err != nil {
			response.Response(c, http.StatusBadRequest, "Error deleting user", nil)
			c.Abort()
			return
		}

		response.Response(c, http.StatusBadRequest, "Token expired", nil)
		c.Abort()
		return
	}

	var user entity.User

	result, err := h.userService.GetUserByEmail(email)
	if err != nil {
		response.Response(c, http.StatusBadRequest, "Email not registered", nil)
		c.Abort()
		return
	}
	user = result

	user.IsActive = true

	if err := h.userService.UpdateUser(user); err != nil {
		response.Response(c, http.StatusBadRequest, "Error updating user", nil)
		c.Abort()
		return
	}

	response.Response(c, http.StatusOK, "Email validated", nil)
	c.Abort()
}

func SendValidationEmail(email string, token string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", from, password, host)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: Email Validation\r\n\r\nDear User,\r\nThank you for registering with our service. To activate your account, please use the following validation token:\r\nLink: https://tengkurafi.aenzt.tech/api/v1/register/validate?email=%s&token=%s\r\nPlease enter this token on the validation page to complete your registration\r\n\r\nThank you for your cooperation\r\nQuill", from, email, email, token)

	err := smtp.SendMail(host+":"+port, auth, from, []string{email}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
