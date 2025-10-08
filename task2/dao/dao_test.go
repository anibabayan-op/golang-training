package dao

import (
	"context"
	"testing"
	"time"

	"golang-training/task2/model"

	"go.mongodb.org/mongo-driver/bson"
)

func clearCollections() {
	_, _ = userCol.DeleteMany(context.Background(), bson.M{})
	_, _ = tokenCol.DeleteMany(context.Background(), bson.M{})
}

func TestCreateAndGetUser(t *testing.T) {
	clearCollections()

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

	// Duplicate should fail
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
	clearCollections()

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
	clearCollections()

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

	err = DeleteToken(token)
	if err != nil {
		t.Fatalf("DeleteToken failed: %v", err)
	}

	_, err = GetUserIDByToken(token)
	if err == nil {
		t.Fatalf("Expected error after deleting token")
	}
}

func TestTokenExpirationSimulation(t *testing.T) {
	clearCollections()

	userID := "2"
	exp := time.Now().Add(-1 * time.Minute).Unix() // expired token

	// Create expired token in DB
	expiredToken := model.Token{
		Token:  "expiredtoken123",
		UserID: userID,
		Exp:    exp,
	}
	_, err := tokenCol.InsertOne(context.Background(), expiredToken)
	if err != nil {
		t.Fatalf("Failed to insert expired token: %v", err)
	}

	_, err = GetUserIDByToken("expiredtoken123")
	if err == nil {
		t.Fatalf("Expected error for expired token")
	}
}
