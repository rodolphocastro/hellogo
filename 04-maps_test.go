package main

import (
	"testing"
)

func TestMaps(t *testing.T) {
	const myKey = "awesome-key"
	const myOtherKey = "awful-key"
	myMap := make(map[string]int)

	// Adding items to a map
	myMap[myKey] = 11235
	myMap[myOtherKey] = 53211
	if myMap[myKey] != 11235 {
		t.Errorf("Expected %v but found %v", 11235, myMap[myKey])
	}

	// Removing items from a map
	delete(myMap, myKey)

	// Fetching items from a map
	_, isPresent := myMap[myKey]
	if isPresent {
		t.Error("Expected a problem to happen but everything went fine")
	}

	// Foreach item in a Map
	for _, item := range myMap {
		if item == 0 {
			t.Error("Expected anything save zero, but found zero")
		}
	}
}
