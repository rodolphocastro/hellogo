package main

import (
	"testing"
)

func TestSlices(t *testing.T) {
	// Creating a slice
	mySlice := make([]int, 0)
	if len(mySlice) != 0 {
		t.Errorf("Expected slice to be empty but it has %v elements", len(mySlice))
	}

	// Adding stuff to a slice
	mySlice = append(mySlice, -42, 42, 0, 5)
	if len(mySlice) != 4 {
		t.Errorf("Expected slice to have 4 elements but it has %v elements", len(mySlice))
	}

	// Removing stuff from a slice
	mySlice = mySlice[:len(mySlice)-1]
	if len(mySlice) != 3 {
		t.Errorf("Expected slice to have 3 elements but it has %v elements", len(mySlice))
	}

	// Iterating thru a slice
	for idx, v := range mySlice {
		if v == -999 {
			t.Errorf("Expected to see anything other than -999 but it is the %vth element on the slice", idx)
		}
	}
}
