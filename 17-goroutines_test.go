package main

import (
	"go.uber.org/zap"
	"testing"
	"time"
)

func initializeZap() *zap.Logger {
	logger, err := zap.NewDevelopment()

	if err != nil {
		panic("Error while setting up zap")
	}

	return logger
}

func doubleSomething(waitTime time.Duration, subject int) int {
	time.Sleep(waitTime)
	return subject * 2
}

func TestTheGoKeyword(t *testing.T) {
	// Arrange
	var firstGot = 0
	var secondGot = 0
	logger := initializeZap()
	firstExpected := 8
	secondExpected := 66
	firstInput := firstExpected / 2
	secondInput := secondExpected / 2

	// Act
	// the 'go' keyword means 'fire this and forget about it', which is pretty useful to allow the runtime to deal
	// with spawning async work
	go func() {
		logger.Info("First routine launched")
		firstGot = doubleSomething(0, firstInput)
		logger.Info("First routine completed")
	}()
	// giving it some time, just in case
	time.Sleep(time.Millisecond)

	// this one should take at least a second to run
	go func() {
		logger.Info("Second routine launched")
		secondGot = doubleSomething(time.Second, secondInput)
		logger.Info("Second routine launched")
	}()

	// Assert
	if firstGot != firstExpected {
		t.Errorf("Expected %d but found %d instead", firstExpected, firstGot)
	}

	if secondGot == secondExpected {
		t.Errorf("Didn't expect %d, but found it", secondGot)
	}

	// lets give it some more time for the function to return
	time.Sleep(time.Second)
	if secondGot != secondExpected {
		t.Errorf("Expected %d but found %d instead", secondExpected, secondGot)
	}
}
