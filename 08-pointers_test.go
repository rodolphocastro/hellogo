package main

import (
	"testing"
)

// Just like C (which I don't miss lol) whenever you use a pointer you can mutate the object directly
// due to the memory address
func ptrMutateTo30(subject *int) {
	*subject = 30
}

func mutatoTo30(subject int) {
	subject = 30
}

func TestPointers(t *testing.T) {
	const defaultValue = 300
	myValue := defaultValue
	ptrMutateTo30(&myValue)
	if myValue == defaultValue {
		t.Errorf("expected myValue to be %v but found %v", 30, myValue)
	}

	myValue = defaultValue
	mutatoTo30(myValue)
	if myValue != defaultValue {
		t.Errorf("expected myValue to be %v, but found %v", defaultValue, myValue)
	}
}
