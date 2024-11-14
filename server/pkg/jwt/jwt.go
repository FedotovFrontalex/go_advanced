package jwt

import (
	"server/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	Secret string
}

type JWTData struct {
	Email string
}

func NewJWT(secret string) *JWT {
	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) Create(data JWTData) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": data.Email,
	})
	logger.Log("Secret: ", j.Secret)
	s, err := token.SignedString([]byte(j.Secret))

	if err != nil {
		return "", err
	}

	return s, nil
}
