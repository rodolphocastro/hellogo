package main

import (
	"context"
	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
	"time"
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
	i.Logger.Debug("initializing the suite and giving it some time to start")
	SpinUpK8s(i.T(), i.PathTok8sFile)
	time.Sleep(time.Second * 2)
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

// TestSetAndGetUnverifiedValues demonstrates how to set and get (write and read)
// values from immudb.
func (i *ImmudbSuite) TestSetAndGetUnverifiedValues() {
	// Arrange
	key := faker.Word()
	expected := faker.Sentence()
	client := i.Client
	logger := i.Logger.With(zap.String("immudbKey", key), zap.String("immudbValue", expected))

	// Act
	// setting a value
	logger.Debug("setting a value to immudb")
	txHeader, err := client.Set(i.Context, []byte(key), []byte(expected))
	i.Require().Nil(err, "no errors should happen when writing to immudb")
	i.NotNil(txHeader, "a transaction should be started")

	// reading the value we just set
	logger.Debug("reading a value from immudb")
	gotSchema, err := client.Get(i.Context, []byte(key))
	i.Require().Nil(err, "no errors should happen when reading from immudb")
	got := string(gotSchema.Value)
	i.NotEmpty(got, "the result should not be empty")
	i.Equal(expected, got, "the result should match the set value")
}

// TODO: https://docs.immudb.io/master/develop/reading.html#get-and-set for Verifications.

func TestImmudbSuite(t *testing.T) {
	// Arrange
	SkipTestIfMinikubeIsUnavailable(t)
	suite.Run(t, new(ImmudbSuite))
}
