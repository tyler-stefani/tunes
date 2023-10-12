package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	m "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthMongoDB struct {
	client *m.Client
}

func NewAuthDB(uri string) (*AuthMongoDB, error) {
	client, err := m.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return &AuthMongoDB{}, err
	}
	return &AuthMongoDB{
		client,
	}, nil
}

func (db *AuthMongoDB) WriteRefreshToken(id string, expires time.Time) (bool, error) {
	rt := RefreshToken{
		id,
		expires,
	}

	coll := db.client.Database("auth").Collection("refresh_tokens")
	if _, err := coll.InsertOne(context.Background(), rt); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (db *AuthMongoDB) FindRefreshToken(id string) (RefreshToken, error) {
	rt := RefreshToken{}

	filter := bson.D{{Key: "_id", Value: id}}

	coll := db.client.Database("auth").Collection("refresh_tokens")
	if err := coll.FindOne(context.Background(), filter).Decode(&rt); err != nil {
		return RefreshToken{}, err
	} else {
		return rt, nil
	}
}
