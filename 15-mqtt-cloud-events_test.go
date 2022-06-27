package main

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"testing"
	"time"
)

const pathToMQTT = "./environments/development/mqtt.yml"
const topicName = "my-awesome-topic"

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

func TestPublishToATopic(t *testing.T) {
	// Arrange
	setupTestEnvironment(t)

	// Act
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

	// Assert
	CleanUpK8s(t, pathToMQTT)
}
