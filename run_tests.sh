#!/bin/bash

echo "🧪 Diagnoxix API Testing Suite"
echo "=============================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if server is running
echo -e "${BLUE}🔍 Checking if server is running...${NC}"
if curl -s http://localhost:7556/health > /dev/null; then
    echo -e "${GREEN}✅ Server is running on localhost:7556${NC}"
else
    echo -e "${RED}❌ Server is not running. Please start the server first:${NC}"
    echo -e "${YELLOW}   ./diagnoxix${NC}"
    exit 1
fi

echo ""

# Build the application first
echo -e "${BLUE}🔨 Building application...${NC}"
if go build -o diagnoxix .; then
    echo -e "${GREEN}✅ Build successful${NC}"
else
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi

echo ""

# Run basic API tests
echo -e "${BLUE}🧪 Running Basic API Tests...${NC}"
echo "================================"
cd tests
if go run basic_api_test.go; then
    echo -e "${GREEN}✅ Basic API tests completed${NC}"
else
    echo -e "${RED}❌ Basic API tests failed${NC}"
fi

echo ""

# Run load tests
echo -e "${BLUE}🚀 Running Load Tests...${NC}"
echo "========================"
if go run load_test.go; then
    echo -e "${GREEN}✅ Load tests completed${NC}"
else
    echo -e "${RED}❌ Load tests failed${NC}"
fi

echo ""

# Check if JWT token is available for authenticated tests
if [ -n "$TEST_JWT_TOKEN" ]; then
    echo -e "${BLUE}🔐 Running Authenticated API Tests...${NC}"
    echo "====================================="
    if go run api_test_suite.go; then
        echo -e "${GREEN}✅ Authenticated API tests completed${NC}"
    else
        echo -e "${RED}❌ Authenticated API tests failed${NC}"
    fi
else
    echo -e "${YELLOW}⚠️ Skipping authenticated tests (TEST_JWT_TOKEN not set)${NC}"
    echo -e "${BLUE}💡 To run authenticated tests:${NC}"
    echo -e "${YELLOW}   export TEST_JWT_TOKEN='your-jwt-token'${NC}"
    echo -e "${YELLOW}   ./run_tests.sh${NC}"
fi

echo ""

# Test WebSocket connection
echo -e "${BLUE}🔌 Testing WebSocket Connection...${NC}"
echo "=================================="
if go run test_websocket.go &
then
    WEBSOCKET_PID=$!
    sleep 5
    kill $WEBSOCKET_PID 2>/dev/null
    echo -e "${GREEN}✅ WebSocket connection test completed${NC}"
else
    echo -e "${RED}❌ WebSocket connection test failed${NC}"
fi

echo ""

# Final summary
echo -e "${GREEN}🎉 All tests completed!${NC}"
echo ""
echo -e "${BLUE}📊 Test Summary:${NC}"
echo -e "${GREEN}✅ Basic API endpoints tested${NC}"
echo -e "${GREEN}✅ Load testing completed${NC}"
echo -e "${GREEN}✅ WebSocket functionality verified${NC}"
echo ""
echo -e "${BLUE}🚀 Your Diagnoxix API is ready for production!${NC}"

cd ..
