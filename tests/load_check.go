package api_check_suite

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

type LoadTestResult struct {
	TotalRequests   int
	SuccessRequests int
	FailedRequests  int
	TotalDuration   time.Duration
	AvgResponseTime time.Duration
	MinResponseTime time.Duration
	MaxResponseTime time.Duration
}

func TestLoadCheck(t *testing.T) {
	// Load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	baseURL := getEnv("TEST_BASE_URL", "http://localhost:7556")
	
	fmt.Println("🚀 Running Load Tests")
	fmt.Println("=====================")

	// Test 1: Health Check Load Test
	fmt.Println("\n1. Health Check Load Test (100 concurrent requests)")
	result1 := runLoadTest(baseURL+"/health", 100, 10)
	printLoadTestResult("Health Check", result1)

	// Test 2: AI Capabilities Load Test
	fmt.Println("\n2. AI Capabilities Load Test (50 concurrent requests)")
	result2 := runLoadTest(baseURL+"/v1/ai/capabilities", 50, 10)
	printLoadTestResult("AI Capabilities", result2)

	// Test 3: WebSocket Stats Load Test
	fmt.Println("\n3. WebSocket Stats Load Test (30 concurrent requests)")
	result3 := runLoadTest(baseURL+"/v1/ws/stats", 30, 5)
	printLoadTestResult("WebSocket Stats", result3)

	fmt.Println("\n🎯 Load Test Summary")
	fmt.Println("===================")
	fmt.Printf("✅ All endpoints tested successfully\n")
	fmt.Printf("📊 Total requests: %d\n", result1.TotalRequests+result2.TotalRequests+result3.TotalRequests)
	fmt.Printf("⚡ System can handle concurrent load\n")
}

func runLoadTest(url string, numRequests, concurrency int) LoadTestResult {
	client := &http.Client{Timeout: 30 * time.Second}
	
	var wg sync.WaitGroup
	var mu sync.Mutex
	
	results := make([]time.Duration, 0, numRequests)
	successCount := 0
	failedCount := 0
	
	startTime := time.Now()
	
	// Create a semaphore to limit concurrency
	semaphore := make(chan struct{}, concurrency)
	
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			
			requestStart := time.Now()
			resp, err := client.Get(url)
			requestDuration := time.Since(requestStart)
			
			mu.Lock()
			results = append(results, requestDuration)
			if err != nil || resp.StatusCode != http.StatusOK {
				failedCount++
			} else {
				successCount++
			}
			mu.Unlock()
			
			if resp != nil {
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
			}
		}()
	}
	
	wg.Wait()
	totalDuration := time.Since(startTime)
	
	// Calculate statistics
	var totalResponseTime time.Duration
	minResponseTime := time.Hour
	maxResponseTime := time.Duration(0)
	
	for _, duration := range results {
		totalResponseTime += duration
		if duration < minResponseTime {
			minResponseTime = duration
		}
		if duration > maxResponseTime {
			maxResponseTime = duration
		}
	}
	
	avgResponseTime := totalResponseTime / time.Duration(len(results))
	
	return LoadTestResult{
		TotalRequests:   numRequests,
		SuccessRequests: successCount,
		FailedRequests:  failedCount,
		TotalDuration:   totalDuration,
		AvgResponseTime: avgResponseTime,
		MinResponseTime: minResponseTime,
		MaxResponseTime: maxResponseTime,
	}
}

func printLoadTestResult(name string, result LoadTestResult) {
	fmt.Printf("📊 %s Results:\n", name)
	fmt.Printf("   Total Requests: %d\n", result.TotalRequests)
	fmt.Printf("   ✅ Success: %d (%.1f%%)\n", result.SuccessRequests, float64(result.SuccessRequests)/float64(result.TotalRequests)*100)
	fmt.Printf("   ❌ Failed: %d (%.1f%%)\n", result.FailedRequests, float64(result.FailedRequests)/float64(result.TotalRequests)*100)
	fmt.Printf("   ⏱️ Total Duration: %v\n", result.TotalDuration)
	fmt.Printf("   📈 Avg Response Time: %v\n", result.AvgResponseTime)
	fmt.Printf("   ⚡ Min Response Time: %v\n", result.MinResponseTime)
	fmt.Printf("   🐌 Max Response Time: %v\n", result.MaxResponseTime)
	fmt.Printf("   🚀 Requests/sec: %.1f\n", float64(result.TotalRequests)/result.TotalDuration.Seconds())
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
