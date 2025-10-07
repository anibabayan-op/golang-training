package service

import (
	"errors"
	"golang-training/task2/dao"
	"golang-training/task2/model"
)

func RegisterUser(user model.User) error {
	return dao.CreateUser(user)
}

func AuthenticateUser(email, password string) (string, error) {
	user, err := dao.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if user.Password != password {
		return "", errors.New("invalid credentials")
	}
	token, err := dao.CreateToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func LogoutUser(token string) {
	dao.DeleteToken(token)
}

func GetCurrentUser(token string) (model.User, error) {
	userID, err := dao.GetUserIDByToken(token)
	if err != nil {
		return model.User{}, err
	}
	return dao.GetUserByID(userID)
}

func GetAllUsersService(token string) ([]model.User, error) {
	_, err := GetCurrentUser(token)
	if err != nil {
		return nil, err
	}
	return dao.GetAllUsers(), nil
}
