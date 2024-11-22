package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/configs"
	"server/internal/auth"
	"server/internal/user"
	"server/pkg/db"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func bootstrap(t *testing.T) (*auth.AuthHandler, sqlmock.Sqlmock, error) {
	database, mock, err := sqlmock.New()

	if err != nil {
		return nil, nil, err
	}

	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: database,
	}))

	if err != nil {
		return nil, nil, err
	}

	userRepository := user.NewUserRepository(&db.Db{
		DB: gormDb,
	})

	handler := &auth.AuthHandler{
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				Secret: "AlQlILi2vrm9tF+M5E+6SnoYOhvz3V+RkISzvy4vmY=",
			},
		},
		AuthService: auth.NewAuthService(userRepository),
	}

	return handler, mock, nil
}

func TestLoginSuccess(t *testing.T) {
	handler, mock, err := bootstrap(t)

	rows := sqlmock.NewRows([]string{"email", "password"}).
		AddRow("test@test.ru", "$2a$10$G0MfRZoTTUK8r.NVjUialOqVIPvveK24waYSMt7WJEGeEXKXAUdc6")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	if err != nil {
		t.Fatal(err)
		return
	}

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@test.ru",
		Password: "Byntuhfk",
	})
	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", reader)

	handler.Login()(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, expected %d", w.Code, 200)
	}
}

func TestRegisterSuccess(t *testing.T) {
	handler, mock, err := bootstrap(t)

	rows := sqlmock.NewRows([]string{"email", "password", "name"})

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	mock.ExpectBegin()
	insertRows := sqlmock.NewRows([]string{"email", "password", "name"})
	insertRows.AddRow("test@test.ru", "Byntuhfk", "Alex")
	mock.ExpectQuery("INSERT").WillReturnRows(insertRows)
	mock.ExpectCommit()

	if err != nil {
		t.Fatal(err)
		return
	}

	data, _ := json.Marshal(&auth.RegisterRequest{
		Email:    "test@test.ru",
		Password: "Byntuhfk",
		Name:     "Alex",
	})
	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/register", reader)

	handler.Register()(w, req)

	if w.Code != 201 {
		t.Errorf("got %d, expected %d", w.Code, 201)
	}
}
