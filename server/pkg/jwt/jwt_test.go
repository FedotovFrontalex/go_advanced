package jwt_test

import (
	"server/pkg/jwt"
	"testing"
)

func TestJWTCreate(t *testing.T) {
	const email = "test@test.ru"

	jwtService := jwt.NewJWT("RAlQlILi2vrm9tF+M5E+6SnoYOhvz3V+RkISzvy4vmY=")
	token, err := jwtService.Create(jwt.JWTData{
		Email: email,
	})

	if err != nil {
		t.Fatal(err)
	}

	isValid, data := jwtService.Parse(token)

	if !isValid {
		t.Fatal("Token is invalid")
	}

	if data.Email != email {
		t.Fatalf("Email %snot equal %s", data.Email, email)
	}
}
