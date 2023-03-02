package crypto

import (
	"os"
	"quell-api/entity"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(user entity.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return signedToken, err
}
