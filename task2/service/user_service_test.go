package service

import (
	"testing"

	"golang-training/task2/dao"
	"golang-training/task2/model"
)

func resetDAO() {
	dao.Users = make(map[string]model.User)
	dao.Tokens = make(map[string]string)
}

func TestRegisterUser(t *testing.T) {
	resetDAO()

	user := model.User{ID: "1", Name: "Alice", Email: "alice@example.com", Password: "secret"}
	err := RegisterUser(user)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	err = RegisterUser(user)
	if err == nil {
		t.Fatalf("Expected error when registering duplicate user")
	}
}

func TestAuthenticateUser(t *testing.T) {
	resetDAO()

	user := model.User{ID: "1", Name: "Alice", Email: "alice@example.com", Password: "secret"}
	_ = dao.CreateUser(user)

	token, err := AuthenticateUser("alice@example.com", "secret")
	if err != nil {
		t.Fatalf("AuthenticateUser failed: %v", err)
	}
	if token == "" {
		t.Fatalf("Expected a token, got empty string")
	}

	_, err = AuthenticateUser("alice@example.com", "wrong")
	if err == nil {
		t.Fatalf("Expected error for wrong password")
	}

	_, err = AuthenticateUser("bob@example.com", "secret")
	if err == nil {
		t.Fatalf("Expected error for non-existing email")
	}
}

func TestLogoutUser(t *testing.T) {
	resetDAO()

	user := model.User{ID: "1", Name: "Alice", Email: "alice@example.com", Password: "secret"}
	_ = dao.CreateUser(user)
	token, _ := dao.CreateToken("1")

	LogoutUser(token)

	_, err := dao.GetUserIDByToken(token)
	if err == nil {
		t.Fatalf("Expected error after logout, token should be deleted")
	}
}

func TestGetCurrentUser(t *testing.T) {
	resetDAO()

	user := model.User{ID: "1", Name: "Alice", Email: "alice@example.com", Password: "secret"}
	_ = dao.CreateUser(user)
	token, _ := dao.CreateToken("1")

	u, err := GetCurrentUser(token)
	if err != nil {
		t.Fatalf("GetCurrentUser failed: %v", err)
	}
	if u.ID != "1" {
		t.Fatalf("Expected user ID 1, got %s", u.ID)
	}

	_, err = GetCurrentUser("invalid")
	if err == nil {
		t.Fatalf("Expected error for invalid token")
	}
}

func TestGetAllUsersService(t *testing.T) {
	resetDAO()

	u1 := model.User{ID: "1", Name: "Alice", Email: "alice@example.com", Password: "secret"}
	u2 := model.User{ID: "2", Name: "Bob", Email: "bob@example.com", Password: "pass"}
	_ = dao.CreateUser(u1)
	_ = dao.CreateUser(u2)
	token, _ := dao.CreateToken("1")

	users, err := GetAllUsersService(token)
	if err != nil {
		t.Fatalf("GetAllUsersService failed: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("Expected 2 users, got %d", len(users))
	}

	_, err = GetAllUsersService("invalid")
	if err == nil {
		t.Fatalf("Expected error for invalid token")
	}
}
