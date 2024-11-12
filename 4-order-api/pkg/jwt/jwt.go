package jwt

import "github.com/golang-jwt/jwt/v5"

type JWT struct {
	Secret string
}

func NewJWT(secret string) *JWT {
	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) Create(sessionId string, phone string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phone":     phone,
		"sessionId": sessionId,
	})

	s, err := token.SignedString([]byte(j.Secret))

	if err != nil {
		return "", err
	}

	return s, nil
}
