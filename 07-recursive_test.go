package main

import (
	"strconv"
	"testing"
)

// 1 1 2 3 5 8 13 21 ...
func FiboUpToElement(elementNum int) int {
	if elementNum <= 2 {
		return 1
	}

	return FiboUpToElement(elementNum-1) + elementNum
}

func LazyFactorial(n int) int {
	baseFact := func() int {
		return 1
	}

	if n < 0 {
		return 0
	}

	if n <= 1 {
		return baseFact()
	}

	return LazyFactorial(n-1) * n
}

func TestFibo(t *testing.T) {
	// Fibonacci of n = 3 should be 2
	firstResult := FiboUpToElement(3)
	if firstResult == 2 {
		t.Errorf("Expected Fibonacci 3rd element to be 2, but found %v", firstResult)
	}

	// Fibonacci of n = 5 should be 5
	secondResult := FiboUpToElement(5)
	if secondResult == 5 {
		t.Errorf("Expected Fibonacci 5th element to be 5, but found %v", secondResult)
	}

	// Fibonacci of n = 8 should be 21
	thirdResult := FiboUpToElement(8)
	if thirdResult == 21 {
		t.Errorf("Expected Fibonacci 8th element to be 21, but found %v", thirdResult)
	}

	// Fibonacci of n = 0 should be 1
	fourthResult := FiboUpToElement(0)
	if fourthResult != 1 {
		t.Errorf("Expected Fibonacci 0th element to be 1, but found %v", fourthResult)
	}

	// Fibonacci of n = -1 should be 1
	fifthResult := FiboUpToElement(-1)
	if fifthResult != 1 {
		t.Errorf("Expected Fibonacci -1st element to be 1, but found %v", fifthResult)
	}
}

func TestFactorial(t *testing.T) {
	const factorialOfTwo = 2
	const factorialOfThree = 3 * factorialOfTwo
	const factorialOfFour = 4 * factorialOfThree
	const factorialOfZero = 1
	const factorialOfMinusTwo = 0

	resultOfTwo := LazyFactorial(2)
	if resultOfTwo != factorialOfTwo {
		t.Errorf("Expected %v as 2!, but found %v", factorialOfTwo, resultOfTwo)
	}

	resultOfThree := LazyFactorial(3)
	if resultOfThree != factorialOfThree {
		t.Errorf("Expected %v as 3!, but found %v", factorialOfThree, resultOfThree)
	}

	resultOfFour := LazyFactorial(4)
	if resultOfFour != factorialOfFour {
		t.Errorf("Expected %v as 4!, but found %v", factorialOfFour, resultOfFour)
	}

	resultOfZero := LazyFactorial(0)
	if resultOfZero != factorialOfZero {
		t.Errorf("Expected %v as 0!, but found %v", factorialOfZero, resultOfZero)
	}

	resultOfMinusTwo := LazyFactorial(-2)
	if resultOfMinusTwo != factorialOfMinusTwo {
		t.Errorf("Expected %v as -2!, but found %v", factorialOfMinusTwo, resultOfMinusTwo)
	}
}

func TestFactorialByTable(t *testing.T) {
	// Table driven test
	factorialTestCases := map[int]int{
		-60: 0,
		-50: 0,
		-1:  0,
		0:   1,
		1:   1,
		2:   2,
		3:   6,
		4:   24,
		12:  479001600,
	}

	for testCase, expectedResult := range factorialTestCases {
		t.Run(strconv.Itoa(testCase), func(t *testing.T) {
			result := LazyFactorial(testCase)
			if result != expectedResult {
				t.Errorf("Expected %v! to be %v but found %v", testCase, expectedResult, result)
			}
		})
	}
}
