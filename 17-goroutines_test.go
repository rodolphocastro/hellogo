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

// Double something and publish it to a channel
func doubleSomethingWithAChannel(waitTime time.Duration, subject int, output chan int, onComplete func()) {
	output <- doubleSomething(waitTime, subject)
	onComplete()
}

// Outputs a message into a channel
// Note the <- after the channel argument, that means this channel can only be written to and not read from
func ping(message string, output chan<- string) {
	output <- message
}

// Subscribes to a channel and does something upon completion
// Note the <- next to the channel argument, that means this channel can only be read from and not written to
func pong(input <-chan string, onComplete func(msg string)) {
	got := <-input
	onComplete(got)
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

// Channels with a buffer allow goroutines to output more than a single value before blocking!
func TestBufferedChannelsBlockLessOften(t *testing.T) {
	// Arrange
	logger := initializeZap()
	expected := 8
	// the 'chan' keyword indicates we want to create a Channel to allow communication to happen
	// but this time we're also making a buffer of 2 which will allow a goroutine to output 3 times before it gets blocked!
	results := make(chan int, 2)

	// Act
	go func() {
		logger.Info("Firing off a goroutine!")
		results <- 0
		logger.Info("Published 0")
		results <- -1
		logger.Info("Published -1")
		logger.Info("About to publish the expected value!")
		results <- doubleSomething(time.Second, expected/2)
		logger.Info("Expected value published")
	}()

	// Assert
	if <-results == expected {
		t.Error("Wasn't expecting the expected value on the first read")
	}

	if <-results == expected {
		t.Error("Wasn't expecting the expected value on the second read")
	}

	// since the buffer was full until now this will be blocking until the result from doubleSomething is published!
	got := <-results
	if got != expected {
		t.Errorf("Expected %d on the third read, but found %d!", expected, got)
	}
}

// We can use channels to synchronize multiple goroutines!
func TestChannelsCanBeUsedToSynchronizeMultiplesRoutines(t *testing.T) {
	// Arrange
	got := make(chan int, 1)
	logger := initializeZap()
	expected := 8

	// Act
	logger.Info("Launching a new goroutine with a channel")
	go doubleSomethingWithAChannel(time.Second, expected/2, got, func() {
		logger.Info("Done executing the spun-off routine!")
	})
	logger.Info("Goroutine launched, moving forward...")

	// Assert
	result := <-got
	if result != expected {
		t.Errorf("Expected %d but found %d!", expected, result)
	}
}

// We can define the direction of a channel when using it as an argument by suffixing or prefixing it with the <- operator
func TestChannelsCanHaveADirectionWhenUsedAsArguments(t *testing.T) {
	// Arrange
	messageChannel := make(chan string, 1)
	logger := initializeZap()
	expected := "ping!"
	var doneSomething = false
	var got = ""

	// Act
	go pong(messageChannel, func(msg string) {
		logger.Info("pong has been executed!")
		doneSomething = true
		got = msg
	})
	go ping(expected, messageChannel)
	time.Sleep(time.Millisecond * 200)

	// Assert
	if !doneSomething {
		t.Error("Expected something to be done, but nothing was done at all")
	}

	if got != expected {
		t.Errorf("Expected %v but found %v", expected, got)
	}
}

// We can use select to iterate and parallelize work across multiple channels
func TestSelectCanBeUsedToParallelizeChannels(t *testing.T) {
	// Arrange
	logger := initializeZap()
	const firstExpected = "hello"
	const secondExpected = "world"
	firstChannel := make(chan string) // creating two blocking channels
	secondChannel := make(chan string)

	go func() {
		logger.Info("First channel is sleeping")
		time.Sleep(time.Millisecond * 500)
		firstChannel <- firstExpected
		logger.Info("First channel done!")
	}()

	go func() {
		logger.Info("Second channel is sleeping")
		time.Sleep(time.Millisecond * 505)
		secondChannel <- secondExpected
		logger.Info("Second channel done!")
	}()

	// Act
	select {
	case firstGot := <-firstChannel:
		// Assert
		if firstGot != firstExpected {
			t.Errorf("Expected %v but found %v", firstExpected, firstGot)
		}
	case _ = <-secondChannel:
		t.Error("This should never be executed before the firstGot!")
	}

	// Act
	select {
	case _ = <-firstChannel:
		// Since nothing is ever produced again this should remain blocked
		t.Error("This should never have executed again!")
	case secondGot := <-secondChannel:
		// Assert
		// And this should have the extra milliseconds it need to execute!
		if secondGot != secondExpected {
			t.Errorf("Expected %v but found %v", secondExpected, secondGot)
		}
	}
}
