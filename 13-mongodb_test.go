package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"strconv"
	"testing"
)

// mongoLogger provides a single Logger instance for all Mongo Tests
var mongoLogger = InitializeLogger().With(zap.String("testSubject", "mongo"))

// mongoClient provides a single mongodb client for all Mongo Tests
var mongoClient *mongo.Client

const mongoDbCredentials = "root:notsafe"
const databaseName = "integration-tests"
const collectionName = "awesomeThings"
const pathToMongoK8s = "./environments/development/mongo.yaml"

// Creates a mongodb mqttClient for the integration environment.
func createMongoClient(t *testing.T) *mongo.Client {
	if mongoClient != nil {
		return mongoClient
	}

	mongoLogger.Info("a client isn't available, creating a new one")
	newClient, err := mongo.Connect(
		context.TODO(),
		options.Client().
			ApplyURI(fmt.Sprintf("mongodb://%v@%v:27017", mongoDbCredentials, getMinikubeIp())),
	)
	if err != nil {
		t.Errorf("Unable to connect to MongoDB: %v", err)
	}

	mongoClient = newClient
	return mongoClient
}

// Attempt to connect to a mongodb instance
func givenAnEnvironmentWhenAClientIsCreatedThenAPingShouldBePossible(t *testing.T) {
	client := createMongoClient(t)
	// Pinging the database to confirm we have a connection!
	err := client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		t.Errorf("Something went wrote while pinging: %v", err)
	}
}

// Access (or create) a Collection in the database
func givenAClientWhenACollectionIsFetchedThenNoErrorsShouldHappen(t *testing.T) {
	// Arrange
	client := createMongoClient(t)

	// Act
	collection := client.Database(databaseName).Collection(collectionName)

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
// Arrange
func givenACollectionWhenADocumentIsInsertedAndQueriedThenDataShouldBeRecovereable(t *testing.T) {
	anEntity := Book{
		Title:  "At the Mountains of Madness",
		Author: "H.P. Lovecraft",
		Tags:   []string{"Horror", "Lovecraftian"},
		// Act
	}
	collection := createMongoClient(t).Database(databaseName).Collection(collectionName)

	insertResult, err := collection.InsertOne(context.TODO(), anEntity)
	if err != nil {
		t.Errorf("Expected no errors but found: %v", err)
	}
	anEntity.ID = insertResult.InsertedID.(primitive.ObjectID) // A type assertion - checking if the interface{} is actually an ObjectID

	deleteResult, err := collection.DeleteOne(context.TODO(), bson.M{"_id": insertResult.InsertedID})
	if err != nil {
		t.Errorf("Expected no errors but found: %v", err)
	}

	// Assert
	if primitive.ObjectID.IsZero(anEntity.ID) {
		t.Errorf("Expected entity ID to be %v but found Empty", insertResult.InsertedID)
	}

	if deleteResult.DeletedCount != 1 {
		t.Errorf("Expected deleted count to be 1 but found %v", deleteResult.DeletedCount)
	}
}

func TestMongoDbScenarios(t *testing.T) {
	// Arrange
	SkipTestIfMinikubeIsUnavailable(t)
	scenarios := []func(*testing.T){
		givenAnEnvironmentWhenAClientIsCreatedThenAPingShouldBePossible,
		givenAClientWhenACollectionIsFetchedThenNoErrorsShouldHappen,
		givenACollectionWhenADocumentIsInsertedAndQueriedThenDataShouldBeRecovereable,
	}

	scenarioLogger := mongoLogger.
		With(
			zap.Int("totalTestCases", len(scenarios)),
		)

	scenarioLogger.Info("initializing mongo environment")
	SpinUpK8s(t, pathToMongoK8s)
	defer func() {
		scenarioLogger.Info("disconnecting the client")
		err := mongoClient.Disconnect(context.TODO())
		if err != nil {
			scenarioLogger.Error("unexpected error disconnecting", zap.Error(err))
		}
		scenarioLogger.Info("deleting the mongo environment")
		CleanUpK8s(t, pathToMongoK8s)
		scenarioLogger.Info("environment deleted")
	}()
	scenarioLogger.Info("environment initialized, executing tests")

	// Act and Assert
	for idx, scenario := range scenarios {
		scenarioLogger.Info("executing scenario", zap.Int("currentTest", idx+1))
		t.Run(strconv.Itoa(idx), scenario)
		scenarioLogger.Info("scenario completed", zap.Int("currentTest", idx+1))
	}
}
