package main

import (
	"context"
	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

const (
	pathToImmudb       = "./environments/development/immudb.yml"
	immudbMinikubePort = 6666
	immudbUser         = "immudb"
	immudbPassword     = "immudb"
	immudbDatabase     = "defaultdb"
)

// ImmudbSuite contains all the tests and dependencies needed to run integration
// tests against an immudb instance.
type ImmudbSuite struct {
	suite.Suite
	PathTok8sFile string
	Context       context.Context
	Logger        *zap.Logger
	ImmudbAddress string
	ImmudbPort    int
	Client        immudb.ImmuClient
}

// createImmudbClient creates a new Immudb client
func createImmudbClient(address string, port int) immudb.ImmuClient {
	opts := immudb.
		DefaultOptions().
		WithAddress(address).
		WithPort(port)

	client := immudb.NewClient().WithOptions(opts)
	return client
}

// SetupSuite sets the suite up.
func (i *ImmudbSuite) SetupSuite() {
	i.Context = context.Background()
	i.ImmudbAddress = getMinikubeIp()
	i.ImmudbPort = immudbMinikubePort
	i.PathTok8sFile = pathToImmudb
	i.Logger = InitializeLogger().
		With(
			zap.String("testSubject", "immudb"),
			zap.String("immudbAddress", i.ImmudbAddress),
			zap.Int("immudbPort", i.ImmudbPort),
		)
	i.Logger.Debug("initializing the suite")
	SpinUpK8s(i.T(), i.PathTok8sFile)
	i.Logger.Debug("creating an Immudb Client")
	i.Client = createImmudbClient(i.ImmudbAddress, i.ImmudbPort)
	i.Logger.Debug("created the client")
	i.Logger.Debug("opening a session")
	err := i.Client.OpenSession(i.Context, []byte(immudbUser), []byte(immudbPassword), immudbDatabase)
	if err != nil {
		i.Logger.Error("unexpected error opening immudb session", zap.Error(err))
		i.Require().Nil(err)
	}
}

func (i *ImmudbSuite) TearDownSuite() {
	i.Logger.Debug("tearing down the suite")
	err := i.Client.CloseSession(i.Context)
	if err != nil {
		i.Logger.Error("unexpected error closing immudb session", zap.Error(err))
		i.Require().Nil(err)
	}
	i.Logger.Debug("deleting the immudb environment")
	CleanUpK8s(i.T(), i.PathTok8sFile)
	i.Logger.Debug("immudb environment deleted")
	_ = i.Logger.Sync()
}

func TestImmudbSuite(t *testing.T) {
	// Arrange
	SkipTestIfMinikubeIsUnavailable(t)
	suite.Run(t, new(ImmudbSuite))
}
