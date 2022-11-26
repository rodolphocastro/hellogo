package main

import (
	"github.com/tkrajina/gpxgo/gpx"
	"os"
	"testing"
	"time"
)

const pathToGuineaPig = "./14-test-subject.gpx"

// It should be possible to read those gpx files into memory.
func TestReadFile(t *testing.T) {
	// Arrange
	fileBytes, err := os.ReadFile(pathToGuineaPig)
	if err != nil {
		t.Errorf("Expected to open the file successfully, but: %v", err)
	}

	// Act
	got, err := gpx.ParseBytes(fileBytes)
	secondGot := readGuineaPigFile(t)

	// Assert
	if err != nil {
		t.Errorf("Expected no errors parsing the file, but got: %v", err)
	}

	if got.Description != secondGot.Description {
		t.Error("Expected both got to be equal, but they aren't")
	}
}

// Read Guinea Pig File
//
// This func reads the 14-test-subject.gpx file and returns its data.
// If an error happens any test using this func fails automatically.
func readGuineaPigFile(t *testing.T) *gpx.GPX {
	gpxData, err := gpx.ParseFile(pathToGuineaPig)
	if err != nil {
		t.Errorf("Expected no errors, but got: %v", err)
	}
	return gpxData
}

func TestSumTotalDistance(t *testing.T) {
	// Arrange
	gpxData := readGuineaPigFile(t)
	totalDistance := 0.0

	// Act
	for _, track := range gpxData.Tracks {
		totalDistance += track.MovingData().MovingDistance
	}

	// Assert
	t.Log(totalDistance)
	if totalDistance < 21 {
		t.Errorf("Expected at least 21km but found %v", totalDistance)
	}
}

func TestGetDataForATimeStamp(t *testing.T) {
	// Arrange
	timestamp, _ := time.Parse(time.RFC3339, "2022-06-12T06:30:00-03:00")
	testSubject := readGuineaPigFile(t)

	// Act
	for _, track := range testSubject.Tracks {
		result := track.PositionAt(timestamp)
		t.Log(result)

		// Assert
		if result == nil {
			t.Error("Expected a position to exist, but nothing was found")
		}
	}
}
