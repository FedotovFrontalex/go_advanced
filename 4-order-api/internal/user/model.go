package user

import (
	"math/rand/v2"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	SessionId string `json:"sessionId" gorm:"uniqIndex"`
	Phone     string `json:"phone" gorm:"uniqIndex"`
}

func NewUser(phone string) *User {
	return &User{
		Phone: phone,
	}
}

func (user *User) CreateSessionId(n int) {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*")
	sessionId := make([]byte, n)

	for i := 0; i < n; i++ {
		letter := letters[rand.IntN(len(letters))]
		sessionId[i] = byte(letter)
	}

	user.SessionId = string(sessionId)
}
