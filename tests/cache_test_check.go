package api_check_suite

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/diagnoxix/core/domain"
	"github.com/diagnoxix/core/services"
	"github.com/diagnoxix/core/services/cache"
	"github.com/diagnoxix/core/utils"
	"github.com/joho/godotenv"
)
func TestCacheImplementation(t *testing.T) {
	// Load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Initialize logger for testing
	logConfig := utils.LogConfig{
		Level:       "info",
		OutputPath:  "stdout",
		Development: true,
	}
	if err := utils.InitLogger(logConfig); err != nil {
		log.Printf("Warning: Failed to initialize logger: %v", err)
	}

	if err := utils.InitLogger(logConfig); err != nil {
		log.Printf("Warning: Failed to initialize logger: %v", err)
	}

	fmt.Println("🧪 Testing AI Caching Implementation")
	fmt.Println("====================================")

	// Test 1: Cache initialization
	fmt.Println("\n1. Testing Cache Initialization...")
	testCacheInitialization()

	// Test 2: Cache operations
	fmt.Println("\n2. Testing Cache Operations...")
	testCacheOperations()

	// Test 3: AI Service with cache
	fmt.Println("\n3. Testing AI Service with Cache...")
	testAIServiceWithCache()

	// Test 4: Cache endpoints (if server is running)
	fmt.Println("\n4. Testing Cache API Endpoints...")
	testCacheEndpoints()

	fmt.Println("\n🎉 Cache testing completed!")
}

func testCacheInitialization() {
	// Test with Redis URL
	config := cache.CacheConfig{
		RedisURL:     "redis://localhost:6379",
		DefaultTTL:   1 * time.Hour,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
	}

	aiCache, err := cache.NewAICache(config)
	if err != nil {
		fmt.Printf("❌ Redis cache initialization failed (expected if Redis not running): %v\n", err)
	} else {
		fmt.Println("✅ Redis cache initialized successfully")
		defer aiCache.Close()
	}

	// Test fallback to in-memory cache
	config.RedisURL = "invalid://url"
	memCache, err := cache.NewAICache(config)
	if err != nil {
		fmt.Printf("❌ Memory cache fallback failed: %v\n", err)
	} else {
		fmt.Println("✅ Memory cache fallback working")
		defer memCache.Close()
	}
}

func testCacheOperations() {
	// Create in-memory cache for testing
	config := cache.CacheConfig{
		RedisURL:   "", // Empty to force in-memory
		DefaultTTL: 1 * time.Hour,
	}

	aiCache, err := cache.NewAICache(config)
	if err != nil {
		fmt.Printf("❌ Cache creation failed: %v\n", err)
		return
	}
	defer aiCache.Close()

	ctx := context.Background()

	// Test Set and Get
	testData := map[string]interface{}{
		"test": "data",
		"number": 42,
	}

	err = aiCache.Set(ctx, "test-key", testData, 5*time.Minute)
	if err != nil {
		fmt.Printf("❌ Cache set failed: %v\n", err)
		return
	}
	fmt.Println("✅ Cache set operation successful")

	// Test Get
	retrieved, found, err := aiCache.Get(ctx, "test-key")
	if err != nil {
		fmt.Printf("❌ Cache get failed: %v\n", err)
		return
	}

	if !found {
		fmt.Println("❌ Cache miss when hit expected")
		return
	}

	if retrieved == nil {
		fmt.Println("❌ Retrieved data is nil")
		return
	}

	fmt.Println("✅ Cache get operation successful")

	// Test cache key generation
	cacheKey := aiCache.GenerateCacheKey("test-operation", testData)
	if cacheKey == "" {
		fmt.Println("❌ Cache key generation failed")
		return
	}
	fmt.Printf("✅ Cache key generated: %s\n", cacheKey[:16]+"...")

	// Test stats
	stats, err := aiCache.GetStats(ctx)
	if err != nil {
		fmt.Printf("❌ Cache stats failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Cache stats retrieved: %d total keys\n", stats.TotalKeys)
}

func testAIServiceWithCache() {
	openAIKey := os.Getenv("OPEN_API_KEY")
	if openAIKey == "" {
		fmt.Println("⏭️ Skipping AI service test (no OpenAI key)")
		return
	}

	// Create AI service with in-memory cache
	config := cache.CacheConfig{
		RedisURL:   "", // Force in-memory for testing
		DefaultTTL: 1 * time.Hour,
	}

	aiCache, err := cache.NewAICache(config)
	if err != nil {
		fmt.Printf("❌ Cache creation failed: %v\n", err)
		return
	}
	defer aiCache.Close()

	aiService := services.NewAIServiceWithCache(openAIKey, aiCache)

	// Test lab interpretation with caching
	labTest := domain.LabTest{
		PatientID:          "test-patient-cache",
		DiagnosticCentreID: "test-center-cache",
		TestName:           "Cache Test CBC",
		Results: map[string]string{
			"Hemoglobin": "12.5 g/dL",
		},
		ReferenceRanges: map[string]string{
			"Hemoglobin": "12.0-15.5 g/dL",
		},
	}

	ctx := context.Background()

	// First call - should miss cache and call OpenAI
	fmt.Println("Making first AI call (cache miss expected)...")
	start := time.Now()
	_, err = aiService.InterpretLabResults(ctx, labTest)
	firstCallDuration := time.Since(start)

	if err != nil {
		fmt.Printf("❌ First AI call failed: %v\n", err)
		return
	}
	fmt.Printf("✅ First AI call successful (%v)\n", firstCallDuration)

	// Second call - should hit cache
	fmt.Println("Making second AI call (cache hit expected)...")
	start = time.Now()
	_, err = aiService.InterpretLabResults(ctx, labTest)
	secondCallDuration := time.Since(start)

	if err != nil {
		fmt.Printf("❌ Second AI call failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Second AI call successful (%v)\n", secondCallDuration)

	// Check if second call was significantly faster (cache hit)
	if secondCallDuration < firstCallDuration/2 {
		fmt.Println("✅ Cache hit detected - second call much faster!")
	} else {
		fmt.Println("⚠️ Cache might not be working - similar response times")
	}

	// Get cache stats
	stats, err := aiService.GetCacheStats(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to get cache stats: %v\n", err)
		return
	}

	fmt.Printf("✅ Cache stats: %.1f%% hit rate, %d total operations\n", 
		stats.HitRate, stats.Hits+stats.Misses)
}

func testCacheEndpoints() {
	baseURL := "http://localhost:7556"
	client := &http.Client{Timeout: 10 * time.Second}

	// Test cache health endpoint
	resp, err := client.Get(baseURL + "/v1/cache/health")
	if err != nil {
		fmt.Printf("❌ Server not running or cache health endpoint failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Cache health endpoint working")
	} else {
		fmt.Printf("❌ Cache health endpoint returned status: %d\n", resp.StatusCode)
		return
	}

	// Test cache stats endpoint (requires auth, so expect 401)
	resp, err = client.Get(baseURL + "/v1/cache/stats")
	if err != nil {
		fmt.Printf("❌ Cache stats endpoint failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("✅ Cache stats endpoint properly protected")
	} else {
		fmt.Printf("⚠️ Cache stats endpoint returned unexpected status: %d\n", resp.StatusCode)
	}

	fmt.Println("✅ Cache API endpoints are accessible")
}

func makeRequest(method, url string, body interface{}) (*http.Response, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return client.Do(req)
}
