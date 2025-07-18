package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cliffordotieno/ai-context-gap-tracker/internal/config"
	"github.com/go-redis/redis/v8"
)

// Client wraps redis.Client
type Client struct {
	*redis.Client
}

// NewClient creates a new Redis client
func NewClient(cfg config.RedisConfig) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address(),
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	log.Println("Successfully connected to Redis")

	return &Client{Client: rdb}, nil
}

// SetContext stores context data in Redis
func (c *Client) SetContext(ctx context.Context, sessionID string, turnNumber int, data interface{}) error {
	key := fmt.Sprintf("context:%s:%d", sessionID, turnNumber)
	return c.Set(ctx, key, data, 24*time.Hour).Err()
}

// GetContext retrieves context data from Redis
func (c *Client) GetContext(ctx context.Context, sessionID string, turnNumber int) (string, error) {
	key := fmt.Sprintf("context:%s:%d", sessionID, turnNumber)
	return c.Get(ctx, key).Result()
}

// SetSession stores session data in Redis
func (c *Client) SetSession(ctx context.Context, sessionID string, data interface{}) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return c.Set(ctx, key, data, 24*time.Hour).Err()
}

// GetSession retrieves session data from Redis
func (c *Client) GetSession(ctx context.Context, sessionID string) (string, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	return c.Get(ctx, key).Result()
}

// SetMemoryGraph stores memory graph in Redis
func (c *Client) SetMemoryGraph(ctx context.Context, sessionID string, graph interface{}) error {
	key := fmt.Sprintf("memory:%s", sessionID)
	return c.Set(ctx, key, graph, 24*time.Hour).Err()
}

// GetMemoryGraph retrieves memory graph from Redis
func (c *Client) GetMemoryGraph(ctx context.Context, sessionID string) (string, error) {
	key := fmt.Sprintf("memory:%s", sessionID)
	return c.Get(ctx, key).Result()
}

// InvalidateSession removes session-related data from Redis
func (c *Client) InvalidateSession(ctx context.Context, sessionID string) error {
	keys := []string{
		fmt.Sprintf("session:%s", sessionID),
		fmt.Sprintf("memory:%s", sessionID),
	}

	// Get all context keys for the session
	contextKeys, err := c.Keys(ctx, fmt.Sprintf("context:%s:*", sessionID)).Result()
	if err != nil {
		return err
	}

	keys = append(keys, contextKeys...)

	if len(keys) > 0 {
		return c.Del(ctx, keys...).Err()
	}

	return nil
}

// Close closes the Redis connection
func (c *Client) Close() error {
	return c.Client.Close()
}