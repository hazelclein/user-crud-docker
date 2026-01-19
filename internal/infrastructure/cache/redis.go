package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"user-crud/internal/domain"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(host, port string, ttl time.Duration) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Password:     "", // no password
		DB:           0,  // default DB
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisCache{
		client: client,
		ttl:    ttl,
	}, nil
}

// GetUser gets user from cache
func (c *RedisCache) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	key := fmt.Sprintf("user:%d", id)

	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, err
	}

	var user domain.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// SetUser sets user in cache
func (c *RedisCache) SetUser(ctx context.Context, user *domain.User) error {
	key := fmt.Sprintf("user:%d", user.ID)

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, c.ttl).Err()
}

// DeleteUser deletes user from cache
func (c *RedisCache) DeleteUser(ctx context.Context, id int64) error {
	key := fmt.Sprintf("user:%d", id)
	return c.client.Del(ctx, key).Err()
}

// Clear clears all cache
func (c *RedisCache) Clear(ctx context.Context) error {
	return c.client.FlushDB(ctx).Err()
}

// Close closes redis connection
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// Ping checks redis connection
func (c *RedisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}