package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/diagnoxix/core/utils"
	"github.com/redis/go-redis/v9"
)

// CacheConfig holds cache configuration
type CacheConfig struct {
	RedisURL     string        `yaml:"redis_url" default:"redis://localhost:6379"`
	DefaultTTL   time.Duration `yaml:"default_ttl" default:"1h"`
	MaxRetries   int           `yaml:"max_retries" default:"3"`
	DialTimeout  time.Duration `yaml:"dial_timeout" default:"5s"`
	ReadTimeout  time.Duration `yaml:"read_timeout" default:"3s"`
	WriteTimeout time.Duration `yaml:"write_timeout" default:"3s"`
	PoolSize     int           `yaml:"pool_size" default:"10"`
}

// AICache provides caching functionality for AI responses
type AICache struct {
	client     *redis.Client
	config     CacheConfig
	keyPrefix  string
	fallback   bool // Use in-memory fallback if Redis is unavailable
	memCache   map[string]*CacheEntry
	memMutex   sync.RWMutex
}

// CacheEntry represents a cached item
type CacheEntry struct {
	Data      interface{} `json:"data"`
	ExpiresAt time.Time   `json:"expires_at"`
	CreatedAt time.Time   `json:"created_at"`
	Hash      string      `json:"hash"`
}

// NewAICache creates a new AI cache instance
func NewAICache(config CacheConfig) (*AICache, error) {
	// Parse Redis URL
	opt, err := redis.ParseURL(config.RedisURL)
	if err != nil {
		utils.Warn("Failed to parse Redis URL, using in-memory cache only", 
			utils.LogField{Key: "error", Value: err.Error()})
		
		return &AICache{
			config:    config,
			keyPrefix: "diagnoxix:ai:",
			fallback:  true,
			memCache:  make(map[string]*CacheEntry),
		}, nil
	}

	// Configure Redis client
	opt.DialTimeout = config.DialTimeout
	opt.ReadTimeout = config.ReadTimeout
	opt.WriteTimeout = config.WriteTimeout
	opt.PoolSize = config.PoolSize
	opt.MaxRetries = config.MaxRetries

	client := redis.NewClient(opt)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		utils.Warn("Redis connection failed, using in-memory cache only", 
			utils.LogField{Key: "error", Value: err.Error()})
		
		return &AICache{
			config:    config,
			keyPrefix: "diagnoxix:ai:",
			fallback:  true,
			memCache:  make(map[string]*CacheEntry),
		}, nil
	}

	utils.Info("Redis cache initialized successfully")

	return &AICache{
		client:    client,
		config:    config,
		keyPrefix: "diagnoxix:ai:",
		fallback:  false,
		memCache:  make(map[string]*CacheEntry),
	}, nil
}

// Get retrieves a value from cache
func (c *AICache) Get(ctx context.Context, key string) (interface{}, bool, error) {
	fullKey := c.keyPrefix + key

	if c.fallback {
		return c.getFromMemory(fullKey)
	}

	// Try Redis first
	data, err := c.client.Get(ctx, fullKey).Result()
	if err != nil {
		if err == redis.Nil {
			// Cache miss
			return nil, false, nil
		}
		
		utils.Error("Redis get error", utils.LogField{Key: "error", Value: err.Error()})
		
		// Fallback to memory cache
		return c.getFromMemory(fullKey)
	}

	// Deserialize data
	var entry CacheEntry
	if err := json.Unmarshal([]byte(data), &entry); err != nil {
		utils.Error("Cache deserialization error", utils.LogField{Key: "error", Value: err.Error()})
		return nil, false, err
	}

	// Check expiration
	if time.Now().After(entry.ExpiresAt) {
		// Expired, delete from cache
		c.client.Del(ctx, fullKey)
		return nil, false, nil
	}

	return entry.Data, true, nil
}

// Set stores a value in cache
func (c *AICache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	fullKey := c.keyPrefix + key

	if ttl == 0 {
		ttl = c.config.DefaultTTL
	}

	entry := CacheEntry{
		Data:      value,
		ExpiresAt: time.Now().Add(ttl),
		CreatedAt: time.Now(),
		Hash:      c.generateHash(value),
	}

	if c.fallback {
		return c.setInMemory(fullKey, &entry)
	}

	// Serialize data
	data, err := json.Marshal(entry)
	if err != nil {
		utils.Error("Cache serialization error", utils.LogField{Key: "error", Value: err.Error()})
		return err
	}

	// Store in Redis
	if err := c.client.Set(ctx, fullKey, data, ttl).Err(); err != nil {
		utils.Error("Redis set error", utils.LogField{Key: "error", Value: err.Error()})
		
		// Fallback to memory cache
		return c.setInMemory(fullKey, &entry)
	}

	return nil
}

// Delete removes a value from cache
func (c *AICache) Delete(ctx context.Context, key string) error {
	fullKey := c.keyPrefix + key

	if c.fallback {
		return c.deleteFromMemory(fullKey)
	}

	if err := c.client.Del(ctx, fullKey).Err(); err != nil {
		utils.Error("Redis delete error", utils.LogField{Key: "error", Value: err.Error()})
		return err
	}

	// Also delete from memory cache if exists
	c.deleteFromMemory(fullKey)
	return nil
}

// Clear removes all AI cache entries
func (c *AICache) Clear(ctx context.Context) error {
	if c.fallback {
		c.memMutex.Lock()
		c.memCache = make(map[string]*CacheEntry)
		c.memMutex.Unlock()
		return nil
	}

	// Delete all keys with our prefix
	keys, err := c.client.Keys(ctx, c.keyPrefix+"*").Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}

	return nil
}

// GetStats returns cache performance statistics
func (c *AICache) GetStats(ctx context.Context) (*CacheStats, error) {
	stats := &CacheStats{}

	if c.fallback {
		c.memMutex.RLock()
		stats.TotalKeys = int64(len(c.memCache))
		c.memMutex.RUnlock()
		return stats, nil
	}

	// Get Redis info
	_, err := c.client.Info(ctx, "stats", "memory").Result()
	if err != nil {
		return nil, err
	}

	// Parse Redis stats (simplified)
	// In production, you'd parse the info string properly
	keys, err := c.client.DBSize(ctx).Result()
	if err == nil {
		stats.TotalKeys = keys
	}

	return stats, nil
}

// GenerateCacheKey creates a cache key from input data
func (c *AICache) GenerateCacheKey(operation string, input interface{}) string {
	hash := c.generateHash(input)
	return fmt.Sprintf("%s:%s", operation, hash)
}

// Memory cache methods (fallback)
func (c *AICache) getFromMemory(key string) (interface{}, bool, error) {
	c.memMutex.RLock()
	defer c.memMutex.RUnlock()

	entry, exists := c.memCache[key]
	if !exists {
		return nil, false, nil
	}

	// Check expiration
	if time.Now().After(entry.ExpiresAt) {
		delete(c.memCache, key)
		return nil, false, nil
	}

	return entry.Data, true, nil
}

func (c *AICache) setInMemory(key string, entry *CacheEntry) error {
	c.memMutex.Lock()
	defer c.memMutex.Unlock()

	c.memCache[key] = entry
	
	// Simple cleanup: remove expired entries periodically
	if len(c.memCache)%100 == 0 {
		go c.cleanupMemoryCache()
	}

	return nil
}

func (c *AICache) deleteFromMemory(key string) error {
	c.memMutex.Lock()
	defer c.memMutex.Unlock()

	delete(c.memCache, key)
	return nil
}

func (c *AICache) cleanupMemoryCache() {
	c.memMutex.Lock()
	defer c.memMutex.Unlock()

	now := time.Now()
	for key, entry := range c.memCache {
		if now.After(entry.ExpiresAt) {
			delete(c.memCache, key)
		}
	}
}

func (c *AICache) generateHash(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	hash := sha256.Sum256(jsonData)
	return hex.EncodeToString(hash[:])
}

// Close closes the cache connection
func (c *AICache) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
