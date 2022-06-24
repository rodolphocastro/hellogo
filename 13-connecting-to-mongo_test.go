package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
)

const minikubeIp = "192.168.49.2"
const mongoDbCredentials = "root:notsafe"
const databaseName = "integration-tests"
const collectionName = "awesomeThings"
const pathToMongoK8s = "./environments/development/mongo.yaml"

// Creates a mongodb client for the integration environment.
func createMongoClient(t *testing.T) *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%v@%v:27017", mongoDbCredentials, minikubeIp)))
	if err != nil {
		t.Errorf("Unable to connect to MongoDB: %v", err)
	}
	return client
}

// Attempt to create and tear down a k8s deployment.
func TestMongoSetup(t *testing.T) {
	SkipTestIfCI(t)

	SpinUpK8s(t, pathToMongoK8s)

	CleanUpK8s(t, pathToMongoK8s)
}

// Attempt to connect to a mongodb instance
func TestMongoClient(t *testing.T) {
	SkipTestIfCI(t)

	SpinUpK8s(t, pathToMongoK8s)

	client := createMongoClient(t)
	// Pinging the database to confirm we have a connection!
	err := client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		t.Errorf("Something went wrote while pinging: %v", err)
	}

	CleanUpK8s(t, pathToMongoK8s)
}

// Access (or create) a Collection in the database
func TestAccessACollection(t *testing.T) {
	// Arrange
	SkipTestIfCI(t)
	SpinUpK8s(t, pathToMongoK8s)
	client := createMongoClient(t)

	// Act
	collection := client.Database(databaseName).Collection(collectionName)
	CleanUpK8s(t, pathToMongoK8s)

	// Assert
	if collection == nil {
		t.Error("Didn't expect collection to be nil, but found nil")
	}
}

// A book!
type Book struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Title  string             `bson:"title,omitempty"`
	Author string             `bson:"author,omitempty"`
	Tags   []string           `bson:"tags,omitempty"`
}

// Creating and deleting a document within a MongoDB Collection
func TestInsertAndDeleteDocument(t *testing.T) {
	// Arrange
	SkipTestIfCI(t)
	SpinUpK8s(t, pathToMongoK8s)
	anEntity := Book{
		Title:  "At the Mountains of Madness",
		Author: "H.P. Lovecraft",
		Tags:   []string{"Horror", "Lovecraftian"},
	}
	collection := createMongoClient(t).Database(databaseName).Collection(collectionName)

	// Act
	insertResult, err := collection.InsertOne(context.TODO(), anEntity)
	if err != nil {
		t.Errorf("Expected no errors but found: %v", err)
	}
	anEntity.ID = insertResult.InsertedID.(primitive.ObjectID) // A type assertion - checking if the interface{} is actually an ObjectID

	deleteResult, err := collection.DeleteOne(context.TODO(), bson.M{"_id": insertResult.InsertedID})
	if err != nil {
		t.Errorf("Expected no errors but found: %v", err)
	}

	CleanUpK8s(t, pathToMongoK8s)

	// Assert
	if primitive.ObjectID.IsZero(anEntity.ID) {
		t.Errorf("Expected entity ID to be %v but found Empty", insertResult.InsertedID)
	}

	if deleteResult.DeletedCount != 1 {
		t.Errorf("Expected deleted count to be 1 but found %v", deleteResult.DeletedCount)
	}
}
