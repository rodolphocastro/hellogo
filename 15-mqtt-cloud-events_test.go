package main

import (
	"encoding/json"
	"fmt"
	cloudEvents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

const (
	// eventSource is the source of the Cloud Event
	eventSource = "github.com/rodolphocastro/hellogo"
	// eventType is the type of the Cloud Event
	eventType = "series.created"
	// pathToMQTT is the source file for a k8s manifesto to spin up MQTT
	pathToMQTT           = "./environments/development/mqtt.yml"
	simpleTopicName      = "my-awesome-topic"
	cloudEventsTopicName = "cloudy-topic"
	aMessage             = "Hello, take me to your leader"
)

// TvSeries holds data for a TV series.
type TvSeries struct {
	Name         string
	FirstAiredOn int64
}

// mqttClient provides a single, shared, instance of a MQTT Client for integration testing MQTT
var mqttClient mqtt.Client

// getRandomMessage gets a random message.
func getRandomMessage() string {
	return fmt.Sprintf("%v-%v", aMessage, rand.Int())
}

// createMqqtClient creates a MQTT mqttClient - if an error is found it'll error out the test.
func createMqqtClient(t *testing.T) mqtt.Client {
	if mqttClient != nil {
		return mqttClient
	}

	options := mqtt.NewClientOptions()
	options.AddBroker(getMqttAddress())
	options.SetClientID("go-lang-mqtt-test")
	mqttClient = mqtt.NewClient(options)
	token := mqttClient.Connect()
	token.Wait()
	err := token.Error()
	if err != nil {
		t.Fatalf("Expected no errors but found %v", err)
	}
	return mqttClient
}

func TestMqttScenarios(t *testing.T) {
	// Arrange
	// a curated list of tests that need a complete MQTT environment
	testCases := []func(*testing.T){
		testCreateACloudEvent,
		testConnectToBroker,
		testPublishToTopic,
		testPublishAndSubscribeToTopic,
		testSerializeAndDeserializeCloudEvent,
		testPublishAndSubscribeToACloudEvent,
	}

	mqttLogger := InitializeLogger().
		With(zap.String("testSubject", "mqtt")).
		With(zap.Int("totalTestCases", len(testCases)))
	mqttLogger.Info("initializing mqtt environment")
	setupTestEnvironment(t)
	defer func() {
		mqttLogger.Info("disconnecting the client")
		mqttClient.Disconnect(1000)
		mqttClient = nil
		mqttLogger.Info("deleting the environment")
		CleanUpK8s(t, pathToMQTT)
	}()
	time.Sleep(time.Second)
	mqttLogger.Info("environment initialized, executing tests")

	// Act and Assert
	for idx, testCase := range testCases {
		mqttLogger.Info("running a test case", zap.Int("currentTest", idx+1))
		t.Run(strconv.Itoa(idx), testCase)
		mqttLogger.Info("test case completed", zap.Int("currentTest", idx+1))
	}
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
	SpinUpK8s(t, pathToMQTT, time.Second*4)
	time.Sleep(time.Second)
}

func testConnectToBroker(t *testing.T) {
	// Arrange

	// Act
	client := createMqqtClient(t)

	// Assert
	if client == nil {
		t.Error("expected a MQTT Client but found nil")
	}
}

func testPublishToTopic(t *testing.T) {
	// Arrange
	client := createMqqtClient(t)

	// Act
	publishToken := client.Publish(simpleTopicName, 0, true, getRandomMessage())
	publishToken.Wait()
	err := publishToken.Error()

	// Assert
	if err != nil {
		t.Errorf("Expected no errors but found %v", err)
	}
}

func testPublishAndSubscribeToTopic(t *testing.T) {
	// Arrange
	expected := getRandomMessage()
	got := ""
	client := createMqqtClient(t)
	onMessageReceived := func(client mqtt.Client, message mqtt.Message) {
		t.Log("Received a new message")
		got = string(message.Payload())
	}

	// Act
	client.Subscribe(simpleTopicName, 0, onMessageReceived)
	publishToken := client.Publish(simpleTopicName, 0, true, expected)
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
}

func testCreateACloudEvent(t *testing.T) {
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

func testSerializeAndDeserializeCloudEvent(t *testing.T) {
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

func testPublishAndSubscribeToACloudEvent(t *testing.T) {
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
	client.Subscribe(cloudEventsTopicName, 0, onMessageReceived)
	publishToken := client.Publish(cloudEventsTopicName, 0, true, expectedJson)
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
}
