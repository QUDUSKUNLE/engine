package cache

import (
	"sync/atomic"
	"time"

	"github.com/diagnoxix/core/utils"
)

// CacheMetrics tracks cache performance
type (
		CacheMetrics struct {
		hits        int64
		misses      int64
		errors      int64
		sets        int64
		deletes     int64
		totalTime   int64 // in nanoseconds
		operations  int64
	}
	// Enhanced CacheStats with additional metrics
	CacheStats struct {
		Hits            int64   `json:"hits"`
		Misses          int64   `json:"misses"`
		Errors          int64   `json:"errors"`
		Sets            int64   `json:"sets"`
		Deletes         int64   `json:"deletes"`
		HitRate         float64 `json:"hit_rate"`
		TotalOperations int64   `json:"total_operations"`
		AvgResponseTime float64 `json:"avg_response_time_ms"`
		TotalKeys       int64   `json:"total_keys"`
		MemoryUsage     int64   `json:"memory_usage"`
	}
)

// NewCacheMetrics creates a new metrics instance
func NewCacheMetrics() *CacheMetrics {
	return &CacheMetrics{}
}

// RecordHit records a cache hit
func (m *CacheMetrics) RecordHit(duration time.Duration) {
	atomic.AddInt64(&m.hits, 1)
	atomic.AddInt64(&m.totalTime, duration.Nanoseconds())
	atomic.AddInt64(&m.operations, 1)
}

// RecordMiss records a cache miss
func (m *CacheMetrics) RecordMiss(duration time.Duration) {
	atomic.AddInt64(&m.misses, 1)
	atomic.AddInt64(&m.totalTime, duration.Nanoseconds())
	atomic.AddInt64(&m.operations, 1)
}

// RecordError records a cache error
func (m *CacheMetrics) RecordError(duration time.Duration) {
	atomic.AddInt64(&m.errors, 1)
	atomic.AddInt64(&m.totalTime, duration.Nanoseconds())
	atomic.AddInt64(&m.operations, 1)
}

// RecordSet records a cache set operation
func (m *CacheMetrics) RecordSet(duration time.Duration) {
	atomic.AddInt64(&m.sets, 1)
	atomic.AddInt64(&m.totalTime, duration.Nanoseconds())
	atomic.AddInt64(&m.operations, 1)
}

// RecordDelete records a cache delete operation
func (m *CacheMetrics) RecordDelete(duration time.Duration) {
	atomic.AddInt64(&m.deletes, 1)
	atomic.AddInt64(&m.totalTime, duration.Nanoseconds())
	atomic.AddInt64(&m.operations, 1)
}

// GetStats returns current cache statistics
func (m *CacheMetrics) GetStats() CacheStats {
	hits := atomic.LoadInt64(&m.hits)
	misses := atomic.LoadInt64(&m.misses)
	errors := atomic.LoadInt64(&m.errors)
	sets := atomic.LoadInt64(&m.sets)
	deletes := atomic.LoadInt64(&m.deletes)
	totalTime := atomic.LoadInt64(&m.totalTime)
	operations := atomic.LoadInt64(&m.operations)

	total := hits + misses
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(hits) / float64(total) * 100
	}

	avgResponseTime := float64(0)
	if operations > 0 {
		avgResponseTime = float64(totalTime) / float64(operations) / 1e6 // Convert to milliseconds
	}

	return CacheStats{
		Hits:            hits,
		Misses:          misses,
		Errors:          errors,
		Sets:            sets,
		Deletes:         deletes,
		HitRate:         hitRate,
		TotalOperations: operations,
		AvgResponseTime: avgResponseTime,
	}
}

// LogStats logs current cache statistics
func (m *CacheMetrics) LogStats() {
	stats := m.GetStats()
	
	utils.Info("Cache performance statistics",
		utils.LogField{Key: "hits", Value: stats.Hits},
		utils.LogField{Key: "misses", Value: stats.Misses},
		utils.LogField{Key: "errors", Value: stats.Errors},
		utils.LogField{Key: "hit_rate", Value: stats.HitRate},
		utils.LogField{Key: "avg_response_time_ms", Value: stats.AvgResponseTime},
		utils.LogField{Key: "total_operations", Value: stats.TotalOperations})
}

// Reset resets all metrics
func (m *CacheMetrics) Reset() {
	atomic.StoreInt64(&m.hits, 0)
	atomic.StoreInt64(&m.misses, 0)
	atomic.StoreInt64(&m.errors, 0)
	atomic.StoreInt64(&m.sets, 0)
	atomic.StoreInt64(&m.deletes, 0)
	atomic.StoreInt64(&m.totalTime, 0)
	atomic.StoreInt64(&m.operations, 0)
}
