#!/bin/bash

# AI Context Gap Tracker - System Test Script
# This script tests the basic functionality of the system

set -e

echo "🚀 AI Context Gap Tracker - System Test"
echo "========================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if services are running
check_service() {
    local service_name=$1
    local url=$2
    local expected_status=$3
    
    echo -n "🔍 Checking $service_name..."
    
    if curl -s -f "$url" > /dev/null 2>&1; then
        echo -e " ${GREEN}✅ Running${NC}"
        return 0
    else
        echo -e " ${RED}❌ Not running${NC}"
        return 1
    fi
}

# Test API endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo -n "🧪 Testing $description..."
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "%{http_code}" -o /tmp/response.json "$endpoint")
    else
        response=$(curl -s -w "%{http_code}" -X "$method" -H "Content-Type: application/json" -d "$data" -o /tmp/response.json "$endpoint")
    fi
    
    if [ "$response" = "200" ] || [ "$response" = "201" ]; then
        echo -e " ${GREEN}✅ Passed${NC}"
        return 0
    else
        echo -e " ${RED}❌ Failed (HTTP $response)${NC}"
        echo "Response: $(cat /tmp/response.json)"
        return 1
    fi
}

# Wait for services to be ready
wait_for_services() {
    echo "⏳ Waiting for services to start..."
    
    max_attempts=30
    attempt=0
    
    while [ $attempt -lt $max_attempts ]; do
        if check_service "Main Service" "http://localhost:8080/api/v1/health" 200 > /dev/null 2>&1; then
            echo -e "${GREEN}✅ Services are ready!${NC}"
            return 0
        fi
        
        echo -n "."
        sleep 2
        attempt=$((attempt + 1))
    done
    
    echo -e "${RED}❌ Services failed to start within timeout${NC}"
    return 1
}

# Main test function
run_tests() {
    echo "🧪 Running API Tests"
    echo "==================="
    
    # Test health endpoint
    test_endpoint "GET" "http://localhost:8080/api/v1/health" "" "Health Check"
    
    # Test context tracking
    test_endpoint "POST" "http://localhost:8080/api/v1/context/track" '{
        "session_id": "test-session-001",
        "turn_number": 1,
        "user_input": "I want to get a visa for a place I mentioned earlier."
    }' "Context Tracking"
    
    # Test rule evaluation
    test_endpoint "POST" "http://localhost:8080/api/v1/rules/evaluate" '{
        "session_id": "test-session-001",
        "turn_number": 1,
        "user_input": "I want to go there soon, but I am not sure about the requirements.",
        "entities": {},
        "topics": [],
        "timeline": [],
        "assertions": [],
        "ambiguities": [],
        "history": []
    }' "Rule Evaluation"
    
    # Test simple prompt rewrite
    test_endpoint "POST" "http://localhost:8080/api/v1/prompt/simple-rewrite" '{
        "session_id": "test-session-001",
        "turn_number": 2,
        "user_input": "Can you help me with that thing we discussed?"
    }' "Prompt Rewriting"
    
    # Test response auditing
    test_endpoint "POST" "http://localhost:8080/api/v1/audit/response" '{
        "session_id": "test-session-001",
        "turn_number": 1,
        "response_text": "Based on your previous mention, I assume you are asking about visa requirements.",
        "context": {}
    }' "Response Auditing"
    
    # Test pipeline processing
    test_endpoint "POST" "http://localhost:8080/api/v1/pipeline/process" '{
        "session_id": "test-session-001",
        "turn_number": 3,
        "user_input": "I need information about the process for that European country.",
        "system_prompt": "You are a helpful assistant for travel planning."
    }' "Pipeline Processing"
    
    # Test session context retrieval
    test_endpoint "GET" "http://localhost:8080/api/v1/context/session/test-session-001" "" "Session Context Retrieval"
    
    echo ""
    echo -e "${GREEN}✅ All tests completed!${NC}"
}

# Test NLP service if available
test_nlp_service() {
    echo "🧪 Testing NLP Service"
    echo "====================="
    
    if check_service "NLP Service" "http://localhost:5000/health" 200 > /dev/null 2>&1; then
        # Test NLP endpoints
        test_endpoint "POST" "http://localhost:5000/entities" '{
            "text": "I want to visit Paris, France next month for vacation.",
            "session_id": "test-session-001",
            "turn_number": 1
        }' "Entity Extraction"
        
        test_endpoint "POST" "http://localhost:5000/sentiment" '{
            "text": "I am very excited about this trip!",
            "session_id": "test-session-001",
            "turn_number": 1
        }' "Sentiment Analysis"
        
        test_endpoint "POST" "http://localhost:5000/ambiguities" '{
            "text": "Can you help me with that thing we discussed?",
            "session_id": "test-session-001",
            "turn_number": 1
        }' "Ambiguity Detection"
        
        echo -e "${GREEN}✅ NLP service tests completed!${NC}"
    else
        echo -e "${YELLOW}⚠️ NLP service not available, skipping NLP tests${NC}"
    fi
}

# Database test
test_database() {
    echo "🧪 Testing Database Connectivity"
    echo "==============================="
    
    # Test rule initialization
    test_endpoint "POST" "http://localhost:8080/api/v1/rules/initialize" "" "Default Rules Initialization"
    
    # Test rule retrieval
    test_endpoint "GET" "http://localhost:8080/api/v1/rules" "" "Rules Retrieval"
    
    echo -e "${GREEN}✅ Database tests completed!${NC}"
}

# Performance test
run_performance_test() {
    echo "🚀 Running Performance Test"
    echo "========================="
    
    echo "📊 Testing concurrent requests..."
    
    # Create a simple load test
    for i in {1..10}; do
        curl -s -X POST -H "Content-Type: application/json" \
            -d "{\"session_id\":\"perf-test-$i\",\"turn_number\":1,\"user_input\":\"Test message $i\"}" \
            "http://localhost:8080/api/v1/context/track" > /dev/null &
    done
    
    wait
    echo -e "${GREEN}✅ Performance test completed!${NC}"
}

# Main execution
main() {
    echo "Starting system tests..."
    
    # Check if Docker Compose is available
    if command -v docker-compose &> /dev/null; then
        echo "🐳 Docker Compose available"
    else
        echo -e "${YELLOW}⚠️ Docker Compose not found. Make sure services are running manually.${NC}"
    fi
    
    # Wait for services
    if ! wait_for_services; then
        echo -e "${RED}❌ Services not ready. Please start the services first.${NC}"
        echo "Run: docker-compose up -d"
        exit 1
    fi
    
    # Run tests
    run_tests
    
    # Test NLP service
    test_nlp_service
    
    # Test database
    test_database
    
    # Run performance test
    run_performance_test
    
    echo ""
    echo -e "${GREEN}🎉 All tests completed successfully!${NC}"
    echo "✅ AI Context Gap Tracker is working correctly"
}

# Cleanup function
cleanup() {
    echo "🧹 Cleaning up..."
    rm -f /tmp/response.json
}

# Set up cleanup on exit
trap cleanup EXIT

# Run main function
main "$@"