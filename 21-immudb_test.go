package main

import (
	"context"
	immudb "github.com/codenotary/immudb/pkg/client"
	"testing"
)

const (
	pathToImmudb = "./environments/development/immudb.yml"
)

func TestConnectToDatabase(t *testing.T) {
	// Arrange
	SkipTestIfMinikubeIsUnavailable(t)
	SpinUpK8s(t, pathToImmudb)
	//time.Sleep(time.Second*15)

	// Act
	opts := immudb.
		DefaultOptions().
		WithPort(6666).
		WithAddress(getMinikubeIp())

	client := immudb.NewClient().WithOptions(opts)
	err := client.OpenSession(context.TODO(), []byte("immudb"), []byte("immudb"), "defaultdb")

	// Assert
	if err != nil {
		t.Errorf("Expected no errors but found %v", err)
	}

	client.CloseSession(context.TODO())
	CleanUpK8s(t, pathToImmudb)
}
