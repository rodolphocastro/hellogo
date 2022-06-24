package main

import "testing"

const pathToMQTT = "./environments/development/mqtt.yml"

func TestMQTTSetup(t *testing.T) {
	SkipTestIfCI(t)

	SpinUpK8s(t, pathToMQTT)

	CleanUpK8s(t, pathToMQTT)
}
