package dao

import (
	"context"
	"errors"
	"golang-training/task2/model"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	client    *mongo.Client
	userCol   *mongo.Collection
	tokenCol  *mongo.Collection
	jwtSecret []byte
	ctx       context.Context
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not loaded, relying on environment variables")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET not set")
	}
	jwtSecret = []byte(secret)

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	userCol = client.Database("task2").Collection("users")
	tokenCol = client.Database("task2").Collection("tokens")
}

func CreateUser(user model.User) error {
	filter := bson.M{"email": user.Email}
	count, err := userCol.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("user already exists")
	}

	_, err = userCol.InsertOne(ctx, user)
	return err
}

func GetUserByID(id string) (model.User, error) {
	var user model.User
	filter := bson.M{"id": id}
	err := userCol.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, errors.New("user not found")
		}
		return model.User{}, err
	}
	return user, nil
}

func GetAllUsers() []model.User {
	var users []model.User
	cursor, err := userCol.Find(ctx, bson.M{})
	if err != nil {
		log.Println("GetAllUsers error:", err)
		return users
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var u model.User
		if err := cursor.Decode(&u); err != nil {
			log.Println("Decode user error:", err)
			continue
		}
		users = append(users, u)
	}
	return users
}

func GetUserByEmail(email string) (model.User, error) {
	var user model.User
	filter := bson.M{"email": email}
	err := userCol.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, errors.New("user not found")
		}
		return model.User{}, err
	}

	log.Printf("GetUserByEmail: found user ID=%s, Name=%s, Email=%s, PasswordHash=%s\n",
		user.ID, user.Name, user.Email, user.Password)
	return user, nil
}

func CreateToken(userID string) (string, error) {
	exp := time.Now().Add(time.Hour).Unix()

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": exp,
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tokenObj.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	_, err = tokenCol.InsertOne(ctx, model.Token{
		Token:  tokenString,
		UserID: userID,
		Exp:    exp,
	})
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserIDByToken(token string) (string, error) {
	var t model.Token
	err := tokenCol.FindOne(ctx, bson.M{"token": token}).Decode(&t)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", errors.New("invalid token")
		}
		return "", err
	}

	if time.Now().Unix() > t.Exp {
		return "", errors.New("token expired")
	}

	return t.UserID, nil
}

func DeleteToken(token string) error {
	_, err := tokenCol.DeleteOne(ctx, bson.M{"token": token})
	return err
}
