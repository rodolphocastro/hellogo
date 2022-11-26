package main

import (
	"encoding/json"
	"fmt"
	cloudEvents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/eclipse/paho.mqtt.golang"
	"math/rand"
	"testing"
	"time"
)

const (
	// eventSource is the source of the Cloud Event
	eventSource = "github.com/rodolphocastro/hellogo"
	// eventType is the type of the Cloud Event
	eventType = "series.created"
)

// TvSeries holds data for a TV series.
type TvSeries struct {
	Name         string
	FirstAiredOn int64
}

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

// createCloudEvent creates a json CloudEvent from a TV Series.
func createCloudEvent(theBoys TvSeries) (event.Event, error) {
	newCloudEvent := cloudEvents.NewEvent()
	newCloudEvent.SetID(theBoys.Name)
	newCloudEvent.SetSource(eventSource)
	newCloudEvent.SetType(eventType)

	// Act
	err := newCloudEvent.SetData(cloudEvents.ApplicationJSON, theBoys)
	return newCloudEvent, err
}

// Gets the MQTT address from Minikube.
func getMqttAddress() string {
	return fmt.Sprintf("tcp://%v:1883", getMinikubeIp())
}

// Sets up the Environment for these tests.
func setupTestEnvironment(t *testing.T) {
	SkipTestIfMinikubeIsUnavailable(t)

	SpinUpK8s(t, pathToMQTT, time.Second*2)
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

func TestPublishAndSubscribeToTopic(t *testing.T) {
	// Arrange
	expected := getRandomMessage()
	got := ""
	setupTestEnvironment(t)
	client := createMqqtClient(t)
	onMessageReceived := func(client mqtt.Client, message mqtt.Message) {
		t.Log("Received a new message")
		got = string(message.Payload())
	}

	// Act
	client.Subscribe(topicName, 0, onMessageReceived)
	publishToken := client.Publish(topicName, 0, true, expected)
	publishToken.Wait()
	err := publishToken.Error()
	time.Sleep(time.Second / 2) // Waiting for a second to give MQTT some time

	// Assert
	if err != nil {
		t.Errorf("Expected no errors but found %v", err)
	}

	if got != expected {
		t.Errorf("Expected %v but found %v", expected, got)
	}

	CleanUpK8s(t, pathToMQTT)
}

func TestCreateACloudEvent(t *testing.T) {
	// Arrange
	theBoys := TvSeries{
		Name:         "The Boys",
		FirstAiredOn: time.Date(2019, time.July, 26, 0, 0, 0, 0, time.Local).Unix(),
	}

	// Act
	_, err := createCloudEvent(theBoys)

	// Assert
	if err != nil {
		t.Errorf("Expected no errors but found %v", err)
	}
}

func TestSerializeAndDeserializeCloudEvent(t *testing.T) {
	// Arrange
	var gotData TvSeries
	gotEvent := cloudEvents.NewEvent()
	theBoys := TvSeries{
		Name:         "The Boys",
		FirstAiredOn: time.Date(2019, time.July, 26, 0, 0, 0, 0, time.Local).Unix(),
	}
	subject, _ := createCloudEvent(theBoys)

	// Act
	gotBytes, err := json.Marshal(subject)
	if err != nil {
		t.Errorf("Expected no errors marshalling but found %v", err)
	}
	err = json.Unmarshal(gotBytes, &gotEvent)

	// Assert
	if err != nil {
		t.Errorf("Expected no errors unmarshalling but found %v", err)
	}

	err = subject.DataAs(&gotData)
	if err != nil {
		t.Errorf("Expected no error fecthing Data, but found %v", err)
	}

	if gotData != theBoys {
		t.Errorf("Expected %v but found %v", theBoys, gotData)
	}
}

func TestPublishAndSubscribeToACloudEvent(t *testing.T) {
	// Arrange
	var got TvSeries
	expected := TvSeries{
		Name:         "The Boys",
		FirstAiredOn: time.Date(2019, time.July, 26, 0, 0, 0, 0, time.Local).Unix(),
	}
	expectedCloudEvent, _ := createCloudEvent(expected)
	expectedJson, err := json.Marshal(expectedCloudEvent)
	if err != nil {
		t.Errorf("Error while arranging test: %v", err)
	}
	setupTestEnvironment(t)

	client := createMqqtClient(t)
	onMessageReceived := func(client mqtt.Client, message mqtt.Message) {
		received := cloudEvents.NewEvent()
		err2 := json.Unmarshal(message.Payload(), &received)
		if err2 != nil {
			t.Errorf("Expected no errors but found %v", err2)
		}
		err2 = received.DataAs(&got)
		if err2 != nil {
			t.Errorf("Expected no errors retriving CloudEvent.Data, but found: %v", err2)
		}
	}

	// Act
	client.Subscribe(topicName, 0, onMessageReceived)
	publishToken := client.Publish(topicName, 0, true, expectedJson)
	publishToken.Wait()
	err = publishToken.Error()
	time.Sleep(time.Second / 2) // Waiting for a second to give MQTT some time

	// Assert
	if err != nil {
		t.Errorf("Expected no errors but found %v", err)
	}

	if got != expected {
		t.Errorf("Expected %v but found %v", expected, got)
	}

	CleanUpK8s(t, pathToMQTT)
}
