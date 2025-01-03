package jwt

import "github.com/golang-jwt/jwt/v5"

type JWT struct {
	Secret string
}

type JWTData struct {
	Phone     string
	SessionId string
}

func NewJWT(secret string) *JWT {
	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) Create(data JWTData) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phone":     data.Phone,
		"sessionId": data.SessionId,
	})

	s, err := token.SignedString([]byte(j.Secret))

	if err != nil {
		return "", err
	}

	return s, nil
}

func (j *JWT) Parse(token string) (bool, *JWTData) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})

	if err != nil {
		return false, nil
	}

	phone := t.Claims.(jwt.MapClaims)["phone"]
	sessionId := t.Claims.(jwt.MapClaims)["sessionId"]

	return t.Valid, &JWTData{
		Phone:     phone.(string),
		SessionId: sessionId.(string),
	}
}
