package dao

import "golang-training/task2/model"

type DAO interface {
	CreateUser(user model.User) error
	GetUserByID(id string) (model.User, error)
	GetUserByEmail(email string) (model.User, error)
	GetAllUsers() []model.User
	CreateToken(userID string) (string, error)
	GetUserIDByToken(token string) (string, error)
	DeleteToken(token string) error
}
