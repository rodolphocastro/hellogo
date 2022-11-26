package main

import (
	"os"
	"testing"
)

func TestIfNoCIEnvIsSetReturnsFalse(t *testing.T) {
	// Arrange
	const expected = false
	err := os.Setenv(cicdPipelineEnvKey, "")
	if err != nil {
		t.Errorf("Error while changing current Env values: %v", err)
	}

	// Act
	got := isEnvironmentCI()

	// Assert
	if got {
		t.Errorf("expected %v but found %v - env contains %v instead", expected, got, os.Getenv(cicdPipelineEnvKey))
	}
}

func TestIfCIEnvIsSetReturnsTrue(t *testing.T) {
	// Arrange
	const expected = true
	err := os.Setenv(cicdPipelineEnvKey, "pudim")
	if err != nil {
		t.Errorf("Error while changing current Env values: %v", err)
	}

	// Act
	got := isEnvironmentCI()

	// Assert
	if !got {
		t.Errorf("expected %v but found %v - env contains %v instead", expected, got, os.Getenv(cicdPipelineEnvKey))
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
	expectEmpty := !isEnvironmentCI()

	// Act
	got := getMinikubeIp()

	// Assert
	if expectEmpty && got != "" {
		t.Errorf("Expected empty but got %v", got)
	}
}

func TestInitializeLogger(t *testing.T) {
	// Arrange

	// Act
	got := InitializeLogger()
	got.Info("it lives!!")

	// Assert
	if got == nil {
		t.Fatal("expected a Zap logger to be created but none was")
	}
}

func TestGetMinikubeStatus(t *testing.T) {
	// Arrange
	expected := getMinikubeIp() != ""

	// Act
	got, stringResult := GetMinikubeStatus()

	// Assert
	if got != expected {
		t.Errorf("expected minikube running status to be %v but found %v instead. Report from the command reads %v",
			expected, got, stringResult)
	}
}
