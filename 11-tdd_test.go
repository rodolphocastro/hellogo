package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// A dependecy we'll rely on
type thingDoer interface {
	say(name string)
}

// A high-level function that depends on said dependency
func speakWithDoer(name string, doer thingDoer) {
	doer.say(name)
}

// A fake doer, manually built by us
type fakeDoer struct {
	numCalls int
}

// Method attached to our manual fake doer
func (doer *fakeDoer) say(name string) {
	doer.numCalls = doer.numCalls + 1
}

func TestTdd(t *testing.T) {
	// Doing tdd the classic way, not using any libraries and doing a mock manually
	t.Run("manual fake thing doer does something", func(t *testing.T) {
		subject := fakeDoer{
			numCalls: 0,
		}
		speakWithDoer("Heisenberg", &subject)
		if subject.numCalls == 0 {
			t.Errorf("Expected at least 1 call but got %v", subject.numCalls)
		}
	})

	// Doing tdd with an assertion library but a manual mock
	t.Run("manual fake thing doer does something but is asserted with testify", func(t *testing.T) {
		subject := fakeDoer{
			numCalls: 0,
		}
		speakWithDoer("Heisenberg", &subject)
		assert.Greater(t, subject.numCalls, 0, "subject should've been called at least once")
	})
}
