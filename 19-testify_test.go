package main

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestTestifyEquals(t *testing.T) {
	// Arrange
	scenarios := map[int]int{
		1:    1,
		1000: 1000,
		-2:   -2,
	}
	// Act and Assert
	for expected, got := range scenarios {
		t.Run(strconv.Itoa(expected), func(t2 *testing.T) {
			assert.Equal(t2, expected, got, "map key and value should be equal")
		})
	}
}
