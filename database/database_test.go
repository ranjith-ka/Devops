package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestInsertOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		userCollection = mt.Coll
		id := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		insertedUser, err := insert(user{
			ID:    id,
			Name:  "john",
			Email: "john.doe@test.com",
		})
		assert.Nil(t, err)
		assert.Equal(t, &user{
			ID:    id,
			Name:  "john",
			Email: "john.doe@test.com",
		}, insertedUser)
	})

	mt.Run("custom error duplicate", func(mt *mtest.T) {
		userCollection = mt.Coll
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))

		insertedUser, err := insert(user{})

		assert.Nil(t, insertedUser)
		assert.NotNil(t, err)
		assert.True(t, mongo.IsDuplicateKeyError(err))
	})

	mt.Run("simple error", func(mt *mtest.T) {
		userCollection = mt.Coll
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})

		insertedUser, err := insert(user{})

		assert.Nil(t, insertedUser)
		assert.NotNil(t, err)
	})
}

//Add testcase for many
