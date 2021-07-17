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
	var expected = []int{2, 2, 4, 6, 10, 16, 26}
	for i := 0; i < len(subject); i++ {
		if Double(subject[i]) != expected[i] {
			t.Errorf("Expected %v but received %v", expected[i], Double(subject[i]))
		}
	}
}
