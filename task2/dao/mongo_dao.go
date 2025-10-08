package dao

import (
	"golang-training/task2/model"
)

// MongoDAO implements dao.DAO using package-level functions
type MongoDAO struct{}

// CreateUser inserts a new user
func (m *MongoDAO) CreateUser(user model.User) error {
	return CreateUser(user)
}

// GetUserByID fetches a user by ID
func (m *MongoDAO) GetUserByID(id string) (model.User, error) {
	return GetUserByID(id)
}

// GetUserByEmail fetches a user by email
func (m *MongoDAO) GetUserByEmail(email string) (model.User, error) {
	return GetUserByEmail(email)
}

// GetAllUsers fetches all users (matches dao.DAO interface)
func (m *MongoDAO) GetAllUsers() []model.User {
	return GetAllUsers()
}

// CreateToken generates and stores a JWT token for a user
func (m *MongoDAO) CreateToken(userID string) (string, error) {
	return CreateToken(userID)
}

// GetUserIDByToken validates a token and returns the associated userID
func (m *MongoDAO) GetUserIDByToken(token string) (string, error) {
	return GetUserIDByToken(token)
}

// DeleteToken invalidates a token
func (m *MongoDAO) DeleteToken(token string) error {
	return DeleteToken(token)
}
