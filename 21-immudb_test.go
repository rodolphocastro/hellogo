package main

import (
	"context"
	immudb "github.com/codenotary/immudb/pkg/client"
	"testing"
	"time"
)

const (
	pathToImmudb = "./environments/development/immudb.yml"
)

// getImmudbAddress gets the address for the Immudb Instance
func getImmudbAddress() string {
	return getMinikubeIp()
}

func TestConnectToDatabase(t *testing.T) {
	// Arrange
	SkipTestIfMinikubeIsUnavailable(t)
	SpinUpK8s(t, pathToImmudb)
	time.Sleep(time.Second)

	// Act
	opts := immudb.DefaultOptions().WithAddress(getImmudbAddress()).WithPort(32003)
	client := immudb.NewClient().WithOptions(opts)
	err := client.OpenSession(context.TODO(), []byte("immudb"), []byte("immudb"), "defaultDb")

	// Assert
	if err != nil {
		t.Errorf("Expected no errors but found %v", err)
	}

	//CleanUpK8s(t, pathToImmudb)
}
