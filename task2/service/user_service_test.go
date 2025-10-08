package service

import (
	"testing"

	"golang-training/task2/dao"
	"golang-training/task2/model"

	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterUser(t *testing.T) {
	mockDAO := new(dao.MockDAO)
	service := NewUserService(mockDAO)

	input := model.User{ID: "1", Name: "Alice", Email: "alice@example.com", Password: "secret"}

	mockDAO.On("CreateUser", mock.AnythingOfType("model.User")).Return(nil)

	created, err := service.RegisterUser(input)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}
	if created.Email != "alice@example.com" {
		t.Fatalf("Expected email alice@example.com, got %s", created.Email)
	}

	mockDAO.AssertExpectations(t)
}

func TestAuthenticateUser(t *testing.T) {
	mockDAO := new(dao.MockDAO)
	service := NewUserService(mockDAO)

	// Generate a bcrypt hash for "secret"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	user := model.User{ID: "1", Name: "Alice", Email: "alice@example.com", Password: string(hashedPassword)}

	mockDAO.On("GetUserByEmail", "alice@example.com").Return(user, nil)
	mockDAO.On("CreateToken", "1").Return("mocked-token", nil)

	token, err := service.AuthenticateUser("alice@example.com", "secret")
	if err != nil {
		t.Fatalf("AuthenticateUser failed: %v", err)
	}
	if token != "mocked-token" {
		t.Fatalf("Expected mocked-token, got %s", token)
	}

	mockDAO.AssertExpectations(t)
}

func TestLogoutUser(t *testing.T) {
	mockDAO := new(dao.MockDAO)
	service := NewUserService(mockDAO)

	mockDAO.On("DeleteToken", "token123").Return(nil)

	service.LogoutUser("token123")

	mockDAO.AssertCalled(t, "DeleteToken", "token123")
	mockDAO.AssertExpectations(t)
}

func TestGetCurrentUser(t *testing.T) {
	mockDAO := new(dao.MockDAO)
	service := NewUserService(mockDAO)

	user := model.User{ID: "1", Name: "Alice", Email: "alice@example.com"}

	mockDAO.On("GetUserIDByToken", "token123").Return("1", nil)
	mockDAO.On("GetUserByID", "1").Return(user, nil)

	u, err := service.GetCurrentUser("token123")
	if err != nil {
		t.Fatalf("GetCurrentUser failed: %v", err)
	}
	if u.ID != "1" {
		t.Fatalf("Expected user ID 1, got %s", u.ID)
	}

	mockDAO.AssertExpectations(t)
}

func TestGetAllUsersService(t *testing.T) {
	mockDAO := new(dao.MockDAO)
	service := NewUserService(mockDAO)

	u1 := model.User{ID: "1", Name: "Alice", Email: "alice@example.com"}
	u2 := model.User{ID: "2", Name: "Bob", Email: "bob@example.com"}

	// Mock token -> userID
	mockDAO.On("GetUserIDByToken", "token123").Return("1", nil)
	mockDAO.On("GetUserByID", "1").Return(u1, nil) // needed for GetCurrentUser
	mockDAO.On("GetAllUsers").Return([]model.User{u1, u2})

	users, err := service.GetAllUsersService()
	if err != nil {
		t.Fatalf("GetAllUsersService failed: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("Expected 2 users, got %d", len(users))
	}

	mockDAO.AssertExpectations(t)
}
