package main

import (
	"os"
	"testing"
)

func TestIfNoCIEnvIsSetReturnsFalse(t *testing.T) {
	// Arrange
	err := os.Setenv("CI", "")
	if err != nil {
		t.Errorf("Error while changing current Env values: %v", err)
	}

	// Act
	got := isEnvironmentCI()

	// Assert
	if got {
		t.Error("Expected environment not to be CI, but found a CI environment")
	}
}

func TestIfCIEnvIsSetReturnsTrue(t *testing.T) {
	// Arrange
	err := os.Setenv("CI", "pudim")
	if err != nil {
		t.Errorf("Error while changing current Env values: %v", err)
	}

	// Act
	got := isEnvironmentCI()

	// Assert
	if !got {
		t.Error("Expected environment be CI, but found a non-CI environment")
	}
}

func FuzzApplyPathOrDefault(f *testing.F) {
	f.Add("/usr/aFile.yml", false)
	f.Add("", true)
	f.Add(" ", true)
	f.Add("c:\\Users\\AnUser\\k8s.yml", false)
	f.Fuzz(func(t *testing.T, path string, shouldGoToDefault bool) {
		got := getPathOrDefault(path)
		t.Log(got)
		if got != path && !shouldGoToDefault {
			t.Errorf("Expected %v but got %v", path, got)
		}
	})
}

func TestGetMinikubeIp(t *testing.T) {
	// Arrange
	expectEmpty := isEnvironmentCI()

	// Act
	got := getMinikubeIp()
	t.Log(got)

	// Assert
	if expectEmpty && got != "" {
		t.Errorf("Expected empty but got %v", got)
	}
}
