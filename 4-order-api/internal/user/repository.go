package user

import (
	"orderApi/pkg/db"
)

type UserRepository struct {
	Database *db.Db
}

func NewUserRepository(db *db.Db) *UserRepository {
	return &UserRepository{
		Database: db,
	}
}

func (repo *UserRepository) Create(user *User) (*User, error) {
	result := repo.Database.DB.Create(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repo *UserRepository) Update(user *User) (*User, error) {
	result := repo.Database.DB.Updates(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repo *UserRepository) FindByPhone(phone string) (*User, error) {
	var user User
	result := repo.Database.DB.First(&user, "phone=?", phone)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (repo *UserRepository) FindBySessionId(sessionId string) (*User, error) {
	var user User
	result := repo.Database.DB.First(&user, "session_id=?", sessionId)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
