package main

import (
	"testing"
)

// Closure is kinda like a C# Action/Func.
// AKA: A lamba, a functional bit of programming within OO
func LazySum(lhs int, rhs int) func() int {
	return func() int {
		return lhs + rhs
	}
}

func TestClosures(t *testing.T) {
	const lhs = 3
	const rhs = 13
	const expectedSum = 16

	sum := LazySum(lhs, rhs)
	if sum() != expectedSum {
		t.Errorf("Expected %v but received %v as result of the LazySum", expectedSum, sum())
	}

	secondSum := LazySum(sum(), lhs)
	thirdSum := LazySum(sum(), rhs)
	if secondSum() == thirdSum() {
		t.Error("Expected the lazy sums to be different, but they were equal")
	}
}
