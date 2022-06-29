package main

import (
	"go.uber.org/zap"
	"testing"
)

func TestCreateZapDevelopment(t *testing.T) {
	// Arrange

	// Act
	_, err := zap.NewDevelopment()

	// Assert
	if err != nil {
		t.Errorf("Expected no error but found: %v", err)
	}
}

func TestCreateSugarZapLogger(t *testing.T) {
	// Arrange
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Act
	got := logger.Sugar()
	got.Infow("A wild sugar log appeared", "a key", "a value", "another key", 42)

	// Assert
	if got == nil {
		t.Error("Expected a logger but found nil")
	}
}

func TestCreateZapLogger(t *testing.T) {
	// Arrange
	got, _ := zap.NewDevelopment()
	defer got.Sync()

	// Act
	got.Info("A wild sugar log appeared 2", zap.String("a key", "a value"), zap.Int("another key", 42))

	// Assert
	if got == nil {
		t.Error("Expected a logger but found nil")
	}
}
