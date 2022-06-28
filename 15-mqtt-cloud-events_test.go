package main

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"math/rand"
	"testing"
	"time"
)

const pathToMQTT = "./environments/development/mqtt.yml"
const topicName = "my-awesome-topic"
const aMessage = "Hello, take me to your leader"

// getRandomMessage gets a random message.
func getRandomMessage() string {
	return fmt.Sprintf("%v-%v", aMessage, rand.Int())
}

// createMqqtClient creates a MQTT client - if an error is found it'll error out the test.
func createMqqtClient(t *testing.T) mqtt.Client {
	options := mqtt.NewClientOptions()
	options.AddBroker(getMqttAddress())
	options.SetClientID("go-lang-mqtt-test")
	client := mqtt.NewClient(options)
	token := client.Connect()
	token.Wait()
	err := token.Error()
	if err != nil {
		t.Errorf("Expected no errors but found %v", err)
	}

	return client
}

// Gets the MQTT address from Minikube.
func getMqttAddress() string {
	return fmt.Sprintf("tcp://%v:1883", getMinikubeIp())
}

// Sets up the Environment for these tests.
func setupTestEnvironment(t *testing.T) {
	SkipTestIfCI(t)

	SpinUpK8s(t, pathToMQTT)
	time.Sleep(time.Second)
}

func TestMQTTSetup(t *testing.T) {
	setupTestEnvironment(t)

	CleanUpK8s(t, pathToMQTT)
}

func TestConnectToBroker(t *testing.T) {
	// Arrange
	setupTestEnvironment(t)

	// Act
	client := createMqqtClient(t)

	// Assert
	if client == nil {
		t.Error("Expected a MQTT Client but found nil")
	}
	CleanUpK8s(t, pathToMQTT)
}

func TestPublishToTopic(t *testing.T) {
	// Arrange
	setupTestEnvironment(t)
	client := createMqqtClient(t)

	// Act
	publishToken := client.Publish(topicName, 0, true, getRandomMessage())
	publishToken.Wait()
	err := publishToken.Error()

	// Assert
	if err != nil {
		t.Errorf("Expected no errors but found %v", err)
	}
	CleanUpK8s(t, pathToMQTT)
}
