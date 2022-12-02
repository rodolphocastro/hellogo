package main

import (
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

const (
	redisCredentials = "reredisdis"
	pathToRedisK8s   = "./environments/development/redis.yml"
)

// RedisSuite contains all the required tools and information for dealing with
// redis integration tests.
type RedisSuite struct {
	suite.Suite
	credentials   string
	pathToK8sFile string
	Logger        *zap.Logger
}

// SetupSuite sets the suite up by initializing stuff and creating shared instances.
func (s *RedisSuite) SetupSuite() {
	s.Logger = InitializeLogger().
		With(
			zap.String("testSubject", "redis"),
		)
	s.Logger.Info("initializing the suite")
	s.credentials = redisCredentials
	s.pathToK8sFile = pathToRedisK8s
}

func TestRedisSuite(t *testing.T) {
	// Delegate to testify's suite
	suite.Run(t, new(RedisSuite))
}
