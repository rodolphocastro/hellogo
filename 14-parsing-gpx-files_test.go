package main

import (
	"github.com/tkrajina/gpxgo/gpx"
	"os"
	"testing"
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

// TODO: https://github.com/tkrajina/gpxgo/blob/321f19554eecf2c5ba914f2bfad70b4458e2819f/gpx/gpx.go#L117 - Iterate until you figure out the total distance.
