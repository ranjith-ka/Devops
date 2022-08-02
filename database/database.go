package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const shortDuration = 3 * time.Second

type user struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name"`
	Email string             `bson:"email"`
}

func SetupDB() {
	d := time.Now().Add(shortDuration)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	clientOpts := options.Client().ApplyURI("mongodb://root:rootpass@localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	usersCollection := client.Database("testing").Collection("users")
	user := bson.D{{Key: "fullName", Value: "User 1"}, {Key: "age", Value: 30}}
	result, err := usersCollection.InsertOne(context.TODO(), user)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.InsertedID)
}

func insert(userData user) (*user, error) {
	insertedResult, err := userCollection.InsertOne(context.Background(), userData)
	if err != nil {
		return nil, err
	}

	userData.ID = insertedResult.InsertedID.(primitive.ObjectID)
	return &userData, nil
}

func insertMany(usersData []user) error {
	users := make([]interface{}, len(usersData))
	for i, userData := range usersData {
		users[i] = userData
	}

	if _, err := userCollection.InsertMany(context.Background(), users); err != nil {
		return err
	}
	return nil
}
