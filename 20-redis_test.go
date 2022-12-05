package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-faker/faker/v4"
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
	s.RedisAddress = fmt.Sprintf("%v:6379", getMinikubeIp())
	s.Logger = InitializeLogger().
		With(
			zap.String("testSubject", "redis"),
			zap.String("redisAddress", s.RedisAddress),
		)
	s.Logger.Debug("initializing the suite")
	s.Context = context.Background()
	s.PathToK8sFile = pathToRedisK8s
	s.Logger.Debug("creating a RedisClient")
	s.RedisClient = setupRedisClient(s.RedisAddress, redisCredentials)
	s.Logger.Debug("created a RedisClient")
	SpinUpK8s(s.T(), s.PathToK8sFile)
	time.Sleep(time.Second * 8)
}

// TearDownSuite tears down the suite after all tests are executed.
func (s *RedisSuite) TearDownSuite() {
	s.Logger.Debug("tearing down the suite")
	err := s.RedisClient.Close()
	if err != nil {
		s.Logger.Error("unexpected error disconnecting from Redis", zap.Error(err))
	}
	s.Logger.Debug("deleting the redis environment")
	CleanUpK8s(s.T(), pathToRedisK8s)
	s.Logger.Debug("redis environment deleted")
	_ = s.Logger.Sync()
}

// TestSetRedisValue demonstrates how to set a value on a Redis instance.
func (s *RedisSuite) TestSetRedisValue() {
	// Arrange
	redisValue := faker.Sentence()
	redisKey := faker.Word()

	// Act
	s.Logger.Debug("setting a value in Redis", zap.String("key", redisKey), zap.String(
		"value",
		redisValue,
	))
	err := s.RedisClient.Set(s.Context, redisKey, redisValue, time.Millisecond*20).Err()

	// Assert
	s.Assert().Nil(err, "no errors were expected but got one")
}

// TestReadASetValueFromRedis demonstrates how to read a value that was been
// previously stored in a Redis instance.
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
	s.Logger.Debug("reading a value from Redis", zap.String("key", redisKey))
	got, err := s.RedisClient.Get(s.Context, redisKey).Result()

	// Assert
	s.Nil(err, "no errors should happen when reading a known value")
	s.Equal(expected, got, "the returned value should match the set value")
}

// TestReadAValueThatWasntSetReturnsNil demonstrates what happens when you read a
// value that isn't available in a Redis instance.
func (s *RedisSuite) TestReadAValueThatWasntSetReturnsNil() {
	// Arrange
	redisKey := "myAwesomeKey3"

	// Act
	s.Logger.Debug("reading a value from Redis", zap.String("key", redisKey))
	result, err := s.RedisClient.Get(s.Context, redisKey).Result()

	// Assert
	s.NotNil(err, "a redis error should be raised")
	s.Empty(result, "the resulting string should be empty")
}

// TestReadAValueThatWasntSetReturnsNil demonstrates how one can store a struct
// in a Redis instance and they recover it afterwards.
func (s *RedisSuite) TestMarshalSetGetUnmarshalShouldKeepData() {
	// Arrange
	const redisKey = "awesomeKey0099"
	expected := Pet{}
	err := faker.FakeData(&expected)
	if err != nil {
		s.Logger.Error("unexpected error faking data for a pet", zap.Error(err))
	}
	redisValueBytes, err := json.Marshal(expected)
	if err != nil {
		s.Logger.Error("unexpected error faking data for a pet", zap.Error(err))
		s.Require().Nil(err)
	}
	redisValue := string(redisValueBytes)

	// Act
	s.Logger.Debug("setting a value in redis", zap.Any("redisInput", expected))
	err = s.RedisClient.Set(s.Context, redisKey, redisValue, 0).Err()
	if err != nil {
		s.Logger.Error("unexpected error setting a redis value", zap.Error(err))
		s.Require().Nil(err, "no errors should happen while setting a redis value")
	}
	result, err := s.RedisClient.Get(s.Context, redisKey).Result()
	if err != nil {
		s.Logger.Error("unexpected error getting a redis value", zap.Error(err))
		s.Require().Nil(err, "no errors should happen while reading from Redis")
	}

	// Assert
	s.Require().NotEmpty(result, "something should be read from Redis")
	s.Logger.Debug("read a pet from redis", zap.String("redisOutput", result))
	got := Pet{}
	err = json.Unmarshal([]byte(result), &got)
	s.Nil(err, "parsing back from json into Pet should be possible")
	s.Equal(expected, got, "the recovered artifact should be equal to the stored one")
}

// TestGettingAnExpiredValueReturnsNil demonstrates how reading a value that has
// expired returns an empty string.
func (s *RedisSuite) TestGettingAnExpiredValueReturnsNil() {
	expected := faker.Sentence()
	redisKey := faker.Word()
	timeToLive := time.Millisecond * 250
	s.Logger.Debug("setting a value into a key", zap.String("redisKey", redisKey), zap.String(
		"redisInput",
		expected,
	))
	err := s.RedisClient.Set(s.Context, redisKey, expected, timeToLive).Err()
	if err != nil {
		s.Logger.Error("unexpected error setting a redis value", zap.Error(err))
	}

	// Act
	s.Logger.Debug("waiting until the value is expired")
	time.Sleep(timeToLive + 1)
	s.Logger.Debug("reading a value from Redis", zap.String("redisKey", redisKey))
	got, err := s.RedisClient.Get(s.Context, redisKey).Result()

	// Assert
	s.NotNil(err, "no errors should happen when reading a known value")
	s.Empty(got, "the returned value should be empty")
}

// Pet holds data related to pets and Redis testing!
type Pet struct {
	Name    string
	IsAngry bool
}

func TestRedisSuite(t *testing.T) {
	SkipTestIfMinikubeIsUnavailable(t)
	// Delegate to testify's suite
	suite.Run(t, new(RedisSuite))
}
