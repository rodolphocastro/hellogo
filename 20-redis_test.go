package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
	"time"
)

const (
	redisCredentials = "reredisdis"
	pathToRedisK8s   = "./environments/development/redis.yml"
	redisDb          = 0
)

// RedisSuite contains all the required tools and information for dealing with
// redis integration tests.
type RedisSuite struct {
	suite.Suite
	PathToK8sFile string
	Context       context.Context
	Logger        *zap.Logger
	RedisAddress  string
	RedisClient   *redis.Client
}

// setupRedisClient sets up a redis client based on an address and a password
func setupRedisClient(address, password string) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       redisDb,
	})
	return redisClient
}

// SetupSuite sets the suite up by initializing stuff and creating shared instances.
func (s *RedisSuite) SetupSuite() {
	s.Logger = InitializeLogger().
		With(
			zap.String("testSubject", "redis"),
		)
	s.RedisAddress = fmt.Sprintf("%v:6379", getMinikubeIp())
	s.Logger.Info("initializing the suite")
	s.Context = context.Background()
	s.PathToK8sFile = pathToRedisK8s
	s.Logger.Info("creating a RedisClient", zap.String("redisAddress", s.RedisAddress))
	s.RedisClient = setupRedisClient(s.RedisAddress, redisCredentials)
	s.Logger.Info("created a RedisClient", zap.String("redisAddress", s.RedisAddress))
	SpinUpK8s(s.T(), s.PathToK8sFile)
	time.Sleep(time.Second)
}

func (s *RedisSuite) TearDownSuite() {
	s.Logger.Info("tearing down the suite")
	err := s.RedisClient.Close()
	if err != nil {
		s.Logger.Error("unexpected error disconnecting from Redis", zap.Error(err))
	}
	s.Logger.Info("deleting the redis environment")
	CleanUpK8s(s.T(), pathToRedisK8s)
	s.Logger.Info("redis environment deleted")
}

func (s *RedisSuite) TestSetRedisValue() {
	// Arrange
	redisValue := "lorem ipsum dolor sit amet"
	redisKey := "myAwesomeKey"

	// Act
	s.Logger.Info("setting a value in Redis", zap.String("key", redisKey), zap.String(
		"value",
		redisValue,
	))
	err := s.RedisClient.Set(s.Context, redisKey, redisValue, 0).Err()

	// Assert
	s.Assert().Nil(err, "no errors were expected but got one")
}

func (s *RedisSuite) TestReadASetValueFromRedis() {
	// Arrange
	const expected = "lorem ipsum dolor sit amet"
	redisValue := expected
	redisKey := "myAwesomeKey2"
	err := s.RedisClient.Set(s.Context, redisKey, redisValue, 0).Err()
	if err != nil {
		s.Logger.Error("unexpected error setting a redis value", zap.Error(err))
	}

	// Act
	s.Logger.Info("reading a value from Redis", zap.String("key", redisKey))
	got, err := s.RedisClient.Get(s.Context, redisKey).Result()

	// Assert
	s.Nil(err, "no errors should happen when reading a known value")
	s.Equal(expected, got, "the returned value should match the set value")
}

func TestRedisSuite(t *testing.T) {
	SkipTestIfMinikubeIsUnavailable(t)
	// Delegate to testify's suite
	suite.Run(t, new(RedisSuite))
}
