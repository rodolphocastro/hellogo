package main

import (
	"github.com/stretchr/testify/assert"
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
