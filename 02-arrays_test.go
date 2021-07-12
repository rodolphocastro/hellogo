package main

import (
	"testing"
)

func Double(lhs int) int {
	return lhs * 2
}

func TestForEach(t *testing.T) {
	var subject = []int{1, 2, 3, 4, 5}
	var expected = []int{2, 4, 6, 8, 10}
	for idx, element := range subject {
		result := Double(element)
		t.Logf("Doubling %v returned %v", element, result)
		if expected[idx] != result {
			t.Error("Double returned an unexpected value")
		}
	}
}

func TestForLoop(t *testing.T) {
	var subject = []int{1, 1, 2, 3, 5, 8, 13}
	for i := 0; i < len(subject); i++ {
		if i >= len(subject) {
			t.Errorf("%v is greater than the maximum allowed length", i)
		}
	}
}
