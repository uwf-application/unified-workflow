package state

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisState implements StateManagement using Redis
type RedisState struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Addr     string        `json:"addr"`
	Password string        `json:"password"`
	DB       int           `json:"db"`
	Prefix   string        `json:"prefix"`
	TTL      time.Duration `json:"ttl"`
}

// NewRedisState creates a new Redis state management instance
func NewRedisState(config RedisConfig) (*RedisState, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	// Set default TTL if not provided
	ttl := config.TTL
	if ttl == 0 {
		ttl = 24 * time.Hour // Default 24 hours
	}

	return &RedisState{
		client: client,
		prefix: config.Prefix,
		ttl:    ttl,
	}, nil
}

// StoreResult stores an execution result in Redis
func (s *RedisState) StoreResult(ctx context.Context, runID string, result interface{}) error {
	key := s.getResultKey(runID)

	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	if err := s.client.Set(ctx, key, data, s.ttl).Err(); err != nil {
		return fmt.Errorf("failed to store result in Redis: %w", err)
	}

	return nil
}

// GetResult retrieves an execution result from Redis
func (s *RedisState) GetResult(ctx context.Context, runID string, result interface{}) (bool, error) {
	key := s.getResultKey(runID)

	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return false, nil // Result not found
		}
		return false, fmt.Errorf("failed to get result from Redis: %w", err)
	}

	if err := json.Unmarshal(data, result); err != nil {
		return false, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return true, nil
}

// DeleteResult deletes an execution result from Redis
func (s *RedisState) DeleteResult(ctx context.Context, runID string) error {
	key := s.getResultKey(runID)

	if err := s.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete result from Redis: %w", err)
	}

	return nil
}

// StoreExecutionStatus stores execution status in Redis
func (s *RedisState) StoreExecutionStatus(ctx context.Context, runID string, status interface{}) error {
	key := s.getStatusKey(runID)

	data, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("failed to marshal status: %w", err)
	}

	if err := s.client.Set(ctx, key, data, s.ttl).Err(); err != nil {
		return fmt.Errorf("failed to store status in Redis: %w", err)
	}

	return nil
}

// GetExecutionStatus retrieves execution status from Redis
func (s *RedisState) GetExecutionStatus(ctx context.Context, runID string, status interface{}) (bool, error) {
	key := s.getStatusKey(runID)

	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return false, nil // Status not found
		}
		return false, fmt.Errorf("failed to get status from Redis: %w", err)
	}

	if err := json.Unmarshal(data, status); err != nil {
		return false, fmt.Errorf("failed to unmarshal status: %w", err)
	}

	return true, nil
}

// StoreExecutionData stores execution data in Redis
func (s *RedisState) StoreExecutionData(ctx context.Context, runID string, data interface{}) error {
	key := s.getDataKey(runID)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := s.client.Set(ctx, key, jsonData, s.ttl).Err(); err != nil {
		return fmt.Errorf("failed to store data in Redis: %w", err)
	}

	return nil
}

// GetExecutionData retrieves execution data from Redis
func (s *RedisState) GetExecutionData(ctx context.Context, runID string, data interface{}) (bool, error) {
	key := s.getDataKey(runID)

	jsonData, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return false, nil // Data not found
		}
		return false, fmt.Errorf("failed to get data from Redis: %w", err)
	}

	if err := json.Unmarshal(jsonData, data); err != nil {
		return false, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return true, nil
}

// SetTTL sets TTL for a key
func (s *RedisState) SetTTL(ctx context.Context, runID string, ttl time.Duration) error {
	key := s.getResultKey(runID)

	if err := s.client.Expire(ctx, key, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set TTL: %w", err)
	}

	return nil
}

// Close closes the Redis connection
func (s *RedisState) Close() error {
	return s.client.Close()
}

// Helper methods for key generation
func (s *RedisState) getResultKey(runID string) string {
	return fmt.Sprintf("%s:result:%s", s.prefix, runID)
}

func (s *RedisState) getStatusKey(runID string) string {
	return fmt.Sprintf("%s:status:%s", s.prefix, runID)
}

func (s *RedisState) getDataKey(runID string) string {
	return fmt.Sprintf("%s:data:%s", s.prefix, runID)
}

// Implement StateManagement interface methods
func (s *RedisState) Store(ctx context.Context, key string, value interface{}) error {
	return s.StoreExecutionData(ctx, key, value)
}

func (s *RedisState) Retrieve(ctx context.Context, key string, value interface{}) (bool, error) {
	return s.GetExecutionData(ctx, key, value)
}

func (s *RedisState) Delete(ctx context.Context, key string) error {
	return s.DeleteResult(ctx, key)
}

func (s *RedisState) Exists(ctx context.Context, key string) (bool, error) {
	resultKey := s.getResultKey(key)
	exists, err := s.client.Exists(ctx, resultKey).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}
	return exists > 0, nil
}
