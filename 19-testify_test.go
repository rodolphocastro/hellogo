package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"strconv"
	"testing"
)

// a single shared testifyLogger instance for testify tests
var testifyLogger = InitializeLogger().
	With(
		zap.String("testSubject", "testify"),
	)

// TestTestifyAssertions uses the assert function provided by testify to sugar-coat goLang's testing tools.
func TestTestifyAssertions(t *testing.T) {
	// Arrange
	scenarios := map[int]int{
		1:    1,
		1000: 1000,
		-2:   -2,
		42:   42,
	}
	assertionsLogger := testifyLogger.With(zap.Int("totalTestCases", len(scenarios)))

	// Act and Assert
	for expected, got := range scenarios {
		assertionsLogger.Info("starting assertions for a new scenario", zap.Int("currentTestKey", expected), zap.Int("currentTestValue", got))
		t.Run(strconv.Itoa(expected), func(t2 *testing.T) {
			//goland:noinspection GoImportUsedAsName
			assert := assert.New(t2) // creating a local assert to allow us to save a few keys passing in t2 everytime
			assert.Equal(expected, got, "map key and value should be equal")
			assert.NotEmpty(got, "map value should not be empty, ever")
			assert.NotEqual(expected, got+1, "map and value+1 should not be equal ever")
			assert.NotZero(got, "the value shouldn't be the default int")
			assert.NotNil(got, "the value shouldn't be nil neither")
			assert.NotSame(expected, got, "the key and its value should be different pointers")
		})
		assertionsLogger.Info("done asserting", zap.Int("currentTestKey", expected), zap.Int("currentTestValue", got))
	}
}

// FriendlyStructFetcherMocked is a mocked version of FriendlyStructFetcherImpl
type FriendlyStructFetcherMocked struct {
	mock.Mock
}

// fetch is a stub implementation to allow FriendlyStructFetcherMocked to match the FriendlyStructFetcher interface.
func (receiver *FriendlyStructFetcherMocked) fetch() int {
	return receiver.Called().Int(0)
}

// FriendlyStructFetcherImpl is a very real struct that does something. trust me, I'm an engineer.
type FriendlyStructFetcherImpl struct {
	returnValue int
}

// FriendlyStructFetcher defines methods that need to exist in order for a
// something to be fetched. Nasty stuff.
type FriendlyStructFetcher interface {
	// fetch returns an int for something.
	fetch() int
}

// fetch returns the return value for the real deal struct!
func (receiver FriendlyStructFetcherImpl) fetch() int {
	return receiver.returnValue
}

// TestTestifyMocks uses the mock object from testify to expedite creating mocks
// for interfaces and types. this test relies on implementation, mocks and stubs
// for the FriendlyStructFetcher interface. Those implementations and mocks are
// the FriendlyStructFetcherImpl and FriendlyStructFetcherMocked structs.
func TestTestifyMocks(t *testing.T) {
	// Arrange
	scenarios := map[int]int{
		42:   24,
		55:   55,
		-10:  10,
		5555: 1231,
	}
	assertionsLogger := testifyLogger.With(zap.Int("totalTestCases", len(scenarios)))

	for input, expected := range scenarios {
		assertionsLogger.Info("beginning a new scenario", zap.Int("scenarioInput", input), zap.Int("scenarioOutput", expected))
		t.Run(strconv.Itoa(input), func(t2 *testing.T) {
			// Arrange
			assertions := assert.New(t2)
			var nonMocked FriendlyStructFetcher = FriendlyStructFetcherImpl{returnValue: expected}
			var mocked = new(FriendlyStructFetcherMocked)
			mocked.
				On("fetch").                    // telling that whenever the 'fetch' method is called
				Run(func(args mock.Arguments) { // we should run a specific closure
					assertionsLogger.Info("mocking and logging - oh yeah")
				}).
				Return(expected) // and finally return something

			// Act
			gotReal := nonMocked.fetch()
			gotMocked := mocked.fetch()

			// Assert
			assertions.Equal(expected, gotReal, "the real implementation should behave as expected")
			assertions.Equal(expected, gotMocked, "the mocked implementation should also match its settings")
			assertions.Equal(gotMocked, gotReal, "the real and mocked results should also match")
		})
		assertionsLogger.Info("done with the scenario", zap.Int("scenarioInput", input), zap.Int("scenarioOutput", expected))
	}
}

// TODO: https://github.com/stretchr/testify#require-package
// TODO: https://github.com/stretchr/testify#suite-package
