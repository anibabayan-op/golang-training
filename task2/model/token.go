package model

type Token struct {
	Token  string `bson:"token"`
	UserID string `bson:"user_id"`
	Exp    int64  `bson:"exp"`
}
