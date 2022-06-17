package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"os/exec"
	"testing"
)

const pathToK8s = "./k8s.yaml"
const ciEnvKey = "CI"
const minikubeIp = "192.168.49.2"
const mongoDbCredentials = "root:notsafe"

// Skips a test if the current environment is a CI pipeline.
func skipTestIfCI(t *testing.T) {
	if os.Getenv(ciEnvKey) != "" {
		t.Skip("Skipping this test - we're running in a CI environment")
	}
}

// Quick and Dirty way to spin up the deployment - invoking kubectl in the os' console.
func spinUpMongoK8s(t *testing.T) {
	kubeApply := exec.Command("kubectl", "apply", "-f", pathToK8s)
	if kubeApply.Run() != nil {
		t.Error("Error while spinning up MongoDb")
	}
}

// Quick and Dirty way to delete the deployment - invoking kubectl in the os' console.
func cleanUpMongoK8s(t *testing.T) {
	kubeDelete := exec.Command("kubectl", "delete", "-f", pathToK8s)
	if kubeDelete.Run() != nil {
		t.Error("Error while cleaning up MongoDb")
	}
}

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
	skipTestIfCI(t)

	spinUpMongoK8s(t)

	cleanUpMongoK8s(t)
}

// Attempt to connect to a mongodb instance
func TestMongoClient(t *testing.T) {
	skipTestIfCI(t)

	spinUpMongoK8s(t)

	client := createMongoClient(t)
	// Pinging the database to confirm we have a connection!
	err := client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		t.Errorf("Something went wrote while pinging: %v", err)
	}

	cleanUpMongoK8s(nil)
}
