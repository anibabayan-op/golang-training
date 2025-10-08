package service

import (
	"errors"
	"golang-training/task2/dao"
	"golang-training/task2/model"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	dao dao.DAO
}

func NewUserService(d dao.DAO) *UserService {
	return &UserService{dao: d}
}

func (s *UserService) RegisterUser(input model.User) (model.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.UserResponse{}, err
	}

	input.Password = string(hashedPassword)
	if err := s.dao.CreateUser(input); err != nil {
		return model.UserResponse{}, err
	}

	user := model.UserResponse{
		ID:    input.ID,
		Name:  input.Name,
		Email: input.Email,
	}
	return user, nil
}

func (s *UserService) AuthenticateUser(email, password string) (string, error) {
	user, err := s.dao.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := s.dao.CreateToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *UserService) LogoutUser(token string) {
	_ = s.dao.DeleteToken(token)
}

func (s *UserService) GetUserByID(id string) (model.User, error) {
	return s.dao.GetUserByID(id)
}

func (s *UserService) GetCurrentUser(token string) (model.User, error) {
	userID, err := s.dao.GetUserIDByToken(token)
	if err != nil {
		return model.User{}, err
	}
	return s.dao.GetUserByID(userID)
}

func (s *UserService) GetAllUsersService() ([]model.User, error) {
	return s.dao.GetAllUsers(), nil
}
