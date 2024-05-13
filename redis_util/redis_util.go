package redis_util

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient represents a Redis client instance.
type RedisClient struct {
	client *redis.Client
}

var redisClient *RedisClient

// Initialize initializes the Redis client.
func Initialize(addr, password string, db int) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	redisClient = &RedisClient{client: rdb}

	// Listen for interrupt signals to close the Redis connection gracefully
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-interrupt
		if err := redisClient.Close(); err != nil {
			fmt.Println("Error closing Redis connection:", err)
		}
	}()
}

// NewRedisClient creates a new Redis client instance.
func NewRedisClient() *RedisClient {
	return redisClient
}

// Ping pings the Redis server to check if the connection is successful.
func (rc *RedisClient) Ping() error {
	pong, err := rc.client.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("failed to ping Redis server: %v", err)
	}
	fmt.Println("Connected to Redis:", pong)
	return nil
}

// Set sets a key-value pair in Redis with an optional expiration time.
func (rc *RedisClient) Set(key, value string, expiration time.Duration) error {
	err := rc.client.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key '%s': %v", key, err)
	}
	return nil
}

// Get gets the value of a key from Redis.
func (rc *RedisClient) Get(key string) (string, error) {
	val, err := rc.client.Get(context.Background(), key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get key '%s': %v", key, err)
	}
	return val, nil
}

// Invalidate invalidates a key in Redis by deleting it.
func (rc *RedisClient) Invalidate(key string) error {
	err := rc.client.Del(context.Background(), key).Err()
	if err != nil {
		return fmt.Errorf("failed to invalidate key '%s': %v", key, err)
	}
	return nil
}

// Close closes the Redis client connection.
func (rc *RedisClient) Close() error {
	err := rc.client.Close()
	if err != nil {
		return fmt.Errorf("failed to close Redis connection: %v", err)
	}
	return nil
}
