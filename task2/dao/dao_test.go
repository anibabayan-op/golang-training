package dao

import (
	"testing"

	"golang-training/task2/model"

	"github.com/golang-jwt/jwt/v5"
)

func TestCreateAndGetUser(t *testing.T) {
	Users = make(map[string]model.User)

	user := model.User{
		ID:       "1",
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "secret",
	}

	err := CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	err = CreateUser(user)
	if err == nil {
		t.Fatalf("Expected error when creating duplicate user")
	}

	u, err := GetUserByID("1")
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if u.Email != "alice@example.com" {
		t.Fatalf("Expected email alice@example.com, got %s", u.Email)
	}

	u2, err := GetUserByEmail("alice@example.com")
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}
	if u2.ID != "1" {
		t.Fatalf("Expected ID 1, got %s", u2.ID)
	}

	_, err = GetUserByID("999")
	if err == nil {
		t.Fatalf("Expected error for non-existing user ID")
	}
}

func TestGetAllUsers(t *testing.T) {
	Users = make(map[string]model.User)

	u1 := model.User{ID: "1", Name: "Alice", Email: "a@example.com", Password: "p"}
	u2 := model.User{ID: "2", Name: "Bob", Email: "b@example.com", Password: "p"}
	_ = CreateUser(u1)
	_ = CreateUser(u2)

	all := GetAllUsers()
	if len(all) != 2 {
		t.Fatalf("Expected 2 Users, got %d", len(all))
	}
}

func TestTokenOperations(t *testing.T) {
	Tokens = make(map[string]string)

	userID := "1"

	token, err := CreateToken(userID)
	if err != nil {
		t.Fatalf("CreateToken failed: %v", err)
	}

	uid, err := GetUserIDByToken(token)
	if err != nil {
		t.Fatalf("GetUserIDByToken failed: %v", err)
	}
	if uid != userID {
		t.Fatalf("Expected userID %s, got %s", userID, uid)
	}

	DeleteToken(token)
	_, err = GetUserIDByToken(token)
	if err == nil {
		t.Fatalf("Expected error after deleting token")
	}
}

func TestTokenExpirationSimulation(t *testing.T) {
	Tokens = make(map[string]string)

	userID := "2"
	token, err := generateToken(userID)
	if err != nil {
		t.Fatalf("generateToken failed: %v", err)
	}

	Tokens[token] = userID

	uid, err := GetUserIDByToken(token)
	if err != nil {
		t.Fatalf("Expected token to exist, got error: %v", err)
	}
	if uid != userID {
		t.Fatalf("Expected userID %s, got %s", userID, uid)
	}

	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		t.Fatalf("Failed to parse JWT: %v", err)
	}
	if !parsed.Valid {
		t.Fatalf("Token should be valid immediately after creation")
	}
}
