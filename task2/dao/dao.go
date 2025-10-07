package dao

import (
	"errors"
	"golang-training/task2/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	Users     = make(map[string]model.User)
	Tokens    = make(map[string]string)
	jwtSecret = []byte("supersecretkey")
)

func CreateUser(user model.User) error {
	if _, exists := Users[user.ID]; exists {
		return errors.New("user already exists")
	}
	Users[user.ID] = user
	return nil
}

func GetUserByID(id string) (model.User, error) {
	user, exists := Users[id]
	if !exists {
		return model.User{}, errors.New("user not found")
	}
	return user, nil
}

func GetAllUsers() []model.User {
	var all []model.User
	for _, u := range Users {
		all = append(all, u)
	}
	return all
}

func GetUserByEmail(email string) (model.User, error) {
	for _, u := range Users {
		if u.Email == email {
			return u, nil
		}
	}
	return model.User{}, errors.New("user not found")
}

func generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func CreateToken(userID string) (string, error) {
	token, err := generateToken(userID)
	if err != nil {
		return "", err
	}
	Tokens[token] = userID
	return token, nil
}

func GetUserIDByToken(token string) (string, error) {
	userID, ok := Tokens[token]
	if !ok {
		return "", errors.New("invalid token")
	}
	return userID, nil
}

func DeleteToken(token string) {
	delete(Tokens, token)
}
