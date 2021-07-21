package main

import (
	"testing"
)

// Declaring a function that takes an int and returns another int
func DoubleAnInt(lhs int) int {
	return lhs * 2
}

// Declaring a function that takes two ints and returns a third int
func AddTwoInts(lhs int, rhs int) int {
	return lhs + rhs
}

// Declaring a function that takes three ints but typing they at once
func AddThreeInts(lhs, mid, rhs int) int {
	return AddTwoInts(AddTwoInts(lhs, rhs), mid)
}

// A function may return more than one value
func IsEven(lhs int) (int, bool) {
	if lhs%2 == 0 {
		return lhs, true
	}

	return lhs, false
}

// A function may have variadic parameters (ie: N parameters of a type)
func SumInts(nums ...int) int {
	result := 0
	for _, i := range nums {
		result += i
	}
	return result
}

// A function may have a named return value
func DoPi() (pi float32) {
	pi = 3.14
	return
}

func TestFuncs(t *testing.T) {
	const lhs = 5
	const rhs = 8
	const mid = 10
	const anOddNumber = 7
	const anEvenNumber = 10
	const expectedDoubleLhs = lhs * 2
	const expectedSum = lhs + rhs
	const expectedTripleSum = expectedSum + mid
	if DoubleAnInt(lhs) != expectedDoubleLhs {
		t.Errorf("Expected %v as a double for %v but that didn't happen", expectedDoubleLhs, lhs)
	}

	if AddTwoInts(lhs, rhs) != expectedSum {
		t.Errorf("Expected %v as the sum of %v and %v, but that didn't happen", expectedSum, lhs, rhs)
	}

	if AddThreeInts(lhs, mid, rhs) != expectedTripleSum {
		t.Errorf("Expected %v as the sum of %v, %v and %v, but that didn't happen", expectedTripleSum, lhs, mid, rhs)
	}

	if AddTwoInts(lhs, rhs) != AddTwoInts(rhs, lhs) {
		t.Error("Additions shouldn't care about the order of its elements")
	}

	if AddThreeInts(lhs, mid, rhs) != AddThreeInts(mid, rhs, lhs) {
		t.Error("Additions shouldn't care about the order of its elements")
	}

	_, isOddEven := IsEven(anOddNumber)
	if isOddEven {
		t.Errorf("Expected %v to be odd but, somehow, it's even", anOddNumber)
	}

	_, isEvenEven := IsEven(anEvenNumber)
	if !isEvenEven {
		t.Errorf("Expected %v to be even but, somehow, it's odd", anEvenNumber)
	}

	if SumInts(lhs, rhs, mid) != AddThreeInts(mid, rhs, lhs) {
		t.Error("Additions shouldn't care about the order of its elements")
	}

	if SumInts(lhs, mid) != AddTwoInts(mid, lhs) {
		t.Error("Additions shouldn't care about the order of its elements")
	}

	// A slice/array may be passed in as a variadic parameter
	sliceParam := []int{lhs, rhs, mid}
	if SumInts(lhs, rhs, mid) != SumInts(sliceParam...) {
		t.Error("Additions shouldn't care about the order of its elements")
	}

	if DoPi() != 3.14 {
		t.Error("Expected pi (with 2 digits), found something else")
	}
}
