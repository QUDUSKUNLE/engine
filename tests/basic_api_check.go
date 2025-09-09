package api_check_suite

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
)

func RunBasicAPITests() {
	// Load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	baseURL := getEnv("TEST_BASE_URL", "http://localhost:7556")
	client := &http.Client{Timeout: 10 * time.Second}

	fmt.Println("🧪 Running Basic API Tests")
	fmt.Println("==========================")

	// Test 1: Health Check
	fmt.Print("1. Testing Health Check... ")
	if testHealthCheck(client, baseURL) {
		fmt.Println("✅ PASS")
	} else {
		fmt.Println("❌ FAIL")
	}

	// Test 2: Home Page
	fmt.Print("2. Testing Home Page... ")
	if testHomePage(client, baseURL) {
		fmt.Println("✅ PASS")
	} else {
		fmt.Println("❌ FAIL")
	}

	// Test 3: AI Capabilities (no auth required)
	fmt.Print("3. Testing AI Capabilities... ")
	if testAICapabilities(client, baseURL) {
		fmt.Println("✅ PASS")
	} else {
		fmt.Println("❌ FAIL")
	}

	// Test 4: WebSocket Stats
	fmt.Print("4. Testing WebSocket Stats... ")
	if testWebSocketStats(client, baseURL) {
		fmt.Println("✅ PASS")
	} else {
		fmt.Println("❌ FAIL")
	}

	// Test 5: Swagger Documentation
	fmt.Print("5. Testing Swagger Docs... ")
	if testSwaggerDocs(client, baseURL) {
		fmt.Println("✅ PASS")
	} else {
		fmt.Println("❌ FAIL")
	}

	// Test 6: Metrics Endpoint
	fmt.Print("6. Testing Metrics Endpoint... ")
	if testMetrics(client, baseURL) {
		fmt.Println("✅ PASS")
	} else {
		fmt.Println("❌ FAIL")
	}

	fmt.Println("\n🎉 Basic API tests completed!")
	fmt.Println("\n💡 To test authenticated endpoints, run:")
	fmt.Println("   go run api_test_suite.go")
	fmt.Println("   (Make sure to set TEST_JWT_TOKEN environment variable)")
}

func testHealthCheck(client *http.Client, baseURL string) bool {
	resp, err := client.Get(baseURL + "/health")
	if err != nil {
		fmt.Printf("❌ Error: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Status: %d", resp.StatusCode)
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Read error: %v", err)
		return false
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("❌ JSON error: %v", err)
		return false
	}

	return true
}

func testHomePage(client *http.Client, baseURL string) bool {
	resp, err := client.Get(baseURL + "/")
	if err != nil {
		fmt.Printf("❌ Error: %v", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func testAICapabilities(client *http.Client, baseURL string) bool {
	resp, err := client.Get(baseURL + "/v1/ai/capabilities")
	if err != nil {
		fmt.Printf("❌ Error: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Status: %d", resp.StatusCode)
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Read error: %v", err)
		return false
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("❌ JSON error: %v", err)
		return false
	}

	// Check if we have the expected AI features
	features, ok := response["features"].([]interface{})
	if !ok {
		fmt.Printf("❌ No features found")
		return false
	}

	if len(features) < 7 {
		fmt.Printf("❌ Expected at least 7 features, got %d", len(features))
		return false
	}

	// Print feature count for verification
	fmt.Printf("(%d features) ", len(features))
	return true
}

func testWebSocketStats(client *http.Client, baseURL string) bool {
	resp, err := client.Get(baseURL + "/v1/ws/stats")
	if err != nil {
		fmt.Printf("❌ Error: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Status: %d", resp.StatusCode)
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Read error: %v", err)
		return false
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("❌ JSON error: %v", err)
		return false
	}

	// Check if response has expected structure
	if data, ok := response["data"].(map[string]interface{}); ok {
		if _, hasClients := data["connected_clients"]; hasClients {
			return true
		}
	}

	fmt.Printf("❌ Invalid response structure")
	return false
}

func testSwaggerDocs(client *http.Client, baseURL string) bool {
	resp, err := client.Get(baseURL + "/swagger/")
	if err != nil {
		fmt.Printf("❌ Error: %v", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func testMetrics(client *http.Client, baseURL string) bool {
	resp, err := client.Get(baseURL + "/metrics")
	if err != nil {
		fmt.Printf("❌ Error: %v", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

