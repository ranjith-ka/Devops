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

var userCollection *mongo.Collection

type user struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name"`
	Email string             `bson:"email"`
}

// SetupDB function to set context and ping status, will close the connection after certiain amount of time
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

	collection := client.Database("testing").Collection("users")

	Firstuser := bson.D{{Key: "fullName", Value: "User_1"}, {Key: "age", Value: 30}}
	result, err := collection.InsertOne(context.Background(), Firstuser)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.InsertedID)

	userData := user{
		Name:  "ranjith",
		Email: "ok@gmail.com",
	}

	userCollection = collection
	t, err := insert(userData)
	fmt.Println(t.ID)
}

// insert func to support single document in the collection, but test case are take over
func insert(userData user) (*user, error) {

	insertedResult, err := userCollection.InsertOne(context.Background(), userData)
	if err != nil {
		return nil, err
	}
	userData.ID = insertedResult.InsertedID.(primitive.ObjectID)
	return &userData, nil
}

// func insertMany(usersData []user) error {
// 	users := make([]interface{}, len(usersData))
// 	for i, userData := range usersData {
// 		users[i] = userData
// 	}

// 	if _, err := userCollection.InsertMany(context.Background(), users); err != nil {
// 		return err
// 	}
// 	return nil
// }
