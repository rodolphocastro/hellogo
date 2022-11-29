package main

import (
	"go.uber.org/zap"
	"os"
	"testing"
)

func TestIfNoCIEnvIsSetReturnsFalse(t *testing.T) {
	// Arrange
	const expected = false
	t.Setenv(cicdPipelineEnvKey, "")

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
	t.Setenv(cicdPipelineEnvKey, "pudim")

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
	scenarios := []string{
		"",
		"CI",
		"randomEnvValue",
	}

	for _, envValue := range scenarios {
		// Arrange
		t.Setenv(cicdPipelineEnvKey, envValue)

		// Act
		got := InitializeLogger().With(zap.String(cicdPipelineEnvKey, envValue))
		got.Info("it lives!!")

		// Assert
		if got == nil {
			t.Fatal("expected a Zap logger to be created but none was")
		}
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
