package main

import (
	"go.uber.org/zap"
	"reflect"
	"testing"
	"time"
)

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
	logger := InitializeLogger()
	firstExpected := 8
	secondExpected := 66
	firstInput := firstExpected / 2
	secondInput := secondExpected / 2

	// Act
	// the 'go' keyword means 'fire this and forget about it', which is pretty useful to allow the runtime to deal
	// with spawning async work
	go func() {
		logger.Debug("First routine launched")
		firstGot = doubleSomething(0, firstInput)
		logger.Debug("First routine completed")
	}()
	// giving it some time, just in case
	time.Sleep(time.Millisecond)

	// this one should take at least a second to run
	go func() {
		logger.Debug("Second routine launched")
		secondGot = doubleSomething(time.Second, secondInput)
		logger.Debug("Second routine launched")
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
	logger := InitializeLogger()
	expected := 8
	// the 'chan' keyword indicates we want to create a Channel to allow communication to happen
	results := make(chan int)

	// Act
	// firing a goroutine that outputs to the results
	go func() {
		logger.Debug("Starting a goroutine")
		// the right hand side (rhs) <- operator publishes something to a channel
		logger.Debug("Expected value has been published")
		results <- doubleSomething(time.Second, expected/2)
		logger.Debug("An unexpected value has been published now")
		results <- 0
	}()

	// the left hand side (lhs) <- operator means to read something from the channel
	// note: both lhs and rhs are *blocking*. Which means that the routine will halt execution until something is read
	// and the receiving routine will also halt until something is received
	logger.Debug("Waiting for got to be published!")
	got := <-results
	logger.Debug("Got has been updated", zap.Int("got", got))

	// Assert
	// since <- blocks the reader and the publisher this should just work
	if got != expected {
		t.Errorf("Expected %d but found %d!", expected, got)
	}

	// now we're asking for a second value!
	got = <-results
	logger.Debug("Got has been updated", zap.Int("got", got))
	if got == expected {
		t.Errorf("Expected anything other than %d, but found %d", expected, expected)
	}
}

// Channels with a buffer allow goroutines to output more than a single value before blocking!
func TestBufferedChannelsBlockLessOften(t *testing.T) {
	// Arrange
	logger := InitializeLogger()
	expected := 8
	// the 'chan' keyword indicates we want to create a Channel to allow communication to happen
	// but this time we're also making a buffer of 2 which will allow a goroutine to output 3 times before it gets blocked!
	results := make(chan int, 2)

	// Act
	go func() {
		logger.Debug("Firing off a goroutine!")
		results <- 0
		logger.Debug("Published 0")
		results <- -1
		logger.Debug("Published -1")
		logger.Debug("About to publish the expected value!")
		results <- doubleSomething(time.Second, expected/2)
		logger.Debug("Expected value published")
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
	logger := InitializeLogger()
	expected := 8

	// Act
	logger.Debug("Launching a new goroutine with a channel")
	go doubleSomethingWithAChannel(time.Second, expected/2, got, func() {
		logger.Debug("Done executing the spun-off routine!")
	})
	logger.Debug("Goroutine launched, moving forward...")

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
	logger := InitializeLogger()
	expected := "ping!"
	var doneSomething = false
	var got = ""

	// Act
	go pong(messageChannel, func(msg string) {
		logger.Debug("pong has been executed!")
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
	logger := InitializeLogger()
	const firstExpected = "hello"
	const secondExpected = "world"
	firstChannel := make(chan string) // creating two blocking channels
	secondChannel := make(chan string)

	go func() {
		logger.Debug("First channel is sleeping")
		time.Sleep(time.Millisecond * 500)
		firstChannel <- firstExpected
		logger.Debug("First channel done!")
	}()

	go func() {
		logger.Debug("Second channel is sleeping")
		time.Sleep(time.Millisecond * 505)
		secondChannel <- secondExpected
		logger.Debug("Second channel done!")
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

// We can use timeouts to elegantly do something else if a channels takes to long to publish a message!
func TestTimeoutsMayBeUsedWhenReadingFromChannels(t *testing.T) {
	// Arrange
	logger := InitializeLogger()
	gotChannel := make(chan int, 1)
	const expected = 1001

	// Act
	go func(output chan<- int) {
		logger.Debug("Firing off a goroutine")
		time.Sleep(time.Second)
		output <- 2000
		logger.Debug("Done publishing an unexpected value!")
	}(gotChannel)

	// Assert
	select {
	case unexpected := <-gotChannel:
		t.Errorf("Expected not getting anything at all due to the timeout, but got %v", unexpected)
	case _ = <-time.After(time.Millisecond * 250): // time.After creates and publishes to a channel after a specified amount of time! Thus being a "timeout" of sorts!
		logger.Debug("Nothing happened, going to publish the expected result and try again")
		gotChannel <- expected
	}

	select {
	case got := <-gotChannel:
		// Since we 'short-circuited' in the previous select the expected value should have been published sooner than the output
		logger.Debug("Something was available in the channel, getting it!")
		if got != expected {
			t.Errorf("Expected %v but got %v!", expected, got)
		}
	case _ = <-time.After(time.Millisecond):
		t.Error("Expected the channel to have an immediate result but we ended up waiting!")
	}
}

// We can also use a default cause within a select to allow non-blocking operations to execute as part of a channel!
func TestSelectCanBeUsedToCreateNonBlockingOperations(t *testing.T) {
	// Arrange
	const expected = "Olá, mundo!"
	logger := InitializeLogger()
	commsChannel := make(chan string) // note: This channel doesn't have a buffer!

	// Act
	select {
	case _ = <-commsChannel:
		// Assert
		t.Error("didn't expected a message to be available yet")
	default:
		logger.Debug("publishing expected to a channel, if non-blocking")
		select {
		case commsChannel <- expected:
			t.Error("expected to not publish but we were able to publish")
		default:
			logger.Debug("nothing has been published")
		}
	}

	// Act
	select {
	case _ = <-commsChannel:
		// Assert
		t.Error("expected this block to not execute since it is blocking, but it executed")
	default:
		logger.Debug("nothing has been received")
	}
}

// We can manually close() channels to signal that no more work should be done
func TestChannelsCanBeClosedToSignalNoMoreValuesWillBeSent(t *testing.T) {
	// Arrange
	logger := InitializeLogger()
	workQueue := make(chan int, 5)
	done := make(chan bool)

	logger.Debug("initializing a worker goCoroutine")
	go func() {
		for {
			current, isOpen := <-workQueue
			if isOpen {
				logger.Debug("received a new job",
					zap.Int("currentWork", current),
				)
			} else {
				logger.Debug("all jobs have been received, shutting down")
				done <- true
				close(done)
				logger.Debug("closed the return channel")
				return
			}
		}
	}()

	// Act
	for i := 0; i < 7; i++ {
		workQueue <- i
	}
	logger.Debug("closing the workQueue")
	close(workQueue)
	got := <-done

	// Assert
	if !got {
		t.Errorf("expected done to be true but got %v", got)
	}
}

// We can use the range keyword to also iterate over results from a channel!
func TestRangesCanBeUsedToIterateOverAChannelResults(t *testing.T) {
	// Arrange
	const sliceLimit = 5
	logger := InitializeLogger()
	intChannel := make(chan int) // using a non-buffered channel to cause blocks
	gotInts := make([]int, 0)
	expectedInts := make([]int, 0)
	for i := 0; i < sliceLimit; i++ {
		expectedInts = append(expectedInts, i*3)
	}

	act := func(output chan<- int) {
		logger.Debug("beginning to publish")
		for idx, expectedInt := range expectedInts {
			logger.Debug("publishing a new int to the channel",
				zap.Int("currentInt", expectedInt),
				zap.Int("currentPosition", idx),
			)
			output <- expectedInt
		}
		logger.Debug("done publishing")
		close(output)
		return
	}

	// Act
	// firing off a goRoutine
	go act(intChannel)

	logger.Debug("beginning to iterate over the channel")
	for got := range intChannel {
		logger.Debug("received a new int from the channel", zap.Int("currentInt", got))
		gotInts = append(gotInts, got)
	}
	logger.Debug("done receiving ints")

	// Assert
	if !reflect.DeepEqual(gotInts, expectedInts) { // DeepEqual allows us to compare everything in two slices at once
		t.Errorf("expected %v but got %v instead", expectedInts, gotInts)
	}
}
