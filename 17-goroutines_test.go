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

// The "go" keyword allows us to spin off goroutines from anywhere in the code.
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

// Channels are used to allow goroutines to communicate with one another or with a higher level function
func TestChannelsAllowGoroutinesToCommunicate(t *testing.T) {
	// Arrange
	logger := initializeZap()
	expected := 8
	// the 'chan' keyword indicates we want to create a Channel to allow communication to happen
	results := make(chan int)

	// Act
	// firing a goroutine that outputs to the results
	go func() {
		logger.Info("Starting a goroutine")
		// the right hand side (rhs) <- operator publishes something to a channel
		logger.Info("Expected value has been published")
		results <- doubleSomething(time.Second, expected/2)
		logger.Info("An unexpected value has been published now")
		results <- 0
	}()

	// the left hand side (lhs) <- operator means to read something from the channel
	// note: both lhs and rhs are *blocking*. Which means that the routine will halt execution until something is read
	// and the receiving routine will also halt until something is received
	logger.Info("Waiting for got to be published!")
	got := <-results
	logger.Info("Got has been updated", zap.Int("got", got))

	// Assert
	// since <- blocks the reader and the publisher this should just work
	if got != expected {
		t.Errorf("Expected %d but found %d!", expected, got)
	}

	// now we're asking for a second value!
	got = <-results
	logger.Info("Got has been updated", zap.Int("got", got))
	if got == expected {
		t.Errorf("Expected anything other than %d, but found %d", expected, expected)
	}
}
