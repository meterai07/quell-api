package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"quell-api/entity"
	"quell-api/initializers"
	"quell-api/sdk/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		response.Response(c, http.StatusUnauthorized, "Unauthorized", nil)
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		response.Response(c, http.StatusUnauthorized, "Error when parsing the token", nil)
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			response.Response(c, http.StatusUnauthorized, "Token expired", nil)
			c.Abort()
			return
		}

		var user entity.User

		if err := initializers.DB.Where("id = ?", claims["sub"]).First(&user).Error; err != nil {
			response.Response(c, http.StatusUnauthorized, "Unauthorized", nil)
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	} else {
		c.Abort()
	}
}
